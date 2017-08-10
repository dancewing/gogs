package jenkins_client

import (
	"fmt"

	"github.com/kataras/iris/core/errors"
)

//"active" : true,
//"backupVersion" : null,
//"bundled" : false,
//"deleted" : false,
//"dependencies" : [
//{
//
//}
//],
//"downgradable" : false,
//"enabled" : true,
//"hasUpdate" : false,
//"longName" : "JavaScript GUI Lib: ACE Editor bundle plugin",
//"pinned" : false,
//"requiredCoreVersion" : "1.580.1",
//"shortName" : "ace-editor",
//"supportsDynamicLoad" : "YES",
//"url" : "https://wiki.jenkins-ci.org/display/JENKINS/ACE+Editor+Plugin",
//"version" : "1.1"

type Plugin struct {
	Active              bool   `json:"active"`
	BackupVersion       string `json:"backupVersion"`
	Bundled             bool   `json:"bundled"`
	Deleted             bool   `json:"deleted"`
	DownGradable        bool   `json:"downgradable"`
	Enabled             bool   `json:"enabled"`
	HasUpdate           bool   `json:"hasUpdate"`
	LongName            string `json:"longName"`
	Pinned              bool   `json:"pinned"`
	RequiredCoreVersion string `json:"requiredCoreVersion"`
	ShortName           string `json:"shortName"`
	SupportsDynamicLoad string `json:"supportsDynamicLoad"`
	Url                 string   `json:"url"`
	Version             string `json:"version"`
}

type Plugins struct {
	Plugins []Plugin `json:"plugins"`
}

func (plugins Plugins) GetVersion(shortName string) (version string, err error) {
	if len(plugins.Plugins) == 0 {
		return "", errors.New("No Plugins loaded")
	}

	for _, p := range plugins.Plugins {
		if p.ShortName == shortName {
			return p.Version, nil
		}
	}
	return "", errors.New(fmt.Sprintf("No Plugin found with name (%s)", shortName))
}
