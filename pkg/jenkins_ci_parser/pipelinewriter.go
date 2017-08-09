package jenkins_ci_parser

import (
	"bytes"
	"strings"
)

const SPACE_INDENT = "    "

type PipelineWriter struct {
	buffer           bytes.Buffer
	parent           *PipelineWriter
	spaceIndentCount int
}

func (writer *PipelineWriter) String() string {
	return writer.buffer.String()
}

func (writer *PipelineWriter) WriteString(content string) {
	writer.buffer.WriteString(content)
}

func (writer *PipelineWriter) NewTag(tag string) *PipelineWriter {
	writer.WriteString(strings.Repeat(SPACE_INDENT, writer.spaceIndentCount) + tag + " {")
	writer.WriteString("\n")

	child := &PipelineWriter{
		buffer:           bytes.Buffer{},
		parent:           writer,
		spaceIndentCount: writer.spaceIndentCount + 1,
	}
	return child
}

func (writer *PipelineWriter) NewTagName(tag string, name string) *PipelineWriter {
	writer.WriteString(strings.Repeat(SPACE_INDENT, writer.spaceIndentCount) + tag + " ('" + name + "') {")
	writer.WriteString("\n")

	child := &PipelineWriter{
		buffer:           bytes.Buffer{},
		parent:           writer,
		spaceIndentCount: writer.spaceIndentCount + 1,
	}
	return child
}

func (writer *PipelineWriter) TagEnd() *PipelineWriter {

	if writer.parent != nil {
		writer.parent.WriteString(writer.String())
		writer.parent.WriteString(strings.Repeat(SPACE_INDENT, writer.parent.spaceIndentCount) + "}")
		writer.parent.WriteString("\n")
		return writer.parent
	}

	writer.WriteString(strings.Repeat(SPACE_INDENT, writer.spaceIndentCount) + "}")
	writer.WriteString("\n")
	return writer
}

func (writer *PipelineWriter) NewLine(content string) *PipelineWriter {
	writer.WriteString(strings.Repeat(SPACE_INDENT, writer.spaceIndentCount) + content + "")
	writer.WriteString("\n")
	return writer
}

func (writer *PipelineWriter) NewLines(lines []string) *PipelineWriter {

	for _, c := range lines {
		writer.WriteString(strings.Repeat(SPACE_INDENT, writer.spaceIndentCount) + c + "")
		writer.WriteString("\n")
	}

	return writer
}

func NewPipelineWriter() *PipelineWriter {
	return &PipelineWriter{
		buffer: bytes.Buffer{},
	}
}
