package models

import "encoding/json"

type JenkinsServiceConfigLoad struct {
	JenkinsHost  string `json:"jenkins_host"`
	JenkinsUser  string `json:"jenkins_user"`
	JenkinsToken string `json:"jenkins_token"`
	*ServiceConfigLoad
}

func (config *ServiceConfig) ToJenkinsServiceConfigEdit() *JenkinsServiceConfigLoad {
	return ToJenkinsServiceConfigEdit(config)
}

func ToJenkinsServiceConfigEdit(config *ServiceConfig) *JenkinsServiceConfigLoad {
	if config == nil {
		return &JenkinsServiceConfigLoad{ServiceConfigLoad: &ServiceConfigLoad{
			HookEvent: &HookEvent{},
		}}
	}

	jenkinsServiceConfig := &JenkinsServiceConfigLoad{}
	json.Unmarshal([]byte(config.ConfigContent), jenkinsServiceConfig)

	jenkinsServiceConfig.ReadEvent()

	jenkinsServiceConfig.ServiceConfigLoad.IsActive = config.IsActive

	return jenkinsServiceConfig
}
