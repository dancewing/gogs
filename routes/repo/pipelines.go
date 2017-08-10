package repo

import (
	"bytes"
	"fmt"
	"strings"

	gotemplate "html/template"

	"github.com/Unknwon/paginater"
	"github.com/gogits/git-module"
	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
)

const (
	PIPELINE_CREATE   = "repo/pipeline/new_pipeline"
	PIPELINE_VIEW     = "repo/pipeline/view_pipeline"
	PIPELINES         = "repo/pipeline/list"
	PIPELINE_JOBS     = "repo/pipeline/jobs"
	PIPELINE_JOB_VIEW = "repo/pipeline/view_job"
	PIPELINE_ENV      = "repo/pipeline/environments"
	PIPELINE_ENV_VIEW = "repo/pipeline/view_environment"
	PIPELINE_NONE     = "repo/pipeline/no_pipeline"

	PIPELINE_PREVIEW = "repo/pipeline/preview_pipeline"
)

func RetrievePipeline(c *context.Context) {

	c.Data["PageIsPipelines"] = true
	//
	//jobDef, err := models.GetCIFileFromGit(c.Repo.Owner, c.Repo.Repository)
	//if err != nil {
	//	c.Handle(200, "RetrievePipeline.GetCIFileFromGit", err)
	//	return
	//}
	//
	//if jobDef == nil {
	//	c.HTML(200, PIPELINE_NONE)
	//	return
	//}

	//c.Data["JobDef"] = jobDef
	c.Data["PageIsPipelines"] = true
}

func ListPipelines(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	pipelines, err := models.ListPipelines(&models.PipelineOptions{
		RepositoryID: c.Repo.Repository.ID,
		Page:         page,
		PageSize:     10,
	})

	if err != nil {
		c.Handle(500, "ListPipelines", err)
		return
	}

	for _, pipeline := range pipelines {
		err = pipeline.LoadAttributes()

		if err != nil {
			c.Handle(500, "ListPipelines.LoadAttributes", err)
			return
		}
	}

	c.Data["Pipelines"] = pipelines

	count := models.CountPipeline(c.Repo.Repository.ID)

	c.Data["Page"] = paginater.New(int(count), 10, page, 5)

	c.HTML(200, PIPELINES)

}

func ViewPipeline(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

	id := c.ParamsInt64(":id")

	pipeline, err := models.GetPipeline(id)

	if err != nil {
		c.Handle(500, "ViewPipeline.GetPipeline", err)
		return
	}

	err = pipeline.LoadAttributes()

	if err != nil {
		c.Handle(500, "ViewPipeline.LoadAttributes", err)
		return
	}

	c.Data["Pipeline"] = pipeline

	c.HTML(200, PIPELINE_VIEW)

}

func NewPipeline(c *context.Context) {
	c.Data["PageIsPipelineList"] = true

	c.HTML(200, PIPELINE_CREATE)

}

func NewPipelinePost(c *context.Context, form form.NewPipeline) {

	c.Data["PageIsPipelineList"] = true

	if form.Status == "preview" {

		content, err := models.PreviewPipelineScript(c.Repo.Repository, form.Branch)

		if err != nil {
			c.Handle(500, "PreviewPipelineScript", err)

			return
		}

		//c.Data["PreviewContent"] = content
		//
		//var fileContent string
		//if err, content := template.ToUTF8WithErr([]byte(content)); err != nil {
		//	if err != nil {
		//		log.Error(4, "ToUTF8WithErr: %s", err)
		//	}
		//	fileContent = string(buf)
		//} else {
		//	fileContent = content
		//}
		fileContent := content

		var output bytes.Buffer
		lines := strings.Split(fileContent, "\n")
		for index, line := range lines {
			output.WriteString(fmt.Sprintf(`<li class="L%d" rel="L%d">%s</li>`, index+1, index+1, gotemplate.HTMLEscapeString(strings.TrimRight(line, "\r"))) + "\n")
		}
		c.Data["FileContent"] = gotemplate.HTML(output.String())

		output.Reset()
		for i := 0; i < len(lines); i++ {
			output.WriteString(fmt.Sprintf(`<span id="L%d">%d</span>`, i+1, i+1))
		}
		c.Data["LineNums"] = gotemplate.HTML(output.String())

		c.HTML(200, PIPELINE_PREVIEW)

		return
	} else {
		err := models.RunPipeline(c.Repo.Repository, form.Branch, getFakePayroad(c, form.Branch))

		if err != nil {
			c.Handle(500, "PreviewPipelineScript", err)

			return
		}
	}

	c.HTML(200, PIPELINES)

}

