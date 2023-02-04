package pgxdml

import (
	"errors"
	"net/url"
	"strings"
)

func BuildWhere(url *url.URL) []Attr {
	if url == nil {
		return nil
	}
	values := url.Query()
	if len(values) == 0 {
		return nil
	}
	var where []Attr
	for k, v := range values {
		where = append(where, Attr{Name: k, Val: v[0]})
	}
	return where
}

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
