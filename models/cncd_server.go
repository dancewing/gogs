package models

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/frontend"
	"github.com/cncd/pipeline/pipeline/frontend/yaml"
	"github.com/cncd/pipeline/pipeline/frontend/yaml/compiler"
	"github.com/cncd/pipeline/pipeline/frontend/yaml/matrix"
	"github.com/cncd/pipeline/pipeline/rpc"
	"github.com/cncd/pubsub"
	"github.com/drone/envsubst"

	"crypto/sha256"
	"encoding/json"
	"strconv"

	"github.com/cncd/queue"

	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/pkg/setting"
	"github.com/gogits/gogs/pkg/sync"
	log "gopkg.in/clog.v1"
)

var BuildQueue = sync.NewUniqueQueue(setting.Webhook.QueueLength)

// return the metadata from the cli context.
func metadataFromStruct(repo *Repository, build, last *Build, proc *Proc, link string) frontend.Metadata {

	cl := repo.CloneLink()

	var remote string

	if setting.Repository.DisableHTTPGit {
		remote = cl.Git
	} else {
		remote = cl.HTTPS
	}

	return frontend.Metadata{
		Repo: frontend.Repo{
			Name:    repo.FullName(),
			Link:    repo.Link(),
			Remote:  remote,
			Private: repo.IsPrivate,
		},
		Curr: frontend.Build{
			Number:   build.Number,
			Parent:   build.Parent,
			Created:  build.Created,
			Started:  build.Started,
			Finished: build.Finished,
			Status:   build.Status,
			Event:    build.Event,
			Link:     build.Link,
			Target:   build.Deploy,
			Commit: frontend.Commit{
				Sha:     build.Commit,
				Ref:     build.Ref,
				Refspec: build.Refspec,
				Branch:  build.Branch,
				Message: build.Message,
				Author: frontend.Author{
					Name:   build.Author,
					Email:  build.Email,
					Avatar: build.Avatar,
				},
			},
		},
		Prev: frontend.Build{
			Number:   last.Number,
			Created:  last.Created,
			Started:  last.Started,
			Finished: last.Finished,
			Status:   last.Status,
			Event:    last.Event,
			Link:     last.Link,
			Target:   last.Deploy,
			Commit: frontend.Commit{
				Sha:     last.Commit,
				Ref:     last.Ref,
				Refspec: last.Refspec,
				Branch:  last.Branch,
				Message: last.Message,
				Author: frontend.Author{
					Name:   last.Author,
					Email:  last.Email,
					Avatar: last.Avatar,
				},
			},
		},
		Job: frontend.Job{
			Number: proc.PID,
			Matrix: proc.Environ,
		},
		Sys: frontend.System{
			Name: "drone",
			Link: link,
			Arch: "linux/amd64",
		},
	}
}

type cBuilder struct {
	Repo  *Repository
	Curr  *Build
	Last  *Build
	Netrc *Netrc
	Secs  []*Secret
	Regs  []*Registry
	Link  string
	Yaml  string
	Envs  map[string]string
}

type buildItem struct {
	Proc     *Proc
	Platform string
	Labels   map[string]string
	Config   *backend.Config
}

