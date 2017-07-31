package jenkins_test

import (
	"testing"

	"fmt"

	"github.com/gogits/gogs/pkg/jenkins"
)

const file_content = `
#default jenkins host
host: http://localhost:9090/
#
docker: false
# ===================================================================
# Environments
# ===================================================================
environments:
  - sit
  - uat
  - pt
  - production
stages:
  - build
  - testing
  - deploy
######
# Task define below, first task will be run first
######
build:
  stage : build
  job: compile_job
  ### if you want to run next task
  next: testing
  ### if you wanto overwrite default jenkins host
  ###host:
testing:
  stage : testing
  job: testing
  next: deploy_sit
  ###


deploy_sit:
  stage: deploy
  environment: sit
  job : deploy_to_sit

`

func Test_Loading_File(t *testing.T) {

	def, _ := jenkins.Parse([]byte(file_content))

	fmt.Printf("%v \n", def)
}
