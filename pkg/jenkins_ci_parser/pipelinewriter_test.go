package jenkins_ci_parser

import (
	"fmt"
	"testing"
)

func Test_Writer(t *testing.T) {

	writer := NewPipelineWriter()

	writer.NewTag("pipeline").
		NewLine("agent none").
		NewTag("stages").
		NewTagName("stage", "build").
		TagEnd().
		NewTagName("stage", "test").
		TagEnd().
		TagEnd().
		TagEnd()

	fmt.Print(writer.String())
}
