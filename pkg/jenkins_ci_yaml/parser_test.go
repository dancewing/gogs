package jenkins_ci_yaml

import (
	"encoding/json"
	"fmt"
	"testing"
)

const file_content = `
environments:
   - testing : http://localhost:8080
   - sit :  http://localhost:8090
pipeline :
    stages :
        -
          name : Example Build 111
          steps:
            - sh 'mvn -B clean verify'
            - sh 'mvn -B clean verify'
        -
          name: Example Build 22
          when :
             - branch : master
             - environment : production
          steps:
            - sh 'mvn -B clean verify'
            - sh 'mvn -B clean verify'

`

func Test_Loading_File(t *testing.T) {

	parser := NewGitLabCiYamlParser([]byte(file_content))

	ci, err := parser.ParseYaml()

	if err != nil {
		fmt.Printf("error %v", err)
		return
	}

	json, _ := json.MarshalIndent(ci, "", "  ")
	fmt.Printf("%s \n", json)
}