func (b *cBuilder) Build() ([]*buildItem, error) {

	axes, err := matrix.ParseString(b.Yaml)
	if err != nil {
		return nil, err
	}
	if len(axes) == 0 {
		axes = append(axes, matrix.Axis{})
	}

	var items []*buildItem
	for i, axis := range axes {
		proc := &Proc{
			BuildID: b.Curr.ID,
			PID:     i + 1,
			PGID:    i + 1,
			State:   StatusPending,
			Environ: axis,
		}

		metadata := metadataFromStruct(b.Repo, b.Curr, b.Last, proc, b.Link)
		environ := metadata.Environ()
		for k, v := range metadata.EnvironDrone() {
			environ[k] = v
		}
		for k, v := range axis {
			environ[k] = v
		}

		var secrets []compiler.Secret
		for _, sec := range b.Secs {
			if !sec.Match(b.Curr.Event) {
				continue
			}
			secrets = append(secrets, compiler.Secret{
				Name:  sec.Name,
				Value: sec.Value,
				Match: sec.Images,
			})
		}

		y := b.Yaml
		s, err := envsubst.Eval(y, func(name string) string {
			env := environ[name]
			if strings.Contains(env, "\n") {
				env = fmt.Sprintf("%q", env)
			}
			return env
		})
		if err != nil {
			return nil, err
		}
		y = s

		parsed, err := yaml.ParseString(y)
		if err != nil {
			return nil, err
		}
		metadata.Sys.Arch = parsed.Platform
		if metadata.Sys.Arch == "" {
			metadata.Sys.Arch = "linux/amd64"
		}

		//lerr := linter.New(
		//	linter.WithTrusted(b.Repo.IsTrusted),
		//).Lint(parsed)
		//if lerr != nil {
		//	return nil, lerr
		//}

		var registries []compiler.Registry
		for _, reg := range b.Regs {
			registries = append(registries, compiler.Registry{
				Hostname: reg.Address,
				Username: reg.Username,
				Password: reg.Password,
				Email:    reg.Email,
			})
		}

		ir := compiler.New(
			compiler.WithEnviron(environ),
			compiler.WithEnviron(b.Envs),
			compiler.WithEscalated(CNCDConfig.Pipeline.Privileged...),
			compiler.WithResourceLimit(CNCDConfig.Pipeline.Limits.MemSwapLimit, CNCDConfig.Pipeline.Limits.MemLimit, CNCDConfig.Pipeline.Limits.ShmSize, CNCDConfig.Pipeline.Limits.CPUQuota, CNCDConfig.Pipeline.Limits.CPUShares, CNCDConfig.Pipeline.Limits.CPUSet),
			compiler.WithVolumes(CNCDConfig.Pipeline.Volumes...),
			compiler.WithNetworks(CNCDConfig.Pipeline.Networks...),
			compiler.WithLocal(false),
			compiler.WithOption(
				compiler.WithNetrc(
					b.Netrc.Login,
					b.Netrc.Password,
					b.Netrc.Machine,
				),
				b.Repo.IsPrivate,
			),
			compiler.WithRegistry(registries...),
			compiler.WithSecret(secrets...),
			compiler.WithPrefix(
				fmt.Sprintf(
					"%d_%d",
					proc.ID,
					rand.Int(),
				),
			),
			compiler.WithEnviron(proc.Environ),
			compiler.WithProxy(),
			compiler.WithWorkspaceFromURL("/drone", b.Repo.Link()),
			compiler.WithMetadata(metadata),
		).Compile(parsed)

		// for _, sec := range b.Secs {
		// 	if !sec.MatchEvent(b.Curr.Event) {
		// 		continue
		// 	}
		// 	if b.Curr.Verified || sec.SkipVerify {
		// 		ir.Secrets = append(ir.Secrets, &backend.Secret{
		// 			Mask:  sec.Conceal,
		// 			Name:  sec.Name,
		// 			Value: sec.Value,
		// 		})
		// 	}
		// }

		item := &buildItem{
			Proc:     proc,
			Config:   ir,
			Labels:   parsed.Labels,
			Platform: metadata.Sys.Arch,
		}
		if item.Labels == nil {
			item.Labels = map[string]string{}
		}
		items = append(items, item)
	}

	return items, nil
}

