package jenkins_client

import "errors"

type TemplateHolder struct {
}

func NewWorkflowJobTemplate(jenkins *Jenkins) (*WorkflowJobItem, error) {

	plugins, err := jenkins.GetPlugins()

	if err != nil {
		return nil, errors.New("can not load plugins")
	}

	workflow_job, err := plugins.GetVersion("workflow-job")
	if err != nil {
		return nil, errors.New("can not can plugin version of  workflow-job")
	}

	workflow_cps, err := plugins.GetVersion("workflow-cps")

	if err != nil {
		return nil, errors.New("can not can plugin version of  workflow-cps")
	}
	jobItem := WorkflowJobItem{
		Plugin:           "workflow-job@" + workflow_job,
		KeepDependencies: "false",
		//Properties:       properties,
		Disabled: "false",
		Definition: CpsFlowDefinition{
			Class:   "org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition",
			Plugin:  "workflow-cps@" + workflow_cps,
			Sandbox: true,
		},
		plugins: plugins,
	}

	return &jobItem, nil
}
