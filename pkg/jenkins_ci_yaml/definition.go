package jenkins_ci_yaml

type Scripts []string

type PostCondition string

const (
	//always, changed, failure, success, unstable, and aborted
	PostConditonAlways   PostCondition = "always"
	PostConditonChanged  PostCondition = "changed"
	PostConditonFailure  PostCondition = "failure"
	PostConditonSuccess  PostCondition = "success"
	PostConditonUnstable PostCondition = "unstable"
	PostConditonAborted  PostCondition = "aborted"
)

type JenkinsCI struct {
	Environments []Environment `json:"environments"`
	Pipeline     *Pipeline     `json:"scripts"`
}

type Environment struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type Pipeline struct {
	Post    Post       `json:"post,omitempty"`
	Stages  Stages     `json:"stages"`
	Env     JenkinsEnv `json:"env,omitempty"`
	Options Options    `json:"options,omitempty"`
}

type JenkinsEnv struct {
	Scripts Scripts `json:"scripts"`
}

type Options struct {
	Scripts Scripts `json:"scripts,omitempty"`
}

type Post struct {
	Conditions []PostConditions `json:"conditions,omitempty"`
}

type PostConditions struct {
	Condition string  `json:"condition"`
	Scripts   Scripts `json:"scripts"`
}

type Stage struct {
	Name  string     `json:"name"`
	Steps Scripts    `json:"scripts"`
	When  Whens      `json:"whens,omitempty"`
	Env   JenkinsEnv `json:"env,omitempty"`
}

type Stages []Stage

type When struct {
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Whens []When
