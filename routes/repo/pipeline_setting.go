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

func PipelineHooksEditPost(c *context.Context, f form.NewPipelineHook) {
	c.Data["Title"] = c.Tr("repo.settings.add_webhook")
	c.Data["PageIsSettingsHooks"] = true
	c.Data["PageIsSettingsHooksNew"] = true
	c.Data["PipelineHook"] = models.PipelineHook{HookEvent: &models.HookEvent{}}
	c.Data["HookType"] = "jenkins"

	//orCtx, err := getOrgRepoCtx(c)
	orCtx, w := checkPipelineHook(c)
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
		w = &models.PipelineHook{
			RepoID:      orCtx.RepoID,
			ContentType: contentType,
			Secret:      f.Secret,
			HookEvent:   ParseHookEvent(f.Webhook),
			IsActive:    f.Active,
			OrgID:       orCtx.OrgID,
		}

		if err := w.UpdateEvent(); err != nil {
			c.Handle(500, "UpdateEvent", err)
			return
		} else if err := models.CreatePipelineHook(w); err != nil {
			c.Handle(500, "CreatePipelineHook", err)
			return
		}

	} else {
		w.RepoID = orCtx.RepoID
		w.ContentType = contentType
		w.Secret = f.Secret
		w.HookEvent = ParseHookEvent(f.Webhook)
		w.IsActive = f.Active
		w.OrgID = orCtx.OrgID

		if err := w.UpdateEvent(); err != nil {
			c.Handle(500, "UpdateEvent", err)
			return
		} else if err := models.UpdatePipelineHook(w); err != nil {
			c.Handle(500, "CreatePipelineHook", err)
			return
		}
	}

	c.Flash.Success(c.Tr("repo.settings.add_hook_success"))
	c.Redirect(orCtx.Link + "/settings/hooks/pipeline")
}

func PipelineHooksEdit(c *context.Context) {
	c.Data["Title"] = c.Tr("repo.settings.update_webhook")
	c.Data["PageIsSettingsHooks"] = true
	c.Data["PageIsSettingsHooksNew"] = true

	orCtx, w := checkPipelineHook(c)
	if c.Written() {
		return
	}

	if w == nil {
		c.Data["HookFound"] = false
		c.Data["PipelineHook"] = models.PipelineHook{HookEvent: &models.HookEvent{}}
	} else {
		c.Data["HookFound"] = true
		c.Data["PipelineHook"] = w
	}

	c.HTML(200, orCtx.NewTemplate)
}

func checkPipelineHook(c *context.Context) (*OrgRepoCtx, *models.PipelineHook) {
	c.Data["RequireHighlightJS"] = true

	orCtx, err := getPipelineOrgRepoCtx(c)
	if err != nil {
		c.Handle(500, "getOrgRepoCtx", err)
		return nil, nil
	}
	c.Data["BaseLink"] = orCtx.Link

	var w *models.PipelineHook

	if orCtx.RepoID > 0 {
		w, err = models.GetPipelineHookOfRepoByID(c.Repo.Repository.ID, c.ParamsInt64(":id"))
	} else if orCtx.OrgID > 0 {
		w, err = models.GetPipelineHookByOrgID(c.Org.Organization.ID, c.ParamsInt64(":id"))
	}
	if err != nil {
		//c.NotFoundOrServerError("GetPipelineHookOfRepoByID/GetPipelineHookByOrgID", errors.IsWebhookNotExist, err)
		//w = &models.PipelineHook{HookEvent: &models.HookEvent{}}
		return orCtx, nil
	}

	c.Data["History"], err = w.History(1)
	if err != nil {
		c.Handle(500, "History", err)
	}
	return orCtx, w
}

func getPipelineOrgRepoCtx(c *context.Context) (*OrgRepoCtx, error) {
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
