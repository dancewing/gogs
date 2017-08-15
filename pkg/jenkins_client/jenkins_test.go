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


	jobItem, err := NewWorkflowJobTemplate(jenkins)


	jobItem.Definition.Script = "testsssss"

	gogsProperty := &GogsProjectProperty{
		GogsSecret: "fs",
	}

	jobItem.AddProperty(gogsProperty, "gogs-webhook")


	newJobName := fmt.Sprintf("test-%d", time.Now().UnixNano())


	err = jenkins.CreateJob(jobItem, newJobName)

	fmt.Printf("jobName : %s \n", newJobName)

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	//jobs, _ := jenkins.GetJobs()

	job, err := jenkins.GetJob(newJobName)

	if err != nil {
		t.Errorf("error %s not found\n", newJobName)
	}

	fmt.Printf("%v", job)

	jobItem.Definition.Script = "update"

	err = jenkins.UpdateJob(jobItem, newJobName)

	if err != nil {
		t.Errorf("error %s not found\n", newJobName)
	}
	//foundNewJob := false
	//for _, v := range jobs {
	//	if v.Name == newJobName {
	//		foundNewJob = true
	//	}
	//}
	//
	//if !foundNewJob {
	//	t.Errorf("error %s not found\n", newJobName)
	//}
}

func TestGetJobByName(t *testing.T) {

	jenkins := NewJenkinsWithTestData()

	newJobName := "kuwago_jhipster-showcase_master_testing"

	//jobs, _ := jenkins.GetJobs()

	job, err := jenkins.GetJob(newJobName)

	if err != nil {
		t.Errorf("error %s not found\n", newJobName)
	}

	fmt.Printf("%v", job)

	//foundNewJob := false
	//for _, v := range jobs {
	//	if v.Name == newJobName {
	//		foundNewJob = true
	//	}
	//}
	//
	//if !foundNewJob {
	//	t.Errorf("error %s not found\n", newJobName)
	//}
}