func ListJobs(c *context.Context) {
	c.Data["PageIsJobList"] = true

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	jobs, err := models.ListJobs(&models.JobOptions{
		RepositoryID: c.Repo.Repository.ID,
		Page:         page,
		PageSize:     10,
	})

	if err != nil {
		c.Handle(500, "ListPipelines", err)
		return
	}

	//for _, job := range jobs {
	//	//		err = job.LoadAttributes()
	//
	//	if err != nil {
	//		c.Handle(500, "ListPipelines.LoadAttributes", err)
	//		return
	//	}
	//}

	c.Data["Jobs"] = jobs

	count := models.CountJob(c.Repo.Repository.ID)

	c.Data["Page"] = paginater.New(int(count), 10, page, 5)

	c.HTML(200, PIPELINE_JOBS)
}

func ViewJob(c *context.Context) {
	c.Data["PageIsJobList"] = true

	id := c.ParamsInt64(":id")

	job, err := models.GetJob(id)

	if err != nil {
		c.Handle(500, "ViewJob.GetPipeline", err)
		return
	}

	if err != nil {
		c.Handle(500, "ViewPipeline.LoadAttributes", err)
		return
	}

	c.Data["Job"] = job

	c.HTML(200, PIPELINE_JOB_VIEW)
}

func ListEnvironments(c *context.Context) {

	c.Data["PageIsEnvironmentList"] = true

	c.HTML(200, PIPELINE_ENV)
}

func ViewEnvironment(c *context.Context) {

	c.Data["PageIsEnvironmentList"] = true

	c.HTML(200, PIPELINE_ENV_VIEW)
}

func getFakePayroad(c *context.Context, branch string) api.Payloader {
	var authorUsername, committerUsername string

	// Grab latest commit or fake one if it's empty repository.
	// commit := c.Repo.Commit
	commit, err :=  models.GetLastCommit(c.Repo.Repository, branch)

	if commit == nil {
		ghost := models.NewGhostUser()
		commit = &git.Commit{
			ID:            git.MustIDFromString(git.EMPTY_SHA),
			Author:        ghost.NewGitSig(),
			Committer:     ghost.NewGitSig(),
			CommitMessage: "This is a fake commit",
		}
		authorUsername = ghost.Name
		committerUsername = ghost.Name
	} else {
		// Try to match email with a real user.
		author, err := models.GetUserByEmail(commit.Author.Email)
		if err == nil {
			authorUsername = author.Name
		} else if !errors.IsUserNotExist(err) {
			c.Handle(500, "GetUserByEmail.(author)", err)
			return nil
		}

		committer, err := models.GetUserByEmail(commit.Committer.Email)
		if err == nil {
			committerUsername = committer.Name
		} else if !errors.IsUserNotExist(err) {
			c.Handle(500, "GetUserByEmail.(committer)", err)
			return nil
		}
	}



	fileStatus, err := commit.FileStatus()
	if err != nil {
		c.Handle(500, "FileStatus", err)
		return nil
	}

	apiUser := c.User.APIFormat()
	p := &api.PushPayload{
		Ref:    git.BRANCH_PREFIX + c.Repo.Repository.DefaultBranch,
		Before: commit.ID.String(),
		After:  commit.ID.String(),
		Commits: []*api.PayloadCommit{
			{
				ID:      commit.ID.String(),
				Message: commit.Message(),
				URL:     c.Repo.Repository.HTMLURL() + "/commit/" + commit.ID.String(),
				Author: &api.PayloadUser{
					Name:     commit.Author.Name,
					Email:    commit.Author.Email,
					UserName: authorUsername,
				},
				Committer: &api.PayloadUser{
					Name:     commit.Committer.Name,
					Email:    commit.Committer.Email,
					UserName: committerUsername,
				},
				Added:    fileStatus.Added,
				Removed:  fileStatus.Removed,
				Modified: fileStatus.Modified,
			},
		},
		Repo:   c.Repo.Repository.APIFormat(nil),
		Pusher: apiUser,
		Sender: apiUser,
	}
	return p
}
