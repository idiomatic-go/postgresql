package pgxdml

import "strings"

func WriteDelete(sql string, where []Attr) (string, error) {
	var sb strings.Builder

	sb.WriteString(sql)
	sb.WriteString("\n")
	err := WriteWhere(&sb, where)
	return sb.String(), err
}
