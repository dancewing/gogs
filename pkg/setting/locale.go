package setting

import (
	"bytes"
	"strings"
)

const SPACE_INDENT = "\t"

type section struct {
	name    string
	locales []*locale

	buffer           bytes.Buffer
	spaceIndentCount int
}

type locale struct {
	key              string
	value            string
	buffer           bytes.Buffer
	spaceIndentCount int
	children         []*locale
}

func newSection(name string) *section {
	return &section{
		name:             name,
		buffer:           bytes.Buffer{},
		spaceIndentCount: 1,
	}
}

func (s *section) ToString() string {
	s.buffer.WriteString("{")
	s.buffer.WriteString("\n")

	s.buffer.WriteString(strings.Repeat(SPACE_INDENT, s.spaceIndentCount) + "\"" + s.name + "\": {")

	if len(s.locales) > 0 {

		s.buffer.WriteString("\n")

		len := len(s.locales)
		for i, l := range s.locales {
			if i < len-1 {
				s.buffer.WriteString(l.ToString() + ",")
				s.buffer.WriteString("\n")
			} else {
				s.buffer.WriteString(l.ToString())
			}
		}

	}

	s.buffer.WriteString("\n")
	s.buffer.WriteString(strings.Repeat(SPACE_INDENT, s.spaceIndentCount) + "}")

	s.buffer.WriteString("\n")
	s.buffer.WriteString("}")
	return s.buffer.String()
}

func (s *locale) ToString() string {
	s.buffer.WriteString(strings.Repeat(SPACE_INDENT, s.spaceIndentCount) + "\"" + s.key + "\": \"" + s.value + "\"")
	return s.buffer.String()
}

func (s *section) add(key string, value string) {
	//	keys := strings.Split(key, ".")

	if key == "" || value == "" {
		return
	}

	value = strings.Replace(value, "\"", "\\\"", 1000)

	value = strings.Replace(value, "\n", "", 1000)

	l := &locale{
		key:              key,
		value:            value,
		spaceIndentCount: s.spaceIndentCount + 1,
		buffer:           bytes.Buffer{},
	}
	s.locales = append(s.locales, l)
	//
	//if len(keys) == 1 {
	//	l := &locale{
	//		key:              keys[0],
	//		value:            value,
	//		spaceIndentCount: s.spaceIndentCount + 1,
	//		buffer:           bytes.Buffer{},
	//	}
	//	s.locales = append(s.locales, l)
	//	return
	//}
	//
	//length := len(keys)
	//for k := 0; k < length; k++ {
	//
	//}
}
