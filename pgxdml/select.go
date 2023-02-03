package pgxdml

import (
	"errors"
	"strings"
)

const (
	firstParameterRef = "$1"
)

func ExpandSelect(template string, where []Attr) (string, error) {
	if template == "" {
		return template, errors.New("template is empty")
	}
	pos := strings.Index(template, firstParameterRef)
	if pos == -1 {
		return template, nil
	}
	if len(where) == 0 {
		return template, errors.New("where attributes are empty")
	}
	sb := strings.Builder{}
	sb.WriteString(template[:pos])
	WriteWhereAttributes(&sb, where)
	sb.WriteString(template[pos+len(firstParameterRef):])
	return sb.String(), nil
}
