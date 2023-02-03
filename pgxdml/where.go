package pgxdml

import (
	"errors"
	"strings"
)

func WriteWhere(sb *strings.Builder, terminate bool, attrs []Attr) error {
	max := len(attrs) - 1
	if max < 0 {
		return errors.New("invalid update where argument, attrs slice is empty")
	}
	sb.WriteString(Where)
	WriteWhereAttributes(sb, attrs)
	if terminate {
		sb.WriteString(";")
	}
	return nil
}

func WriteWhereAttributes(sb *strings.Builder, attrs []Attr) error {
	max := len(attrs) - 1
	if max < 0 {
		return errors.New("invalid update where argument, attrs slice is empty")
	}
	for i, attr := range attrs {
		s, err := FmtAttr(attr)
		if err != nil {
			return err
		}
		sb.WriteString(s)
		if i < max {
			sb.WriteString(And)
		}
	}
	//sb.WriteString(";")
	return nil
}