func PrepareCNCD(repo *Repository, event HookEventType, p api.Payloader) {

	//remote_ := remote.FromContext(c)

	//tmprepo, build, err := remote_.Hook(c.Request)

	if event != HOOK_EVENT_PUSH {
		log.Error(0, "[PrepareCNCD] Event Type (%s) is not supported", event)
		return
	}

	var build *Build

	//check twice
	if pushPayload, ok := p.(*api.PushPayload); ok {

		build = &Build{
			Event:   EventPush,
			Commit:  pushPayload.After,
			Ref:     pushPayload.Ref,
			Link:    pushPayload.CompareURL,
			Branch:  strings.TrimPrefix(pushPayload.Ref, "refs/heads/"),
			Message: pushPayload.Commits[0].Message,
			//Avatar:    avatar,
			//Author:    author,
			//Email:     hook.Pusher.Email,
			//Timestamp: time.Now().UTC().Unix(),
			//Sender:    sender,
		}
	}

	// fetch the build file from the database
	//confb, err := remote_.File(user, repo, build, repo.Config)
	confb, err := GetFileFromGit(repo, build.Branch, ".drone.yml")

	if err != nil {
		//logrus.Errorf("error: %s: cannot find %s in %s: %s", repo.FullName, repo.Config, build.Ref, err)
		//c.AbortWithError(404, err)
		log.Error(0, "[PrepareCNCD] error: %s: cannot find .drone.yml in %s: %s ", repo.FullName(), build.Ref, err)
		return
	}
	sha := shasum(confb)
	conf, err := ConfigFind(repo, sha)
	if err != nil {
		conf = &Config{
			RepoID: repo.ID,
			Data:   string(confb),
			Hash:   sha,
		}
		err = ConfigCreate(conf)
		if err != nil {
			// retry in case we receive two hooks at the same time
			conf, err = ConfigFind(repo, sha)
			if err != nil {
				log.Error(0, "[PrepareCNCD] failure to find or persist build CNCDConfig for %s. %s ", repo.FullName(), err)
				//c.AbortWithError(500, err)
				return
			}
		}
	}

	build.ConfigID = conf.ID

	// verify the branches can be built vs skipped
	branches, err := yaml.ParseString(conf.Data)
	if err == nil {
		if !branches.Branches.Match(build.Branch) && build.Event != EventTag && build.Event != EventDeploy {
			//c.String(200, "Branch does not match restrictions defined in yaml")
			log.Error(0, "[PrepareCNCD] Branch does not match restrictions defined in yaml ", nil)
			return
		}
	}

	// update some build fields
	build.RepoID = repo.ID
	build.Verified = true
	build.Status = StatusPending

	build.Trim()
	err = CreateBuild(build, build.Procs...)
	if err != nil {
		log.Error(0, "[PrepareCNCD] failure to create build and procs %s. %s", repo.FullName(), err)
		return
	}

	if err == nil {
		go BuildQueue.Add(repo.ID)
	}

}

func shasum(raw []byte) string {
	sum := sha256.Sum256(raw)
	return fmt.Sprintf("%x", sum)
}

