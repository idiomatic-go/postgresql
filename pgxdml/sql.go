package pgxdml

import "fmt"

const (
	//NowFn       = "now()"
	TimestampFn = Function("now()")

	nextValFnFmt = "nextval('%s')"
	valueFmt     = "%v"
	stringFmt    = "'%v'"
	attrFmt      = "%v = %v"
	//Delete    = "DELETE"
	//Update    = "UPDATE"
	//Insert    = "INSERT"
	//Select    = "SELECT"
)

func NextValFn(sequence string) Function {
	return Function(fmt.Sprintf(nextValFnFmt, sequence))
}

//func TimestampFn() Function {
//	return Function(timestampFn)
//}

type Function string

var tokens = []string{"drop table", "delete from", "--", ";", "/*", "*/", "select * from"}

type Attr struct {
	Name string
	Val  any
}
