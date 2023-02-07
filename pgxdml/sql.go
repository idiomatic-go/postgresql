package pgxdml

import "fmt"

const (
	TimestampFn = Function("now()")

	nextValFnFmt = "nextval('%s')"
	valueFmt     = "%v"
	stringFmt    = "'%v'"
	attrFmt      = "%v = %v"
)

func NextValFn(sequence string) Function {
	return Function(fmt.Sprintf(nextValFnFmt, sequence))
}

type Function string

var tokens = []string{"drop table", "delete from", "--", ";", "/*", "*/", "select * from"}

type Attr struct {
	Name string
	Val  any
}
