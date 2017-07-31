package jenkins

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"reflect"

	"strings"

	"github.com/Unknwon/com"
)

type JobDefinition struct {
	Host         string
	Docker       bool
	Environments []string
	Stages       []string

	Jobs     []*Job
	Pipeline *Pipeline
}

func (def *JobDefinition) String() string {
	return fmt.Sprintf("JobDefinition:\n Host: %s \n Docker : %t \n Environments : %s \n Stages: %s \n Jobs: \n %v \n Pipeline: \n %v",
		def.Host, def.Docker, strings.Join(def.Environments, ","), strings.Join(def.Stages, ","), def.Jobs, def.Pipeline)
}

type Pipeline struct {
	Jobs []*Job
}

type Job struct {
	JobName     string
	Stage       string
	Environment string
	JenkinsJob  string
	NextJobName string
	JenkinsHost string
}

func (job *Job) String() string {
	return fmt.Sprintf("\n---- Job --- \n JobName : %s \n Stage : %s \n Environment : %s \n JenkinsJob : %s \n NextJobName : %s \n JenkinsHost : %s \n", job.JobName, job.Stage, job.Environment, job.JenkinsJob, job.NextJobName, job.JenkinsHost)
}

func (def *JobDefinition) CheckPipeline() {
	if def.Jobs != nil && len(def.Jobs) > 0 {
		pJobs := make([]*Job, 0)
		def.Pipeline = &Pipeline{
			Jobs: pJobs,
		}
		def.Pipeline.Jobs = append(def.Pipeline.Jobs, def.Jobs[0])
		def.Pipeline.appendNextJob(def, def.Jobs[0].NextJobName)
	}
}

func (pipeline *Pipeline) appendNextJob(def *JobDefinition, next string) {

	if next == "" {
		return
	}
	nextJob := def.FindJobByName(next)
	if nextJob != nil {
		pipeline.Jobs = append(pipeline.Jobs, nextJob)
		pipeline.appendNextJob(def, nextJob.NextJobName)
	} else {
		return
	}
}

func (def *JobDefinition) FindJobByName(name string) *Job {
	if def.Jobs != nil && len(def.Jobs) > 0 {
		for _, j := range def.Jobs {
			if j.JobName == name {
				return j
			}
		}
	}
	return nil
}

func Parse(in []byte) (*JobDefinition, error) {

	def := &JobDefinition{}

	dataMap := yaml.MapSlice{}

	yaml.Unmarshal(in, &dataMap)

	for _, mi := range dataMap {
		key := strings.ToLower(com.ToStr(mi.Key))

		if key == "host" {
			def.Host = com.ToStr(mi.Value)
		} else if key == "docker" {
			def.Docker = mi.Value.(bool)
		} else if key == "environments" {
			def.Environments = convertStringSlice(mi.Value)
		} else if key == "stages" {
			def.Stages = convertStringSlice(mi.Value)
		} else {
			job := &Job{}
			job.JobName = key

			convertToJob(mi.Value, job)
			if def.Jobs == nil {
				def.Jobs = make([]*Job, 0)
			}
			def.Jobs = append(def.Jobs, job)
		}
	}

	def.CheckPipeline()

	return def, nil
}

func convertToJob(i interface{}, job *Job) {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Slice {
		result := make([]yaml.MapItem, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = val.Index(i).Interface().(yaml.MapItem)
		}
		for _, m := range result {
			key := strings.ToLower(com.ToStr(m.Key))

			if key == "stage" {
				job.Stage = com.ToStr(m.Value)
			} else if key == "environment" {
				job.Environment = com.ToStr(m.Value)
			} else if key == "job" {
				job.JenkinsJob = com.ToStr(m.Value)
			} else if key == "host" {
				job.JenkinsHost = com.ToStr(m.Value)
			} else if key == "next" {
				job.NextJobName = com.ToStr(m.Value)
			}
		}

	}
}

func convertStringSlice(i interface{}) []string {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Slice {
		result := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = com.ToStr(val.Index(i).Interface())
		}
		return result
	}
	return nil
}
