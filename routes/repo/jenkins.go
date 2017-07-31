package repo

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
)

const (
	PIPELINE_HOOK_EDIT = "repo/settings/pipeline_hook_edit"
)

func JenkinsHooksEditPost(c *context.Context, f form.NewWebhook) {
	c.Data["Title"] = c.Tr("repo.settings.add_webhook")
	c.Data["PageIsSettingsHooks"] = true
	c.Data["PageIsSettingsHooksNew"] = true
	c.Data["Webhook"] = models.Webhook{HookEvent: &models.HookEvent{}}
	c.Data["HookType"] = "jenkins"

	//orCtx, err := getOrgRepoCtx(c)
	orCtx, w := checkJenkinsHook(c)
	//if err != nil {
	//	c.Handle(500, "getOrgRepoCtx", err)
	//	return
	//}
	c.Data["BaseLink"] = orCtx.Link

	if c.HasError() {
		c.HTML(200, orCtx.NewTemplate)
		return
	}

	contentType := models.JSON
	if models.HookContentType(f.ContentType) == models.FORM {
		contentType = models.FORM
	}

	if w == nil {
		w = &models.JenkinsHook{
			RepoID:      orCtx.RepoID,
			URL:         f.PayloadURL,
			ContentType: contentType,
			Secret:      f.Secret,
			HookEvent:   ParseHookEvent(f.Webhook),
			IsActive:    f.Active,
			OrgID:       orCtx.OrgID,
		}

		if err := w.UpdateEvent(); err != nil {
			c.Handle(500, "UpdateEvent", err)
			return
		} else if err := models.CreateJenkinsHook(w); err != nil {
			c.Handle(500, "CreateJenkinsHook", err)
			return
		}

	} else {
		w.RepoID = orCtx.RepoID
		w.URL = f.PayloadURL
		w.ContentType = contentType
		w.Secret = f.Secret
		w.HookEvent = ParseHookEvent(f.Webhook)
		w.IsActive = f.Active
		w.OrgID = orCtx.OrgID

		if err := w.UpdateEvent(); err != nil {
			c.Handle(500, "UpdateEvent", err)
			return
		} else if err := models.UpdateJenkinsHook(w); err != nil {
			c.Handle(500, "CreateJenkinsHook", err)
			return
		}
	}

	c.Flash.Success(c.Tr("repo.settings.add_hook_success"))
	c.Redirect(orCtx.Link + "/settings/hooks/pipeline")
}

func JenkinsHooksEdit(c *context.Context) {
	c.Data["Title"] = c.Tr("repo.settings.update_webhook")
	c.Data["PageIsSettingsHooks"] = true
	c.Data["PageIsSettingsHooksNew"] = true

	orCtx, w := checkJenkinsHook(c)
	if c.Written() {
		return
	}

	if w == nil {
		c.Data["HookFound"] = false
		c.Data["Webhook"] = models.JenkinsHook{HookEvent: &models.HookEvent{}}
	} else {
		c.Data["HookFound"] = true
		c.Data["Webhook"] = w
	}

	c.HTML(200, orCtx.NewTemplate)
}

func checkJenkinsHook(c *context.Context) (*OrgRepoCtx, *models.JenkinsHook) {
	c.Data["RequireHighlightJS"] = true

	orCtx, err := getJenkinsOrgRepoCtx(c)
	if err != nil {
		c.Handle(500, "getOrgRepoCtx", err)
		return nil, nil
	}
	c.Data["BaseLink"] = orCtx.Link

	var w *models.JenkinsHook

	if orCtx.RepoID > 0 {
		w, err = models.GetJenkinsHookOfRepoByID(c.Repo.Repository.ID, c.ParamsInt64(":id"))
	} else if orCtx.OrgID > 0 {
		w, err = models.GetJenkinsHookByOrgID(c.Org.Organization.ID, c.ParamsInt64(":id"))
	}
	if err != nil {
		//c.NotFoundOrServerError("GetJenkinsHookOfRepoByID/GetJenkinsHookByOrgID", errors.IsWebhookNotExist, err)
		//w = &models.JenkinsHook{HookEvent: &models.HookEvent{}}
		return orCtx, nil
	}

	c.Data["History"], err = w.History(1)
	if err != nil {
		c.Handle(500, "History", err)
	}
	return orCtx, w
}

func getJenkinsOrgRepoCtx(c *context.Context) (*OrgRepoCtx, error) {
	if len(c.Repo.RepoLink) > 0 {
		c.Data["PageIsRepositoryContext"] = true
		return &OrgRepoCtx{
			RepoID:      c.Repo.Repository.ID,
			Link:        c.Repo.RepoLink,
			NewTemplate: PIPELINE_HOOK_EDIT,
		}, nil
	}

	if len(c.Org.OrgLink) > 0 {
		c.Data["PageIsOrganizationContext"] = true
		return &OrgRepoCtx{
			OrgID:       c.Org.Organization.ID,
			Link:        c.Org.OrgLink,
			NewTemplate: PIPELINE_HOOK_EDIT,
		}, nil
	}

	return nil, errors.New("Unable to set OrgRepo context")
}
