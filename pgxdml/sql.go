package pgxdml

const (
	valueFmt  = "%v"
	stringFmt = "'%v'"
	attrFmt   = "%v = %v"
	Delete    = "DELETE"
	Update    = "UPDATE"
	Insert    = "INSERT"
	Select    = "SELECT"
)

type Function string

var tokens = []string{"drop table", "delete from", "--", ";", "/*", "*/", "select * from"}

type Attr struct {
	Name string
	Val  any
}
