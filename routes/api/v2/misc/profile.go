package misc

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/setting"
)

func GetProfileInfo(c *context.APIContext) {

	data := &ProfileInfoVM{}

	if setting.ProdMode {
		data.ActiveProfiles = append(data.ActiveProfiles, "prod")
	} else {
		data.ActiveProfiles = append(data.ActiveProfiles, "dev")
	}

	c.JSONSuccess(data)
}

type ProfileInfoVM struct {
	ActiveProfiles []string `json:"activeProfiles"       `
	RibbonEnv      string   `json:"ribbonEnv"            `
}