func (build *Build) deliver() {

	build.IsDelivered = true

	envs := map[string]string{}

	netrc := &Netrc{
		Login:    "test",
		Password: "x-oauth-basic",
		Machine:  "localhost",
	}

	repo, _ := GetRepositoryByID(build.RepoID)

	secs, err := SecretListBuild(repo)
	if err != nil {
		log.Error(0, "[Build.deliver] Error getting secrets for %s#%d. %s", repo.FullName(), build.Number, err)
	}

	regs, err := RegistryList(repo)
	if err != nil {
		log.Error(0, "[Build.deliver] Error getting registry credentials for %s#%d. %s", repo.FullName(), build.Number, err)
	}

	// get the previous build so that we can send
	// on status change notifications
	last, err := GetBuildLastBefore(repo, build.Branch, build.ID)
	if err != nil {
		log.Error(0, "[Build.deliver] Error getting last build before for %s#%d. %s", repo.FullName(), build.Number, err)
	}
	//
	// BELOW: NEW
	//

	conf, err := GetConfigByID(build.ConfigID)
	if err != nil {
		log.Error(0, "[Build.deliver] Error getting config data for %s#%d. %s", repo.FullName(), build.Number, err)
		return
	}

	b := cBuilder{
		Repo:  repo,
		Curr:  build,
		Last:  last,
		Netrc: netrc,
		Secs:  secs,
		Regs:  regs,
		Envs:  envs,
		//	Link:  httputil.GetURL(c.Request),
		Yaml: conf.Data,
	}
	items, err := b.Build()
	if err != nil {
		build.Status = StatusError
		build.Started = time.Now().Unix()
		build.Finished = build.Started
		build.Error = err.Error()
		UpdateBuild(build)
		return
	}

	var pcounter = len(items)

	for _, item := range items {
		build.Procs = append(build.Procs, item.Proc)
		item.Proc.BuildID = build.ID

		for _, stage := range item.Config.Stages {
			var gid int
			for _, step := range stage.Steps {
				pcounter++
				if gid == 0 {
					gid = pcounter
				}
				proc := &Proc{
					BuildID: build.ID,
					Name:    step.Alias,
					PID:     pcounter,
					PPID:    item.Proc.PID,
					PGID:    gid,
					State:   StatusPending,
				}
				build.Procs = append(build.Procs, proc)
			}
		}
	}
	err = ProcCreate(build.Procs)
	if err != nil {
		log.Error(0, "[Build.deliver] error persisting procs %s/%d: %s", repo.FullName(), build.Number, err)
	}

	//
	// publish topic
	//
	message := pubsub.Message{
		Labels: map[string]string{
			"repo":    repo.FullName(),
			"private": strconv.FormatBool(repo.IsPrivate),
		},
	}
	buildCopy := *build
	buildCopy.Procs = Tree(buildCopy.Procs)
	message.Data, _ = json.Marshal(Event{
		Type:  Enqueued,
		Repo:  *repo,
		Build: buildCopy,
	})
	// TODO remove global reference
	CNCDConfig.Services.Pubsub.Publish(context.Background(), "topic/events", message)

	//end publish topic

	for _, item := range items {
		task := new(queue.Task)
		task.ID = fmt.Sprint(item.Proc.ID)
		task.Labels = map[string]string{}
		task.Labels["platform"] = item.Platform
		for k, v := range item.Labels {
			task.Labels[k] = v
		}

		task.Data, _ = json.Marshal(rpc.Pipeline{
			ID:      fmt.Sprint(item.Proc.ID),
			Config:  item.Config,
			Timeout: 10000,
		})

		CNCDConfig.Services.Logs.Open(context.Background(), task.ID)
		CNCDConfig.Services.Queue.Push(context.Background(), task)
	}

}

func DeliverBuilds() {
	builds := make([]*Build, 0, 10)
	x.Where("is_delivered = ?", false).Iterate(new(Build),
		func(idx int, bean interface{}) error {
			t := bean.(*Build)
			t.deliver()
			builds = append(builds, t)
			return nil
		})

	// Update hook task status.
	for _, t := range builds {
		if err := UpdateBuild(t); err != nil {
			log.Error(4, "UpdateBuild [%d]: %v", t.ID, err)
		}
	}

	// Start listening on new hook requests.
	for repoID := range BuildQueue.Queue() {
		log.Trace("DeliverBuilds [repo_id: %v]", repoID)
		BuildQueue.Remove(repoID)

		builds = make([]*Build, 0, 5)
		if err := x.Where("repo_id = ?", repoID).And("is_delivered = ?", false).Find(&builds); err != nil {
			log.Error(4, "Get repository [%s] hook tasks: %v", repoID, err)
			continue
		}
		for _, t := range builds {
			t.deliver()
			if err := UpdateBuild(t); err != nil {
				log.Error(4, "UpdateBuild [%d]: %v", t.ID, err)
				continue
			}
		}
	}
}

func InitDeliverBuilds() {
	go DeliverBuilds()
}
