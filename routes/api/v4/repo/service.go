package repo

import (
	"strings"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/context"
)

const (
	SERVICES         = "repo/settings/service/base"
	SERVICE_EDIT     = "repo/settings/service/edit"
	ORG_SERVICE_EDIT = "org/settings/service_edit"
)

func Services(c *context.RestContext) {
	c.Data["Title"] = c.Tr("repo.settings.hooks")
	c.Data["PageIsSettingsServices"] = true
	c.Data["BaseLink"] = c.Repo.RepoLink

	ws, err := models.GetAllServicesByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Handle(500, "GetAllServicesByRepoID", err)
		return
	}
	c.Data["Services"] = ws

	c.HTML(200, SERVICES)
}

func ServicesEdit(c *context.RestContext) {
	c.Data["Title"] = c.Tr("repo.settings.add_webhook")
	c.Data["PageIsSettingsServices"] = true
	c.Data["PageIsSettingsServicesEdit"] = true

	serviceTypeName := strings.ToLower(c.Params(":type"))
	serviceType := models.ToServiceType(serviceTypeName)
	config, err := models.GetConfiguration(serviceType, c.Repo.Repository.ID)
	//Can't get configuration, maybe it's not saved
	if err != nil {
		c.Data["PageIsSettingsServicesIsNew"] = true
	} else {
		c.Data["History"], err = config.History(1)
	}
	switch serviceType {
	case models.JENKINS:
		c.Data["ServiceConfig"] = models.ToJenkinsServiceConfigLoad(config)
	}

	orCtx, err := getOrgRepoCtxForService(c)
	if err != nil {
		c.Handle(500, "getOrgRepoCtxForService", err)
		return
	}

	c.Data["ServiceType"] = serviceTypeName
	if c.Written() {
		return
	}
	c.Data["BaseLink"] = orCtx.Link

	c.HTML(200, orCtx.NewTemplate)
}

// getOrgRepoCtx determines whether this is a repo context or organization context.
func getOrgRepoCtxForService(c *context.RestContext) (*OrgRepoCtx, error) {
	if len(c.Repo.RepoLink) > 0 {
		c.Data["PageIsRepositoryContext"] = true
		return &OrgRepoCtx{
			RepoID:      c.Repo.Repository.ID,
			Link:        c.Repo.RepoLink,
			NewTemplate: SERVICE_EDIT,
		}, nil
	}

	if len(c.Org.OrgLink) > 0 {
		c.Data["PageIsOrganizationContext"] = true
		return &OrgRepoCtx{
			OrgID:       c.Org.Organization.ID,
			Link:        c.Org.OrgLink,
			NewTemplate: ORG_SERVICE_EDIT,
		}, nil
	}

	return nil, errors.New("Unable to set OrgRepo context")
}

func ServicesJenkinsPost(c *context.RestContext, edit models.JenkinsServiceConfigLoad) {

	orCtx, err := getOrgRepoCtxForService(c)
	if err != nil {
		c.Handle(500, "getOrgRepoCtxForService", err)
		return
	}

	c.Data["BaseLink"] = orCtx.Link

	serviceTypeName := "jenkins"
	serviceType := models.ToServiceType(serviceTypeName)

	c.Data["ServiceType"] = serviceTypeName

	config, err := models.GetConfiguration(serviceType, c.Repo.Repository.ID)

	edit.UpdateEvent()

	if err != nil {
		c.Data["PageIsSettingsServicesIsNew"] = true
		config = &models.ServiceConfig{
			Type:     serviceType,
			RepoID:   orCtx.RepoID,
			OrgID:    orCtx.OrgID,
			IsActive: edit.IsActive,
		}
		config, err = models.CreateConfiguration(config, edit)
	} else {

		config.IsActive = edit.IsActive
		config, err = models.UpdateConfiguration(config, edit)

		c.Data["History"], err = config.History(1)
	}

	c.Data["ServiceConfig"] = models.ToJenkinsServiceConfigLoad(config)

}

func ServicesPost(c *context.RestContext) {
	c.Data["Title"] = c.Tr("repo.settings.add_webhook")
	c.Data["PageIsSettingsServices"] = true
	c.Data["PageIsSettingsServicesEdit"] = true

	orCtx, err := getOrgRepoCtxForService(c)
	if err != nil {
		c.Handle(500, "getOrgRepoCtxForService", err)
		return
	}

	c.Data["BaseLink"] = orCtx.Link

	if c.Written() {
		return
	}

	c.HTML(200, orCtx.NewTemplate)
}
