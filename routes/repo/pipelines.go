package repo

import (
	"github.com/Unknwon/paginater"
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
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
)

func RetrievePipeline(c *context.Context) {

	c.Data["PageIsPipelines"] = true

	jobDef, err := models.GetCIFileFromGit(c.Repo.Owner, c.Repo.Repository)
	if err != nil {
		c.Handle(200, "RetrievePipeline.GetCIFileFromGit", err)
		return
	}

	if jobDef == nil {
		c.HTML(200, PIPELINE_NONE)
		return
	}

	c.Data["JobDef"] = jobDef
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

func NewPipelinePost(c *context.Context) {

	c.Data["PageIsPipelineList"] = true

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
