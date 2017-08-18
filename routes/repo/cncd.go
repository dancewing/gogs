package repo

import (
	"github.com/Unknwon/paginater"
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
)

const (
	BUILD_LIST = "repo/cncd/build_list"
	BUILD_VIEW = "repo/cncd/build_view"
)

func GetBuildList(c *context.Context) {

	c.Data["PageIsCNCD"] = true
	c.Data["PageIsCNCDBuildList"] = true

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	builds, err := models.ListBuilds(&models.BuildOptions{
		RepositoryID: c.Repo.Repository.ID,
		Page:         page,
		PageSize:     10,
	})

	if err != nil {
		c.Handle(500, "GetBuildList", err)
		return
	}

	//for _, pipeline := range builds {
	//	err = pipeline.LoadAttributes()
	//
	//	if err != nil {
	//		c.Handle(500, "ListPipelines.LoadAttributes", err)
	//		return
	//	}
	//}

	c.Data["Builds"] = builds

	count := models.CountBuild(c.Repo.Repository.ID)

	c.Data["Page"] = paginater.New(int(count), 10, page, 5)

	c.HTML(200, BUILD_LIST)
}

func GetBuild(c *context.Context) {

	c.Data["PageIsCNCD"] = true
	c.Data["PageIsCNCDBuildList"] = true

	buildID := c.ParamsInt64(":id")

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	build, err := models.GetBuild(buildID)

	if err != nil {
		c.Handle(500, "GetBuild", err)
		return
	}

	c.Data["Build"] = build

	builds, err := models.ProcList(build)

	if err != nil {
		c.Handle(500, "ProcList", err)
		return
	}

	//for _, pipeline := range builds {
	//	err = pipeline.LoadAttributes()
	//
	//	if err != nil {
	//		c.Handle(500, "ListPipelines.LoadAttributes", err)
	//		return
	//	}
	//}

	c.Data["Procs"] = builds

	c.HTML(200, BUILD_VIEW)
}
