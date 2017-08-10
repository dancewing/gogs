package jenkins_client

import (
	"fmt"
	"testing"
	"time"
)

func NewJenkinsWithTestData() *Jenkins {
	var auth Auth
	auth = Auth{
		Username: "admin",
		ApiToken: "4f74a2b22446a37bb64b959edbb3397f",
	}
	return NewDebugJenkins(&auth, "http://localhost:9090")
}

func Test(t *testing.T) {
	jenkins := NewJenkinsWithTestData()
	jobs, err := jenkins.GetJobs()

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	if len(jobs) == 0 {
		t.Errorf("return no jobs\n")
	}

	fmt.Printf("%v", jobs)
}

func TestAddJobToView(t *testing.T) {
	jenkins := NewJenkinsWithTestData()

	scm := Scm{
		Class: "hudson.scm.SubversionSCM",
	}
	jobItem := MavenJobItem{
		Plugin:               "maven-plugin@2.7.1",
		Description:          "test description",
		Scm:                  scm,
		Triggers:             Triggers{},
		RunPostStepsIfResult: RunPostStepsIfResult{},
		Settings:             JobSettings{Class: "jenkins.mvn.DefaultSettingsProvider"},
		GlobalSettings:       JobSettings{Class: "jenkins.mvn.DefaultSettingsProvider"},
	}
	newJobName := fmt.Sprintf("test-with-view-%d", time.Now().UnixNano())
	newViewName := fmt.Sprintf("test-view-%d", time.Now().UnixNano())
	jenkins.CreateJob(jobItem, newJobName)
	jenkins.CreateView(NewListView(newViewName))

	job := Job{Name: newJobName}
	err := jenkins.AddJobToView(newViewName, job)

	if err != nil {
		t.Errorf("error %v\n", err)
	}
}

func TestCreateView(t *testing.T) {
	jenkins := NewJenkinsWithTestData()

	newViewName := fmt.Sprintf("test-view-%d", time.Now().UnixNano())
	err := jenkins.CreateView(NewListView(newViewName))

	if err != nil {
		t.Errorf("error %v\n", err)
	}
}

func TestCreateJobItem(t *testing.T) {

	jenkins := NewJenkinsWithTestData()

	//properties := []JobProperty{
	//	PipelineTriggersJobProperty{
	//
	//	}
	//}

	plugins, err := jenkins.GetPlugins()

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	workflow_job, err := plugins.GetVersion("workflow-job")
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	workflow_cps, err := plugins.GetVersion("workflow-cps")

	if err != nil {
		t.Errorf("error %v\n", err)
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
	}

	newJobName := fmt.Sprintf("test-%d", time.Now().UnixNano())
	err = jenkins.CreateJob(jobItem, newJobName)

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	jobs, _ := jenkins.GetJobs()
	foundNewJob := false
	for _, v := range jobs {
		if v.Name == newJobName {
			foundNewJob = true
		}
	}

	if !foundNewJob {
		t.Errorf("error %s not found\n", newJobName)
	}
}
