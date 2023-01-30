package sqldml

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"reflect"
	"strings"
	"time"
)

func TrimDoubleSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func FmtValue(v any) (string, error) {
	if template.IsNil(v) {
		return "NULL", nil
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		return "", errors.New(fmt.Sprintf("invalid argument : pointer types are not supported : %v", t.String()))
	}
	// Process time.Time first
	if t, ok := v.(time.Time); ok {
		return FmtTimestamp(t), nil
	}
	if t.Kind() != reflect.String {
		return fmt.Sprintf(valueFmt, v), nil
	}
	if _, function := v.(Function); function {
		return fmt.Sprintf(valueFmt, v), nil
	}
	err := SanitizeString(v.(string))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(stringFmt, v.(string)), nil
}

func FmtAttr(attr Attr) (string, error) {
	if attr.Name == "" {
		return "", errors.New("invalid attribute argument, attribute name is empty")
	}
	s, err := FmtValue(attr.Val)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(attrFmt, attr.Name, s), nil
}
