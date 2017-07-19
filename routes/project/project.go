package project

import (
	"github.com/gogits/gogs/pkg/context"

	"github.com/Unknwon/paginater"
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/form"
	"github.com/gogits/gogs/pkg/setting"
	log "gopkg.in/clog.v1"
)

const (
	CREATE = "project/create"
	LIST   = "project/projects"
)

//Create project, get
func Create(c *context.Context) {
	c.Data["Title"] = c.Tr("new_project")

	c.User.GetRepositories(0, 20)
	c.User.GetTeams()
	c.Data["Repos"] = c.User.Repos
	c.Data["Teams"] = c.User.Teams

	c.Data["ContextUser"] = c.User

	c.HTML(200, CREATE)
}

//Create projects
func CreatePost(c *context.Context, f form.CreateProject) {
	c.Data["Title"] = c.Tr("new_repo")

	c.User.GetRepositories(0, 20)
	c.User.GetTeams()
	c.Data["Repos"] = c.User.Repos
	c.Data["Teams"] = c.User.Teams

	ctxUser := c.User
	if c.Written() {
		return
	}
	c.Data["ContextUser"] = ctxUser

	if c.HasError() {
		c.HTML(200, CREATE)
		return
	}

	prg, err := models.CreateProject(ctxUser, models.CreateProjectOptions{
		Name:        f.ProjectName,
		Description: f.Description,
	})

	if err == nil {
		log.Trace("Project created [%d]: %s/%s", prg.ID, ctxUser.Name, prg.Name)
		c.Redirect(setting.AppSubURL + "/" + ctxUser.Name + "/" + prg.Name)
		return
	}

	handleCreateError(c, ctxUser, err, "CreatePost", CREATE, &f)
}

func handleCreateError(c *context.Context, owner *models.User, err error, name, tpl string, form interface{}) {
	switch {
	case errors.IsReachLimitOfRepo(err):
		c.RenderWithErr(c.Tr("repo.form.reach_limit_of_creation", owner.RepoCreationNum()), tpl, form)
	case models.IsErrRepoAlreadyExist(err):
		c.Data["Err_RepoName"] = true
		c.RenderWithErr(c.Tr("form.repo_name_been_taken"), tpl, form)
	case models.IsErrNameReserved(err):
		c.Data["Err_RepoName"] = true
		c.RenderWithErr(c.Tr("repo.form.name_reserved", err.(models.ErrNameReserved).Name), tpl, form)
	case models.IsErrNamePatternNotAllowed(err):
		c.Data["Err_RepoName"] = true
		c.RenderWithErr(c.Tr("repo.form.name_pattern_not_allowed", err.(models.ErrNamePatternNotAllowed).Pattern), tpl, form)
	default:
		c.Handle(500, name, err)
	}
}

func ListProject(c *context.Context) {
	c.Data["Title"] = c.Tr("projects")

	//c.Data["PageIsExplore"] = true
	//c.Data["PageIsExploreRepositories"] = true

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	keyword := c.Query("q")
	searchOpt := &models.SearchProjectOptions{
		Keyword:  keyword,
		UserID:   -1,
		OrderBy:  "updated_unix DESC",
		Page:     page,
		PageSize: setting.UI.ExplorePagingNum,
	}

	prgs, count, err := models.SearchProjectByName(searchOpt)

	if err != nil {
		c.Handle(500, "SearchProjectByName", err)
		return
	}
	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), setting.UI.ExplorePagingNum, page, 5)

	c.Data["Projects"] = prgs

	c.HTML(200, LIST)
}

func ListMyProject(c *context.Context) {

	c.Data["Title"] = c.Tr("my_projects")

	//c.Data["PageIsExplore"] = true
	//c.Data["PageIsExploreRepositories"] = true

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	keyword := c.Query("q")
	searchOpt := &models.SearchProjectOptions{
		Keyword:  keyword,
		UserID:   -1,
		OrderBy:  "updated_unix DESC",
		Page:     page,
		PageSize: setting.UI.ExplorePagingNum,
	}

	if c.IsLogged {
		searchOpt.UserID = c.User.ID
	}

	prgs, count, err := models.SearchProjectByName(searchOpt)

	if err != nil {
		c.Handle(500, "SearchProjectByName", err)
		return
	}
	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), setting.UI.ExplorePagingNum, page, 5)

	c.Data["Projects"] = prgs

	c.HTML(200, LIST)
}
