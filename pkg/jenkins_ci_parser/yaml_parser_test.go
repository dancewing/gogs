package jenkins_ci_parser

import (
	"encoding/json"
	"fmt"
	"strings"
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
          name: Test in only in master
          when :
            - branch : master
          steps:
            - sh 'mvn -B clean verify'
            - sh 'mvn -B clean verify'
        -
          name : deploy to production
          when :
             - environment : prod
          steps:
             - sh 'mvn compile'
`

func Test_Loading_File(t *testing.T) {

	content := file_content

	if ok := strings.Contains(content, "\t"); ok {
		content = strings.Replace(content, "\t", "", 1)
		fmt.Print("replace tab \n")
	}

	parser := NewGitLabCiYamlParser([]byte(content))

	ci, err := parser.ParseYaml()

	if err != nil {
		fmt.Printf("error %v", err)
		return
	}

	json, _ := json.MarshalIndent(ci, "", "  ")
	fmt.Printf("%s \n", json)

	writer := NewPipelineWriter(true)

	p := ci.Pipeline.FilterStages("branch1", "prod")

	p.Writer(writer)

	fmt.Printf("%s \n", writer.String())

}
