package repo

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
)

const (
	SETTINGS_SECRET_LIST = "repo/settings/cncd/secret"
)

func GetSecretList(c *context.RestContext) {
	c.Data["Title"] = c.Tr("repo.settings.githooks")
	c.Data["PageIsSettingsCNCD"] = true

	secrets, err := models.SecretListBuild(c.Repo.Repository)

	if err != nil {
		c.Handle(500, "Hooks", err)
		return
	}

	c.Data["Secrets"] = secrets

	c.HTML(200, SETTINGS_SECRET_LIST)
}
