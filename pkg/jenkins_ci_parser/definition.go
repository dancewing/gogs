package jenkins_ci_parser

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

func (ci *JenkinsCI) DefaultEnvironment() Environment {
	return ci.Environments[0]
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

func (pipeline *Pipeline) Writer(writer *PipelineWriter) *PipelineWriter {
	var sub = writer.NewTag("pipeline")

	pipeline.Stages.Writer(sub)

	sub.TagEnd()

	return writer
}

func (pipeline *Pipeline) FilterStages(branch string, environment string) *Pipeline {
	result := &Pipeline{}

	result.Env = pipeline.Env
	result.Post = pipeline.Post
	result.Options = pipeline.Options

	result.Stages = make([]Stage, 0)

	for _, stage := range pipeline.Stages {
		if stage.Match(branch, environment) {
			result.Stages = append(result.Stages, stage)
		}
	}

	return result
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

func (stage Stage) Writer(writer *PipelineWriter) *PipelineWriter {

	sub := writer.NewTagName("stage", stage.Name)

	sub.NewTagName("gogsReportStatus", stage.Name).NewTag("steps").NewLines(stage.Steps).TagEnd().TagEnd()

	sub.TagEnd()

	return sub
}

func (stage Stage) Match(branch string, environment string) bool {

	if branch == "" && environment == "" {
		return true
	}

	if len(stage.When) > 0 {

		var (
			matchBranch bool
			foundBranch bool
			matchEnv    bool
			foundEnv    bool
		)

		for _, w := range stage.When {
			if w.Condition == "branch" {
				foundBranch = true
				if w.Value == branch {
					matchBranch = true
				}
			}
			if w.Condition == "environment" {
				foundEnv = true
				if w.Value == environment {
					matchEnv = true
				}
			}
		}

		if branch != "" && environment != "" {

			if foundBranch && foundEnv {
				return matchBranch && matchEnv
			} else if foundBranch {
				return matchBranch
			} else if foundEnv {
				return matchEnv
			}

		} else if branch != "" {
			if foundBranch {
				return matchBranch
			}
			return true
		} else if environment != "" {
			if foundEnv {
				return matchEnv
			}
			return true
		}

	} else {
		return true
	}
	return false
}

type Stages []Stage

func (stages Stages) Writer(writer *PipelineWriter) *PipelineWriter {

	sub := writer.NewTag("stages")

	for _, s := range stages {
		s.Writer(sub)
	}

	sub.TagEnd()
	return sub
}

type When struct {
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Whens []When
