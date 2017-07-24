package project

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/setting"

	"fmt"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/form"

	log "gopkg.in/clog.v1"
)

const (
	CREATE     = "project/create"
	LIST       = "project/projects"
	INITIALIZE = "project/initialize"
	MIGRATE    = "project/migrate"
)

func checkContextUser(c *context.Context, uid int64) *models.User {
	orgs, err := models.GetOwnedOrgsByUserIDDesc(c.User.ID, "updated_unix")
	if err != nil {
		c.Handle(500, "GetOwnedOrgsByUserIDDesc", err)
		return nil
	}
	c.Data["Orgs"] = orgs

	// Not equal means current user is an organization.
	if uid == c.User.ID || uid == 0 {
		return c.User
	}

	org, err := models.GetUserByID(uid)
	if errors.IsUserNotExist(err) {
		return c.User
	}

	if err != nil {
		c.Handle(500, "GetUserByID", fmt.Errorf("[%d]: %v", uid, err))
		return nil
	}

	// Check ownership of organization.
	if !org.IsOrganization() || !(c.User.IsAdmin || org.IsOwnedBy(c.User.ID)) {
		c.Error(403)
		return nil
	}
	return org
}

func Create(c *context.Context) {
	c.Data["Title"] = c.Tr("new_repo")

	// Give default value for template to render.
	c.Data["Gitignores"] = models.Gitignores
	c.Data["Licenses"] = models.Licenses
	c.Data["Readmes"] = models.Readmes
	c.Data["readme"] = "Default"
	c.Data["private"] = c.User.LastRepoVisibility
	c.Data["IsForcedPrivate"] = setting.Repository.ForcePrivate

	ctxUser := checkContextUser(c, c.QueryInt64("org"))
	if c.Written() {
		return
	}
	c.Data["ContextUser"] = ctxUser

	projects, _ := models.ListUserTopProjects(c.User)

	c.Data["Projects"] = projects

	c.HTML(200, CREATE)
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

func CreatePost(c *context.Context, f form.CreateProject) {
	c.Data["Title"] = c.Tr("new_project")

	c.Data["Gitignores"] = models.Gitignores
	c.Data["Licenses"] = models.Licenses
	c.Data["Readmes"] = models.Readmes

	ctxUser := checkContextUser(c, f.UserID)
	if c.Written() {
		return
	}
	c.Data["ContextUser"] = ctxUser

	projects, _ := models.ListUserTopProjects(c.User)

	c.Data["Projects"] = projects

	if c.HasError() {
		c.HTML(200, CREATE)
		return
	}

	repo, err := models.CreateProject(c.User, ctxUser, models.CreateProjectOptions{
		Name:          f.Name,
		Description:   f.Description,
		Gitignores:    f.Gitignores,
		License:       f.License,
		Readme:        f.Readme,
		IsPrivate:     f.Private || setting.Repository.ForcePrivate,
		AutoInit:      f.AutoInit,
		CreateGitRepo: f.CreateGitRepo,
	})
	if err == nil {
		log.Trace("Repository created [%d]: %s/%s", repo.ID, ctxUser.Name, repo.Name)
		c.Redirect(setting.AppSubURL + "/" + ctxUser.Name + "/" + repo.Name)
		return
	} else {
		log.Error(4, "CreateProject: %v", err)
	}

	if repo != nil && f.CreateGitRepo {
		if errDelete := models.DeleteRepository(ctxUser.ID, repo.ID); errDelete != nil {
			log.Error(4, "DeleteRepository: %v", errDelete)
		}
	}

	handleCreateError(c, ctxUser, err, "CreatePost", CREATE, &f)
}

func InitializeGit(c *context.Context) {

	c.Data["Gitignores"] = models.Gitignores
	c.Data["Licenses"] = models.Licenses
	c.Data["Readmes"] = models.Readmes
	c.Data["readme"] = "Default"

	c.HTML(200, INITIALIZE)
}
func InitializeGitPost(c *context.Context, f form.InitializeGit) {
	//TODO check
	c.Data["Gitignores"] = models.Gitignores
	c.Data["Licenses"] = models.Licenses
	c.Data["Readmes"] = models.Readmes
	c.Data["readme"] = "Default"

	c.Repo.Repository.LoadAttributes()

	err := c.Repo.Repository.InitializeGit(c.User, models.CreateProjectOptions{
		Gitignores: f.Gitignores,
		License:    f.License,
		Readme:     f.Readme,
		AutoInit:   f.AutoInit,
	})

	if err == nil {
		log.Trace("Repository Initialized [%d]: %s/%s", c.Repo.Repository.ID, c.Repo.Repository.Owner.Name, c.Repo.Repository.Name)
		c.Redirect(setting.AppSubURL + "/" + c.Repo.Repository.Owner.Name + "/" + c.Repo.Repository.Name)
		return
	} else {
		log.Error(4, "CreateProject: %v", err)

		c.RenderWithErr(c.Tr("form.repo_name_been_taken"), INITIALIZE, &f)
	}

}

func Migrate() {

}
func MigratePost() {

}
