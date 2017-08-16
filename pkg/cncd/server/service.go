package server

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/frontend"
	"github.com/cncd/pipeline/pipeline/frontend/yaml"
	"github.com/cncd/pipeline/pipeline/frontend/yaml/compiler"
	"github.com/cncd/pipeline/pipeline/frontend/yaml/matrix"
	"github.com/drone/envsubst"

	"github.com/gogits/gogs/models"
)

// return the metadata from the cli context.
func metadataFromStruct(repo *models.Repository, build, last *models.Build, proc *models.Proc, link string) frontend.Metadata {
	return frontend.Metadata{
		Repo: frontend.Repo{
			Name:    repo.FullName(),
			Link:    repo.Link(),
			Remote:  repo.FullAddress(),
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

type builder struct {
	Repo  *models.Repository
	Curr  *models.Build
	Last  *models.Build
	Netrc *models.Netrc
	Secs  []*models.Secret
	Regs  []*models.Registry
	Link  string
	Yaml  string
	Envs  map[string]string
}

type buildItem struct {
	Proc     *models.Proc
	Platform string
	Labels   map[string]string
	Config   *backend.Config
}

func (b *builder) Build() ([]*buildItem, error) {

	axes, err := matrix.ParseString(b.Yaml)
	if err != nil {
		return nil, err
	}
	if len(axes) == 0 {
		axes = append(axes, matrix.Axis{})
	}

	var items []*buildItem
	for i, axis := range axes {
		proc := &models.Proc{
			BuildID: b.Curr.ID,
			PID:     i + 1,
			PGID:    i + 1,
			State:   models.StatusPending,
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
			compiler.WithEscalated(Config.Pipeline.Privileged...),
			compiler.WithResourceLimit(Config.Pipeline.Limits.MemSwapLimit, Config.Pipeline.Limits.MemLimit, Config.Pipeline.Limits.ShmSize, Config.Pipeline.Limits.CPUQuota, Config.Pipeline.Limits.CPUShares, Config.Pipeline.Limits.CPUSet),
			compiler.WithVolumes(Config.Pipeline.Volumes...),
			compiler.WithNetworks(Config.Pipeline.Networks...),
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
