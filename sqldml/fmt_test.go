package sqldml

import (
	"fmt"
)

/*
func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}


*/
func ExampleFmtValues() {
	var ptr *int

	v, err := FmtValue(nil)
	fmt.Printf("Value  [nil]  : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	v, err = FmtValue(ptr)
	fmt.Printf("Value  [ptr]  : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	var n = 123
	v, err = FmtValue(&n)
	fmt.Printf("Value  [ptr]  : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	v, err = FmtValue(true)
	fmt.Printf("Value  [true] : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	v, err = FmtValue(1001)
	fmt.Printf("Value  [1001] : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	v, err = FmtValue("")
	fmt.Printf("Value  [\"\"]   : %v\n", NilEmpty(v))
	fmt.Printf("Error         : %v\n", err)

	//t := time.Now()
	//v, err = FmtValue(t)
	//fmt.Printf("Value  [time.Now] : %v\n", NilEmpty(v[:19]))
	//fmt.Printf("Error             : %v\n", err)

	v, err = FmtValue("test string")
	fmt.Printf("Value  [test string]  : %v\n", NilEmpty(v))
	fmt.Printf("Error                 : %v\n", err)

	v, err = FmtValue(Function("now()"))
	fmt.Printf("Value  [now()]  : %v\n", NilEmpty(v))
	fmt.Printf("Error           : %v\n", err)

	v, err = FmtValue("drop table")
	fmt.Printf("Value  [drop table]  : %v\n", NilEmpty(v))
	fmt.Printf("Error                : %v\n", err)

	//Output:
	//Value  [nil]  : NULL
	//Error         : <nil>
	//Value  [ptr]  : NULL
	//Error         : <nil>
	//Value  [ptr]  : <nil>
	//Error         : invalid argument : pointer types are not supported : *int
	//Value  [true] : true
	//Error         : <nil>
	//Value  [1001] : 1001
	//Error         : <nil>
	//Value  [""]   : ''
	//Error         : <nil>
	//Value  [test string]  : 'test string'
	//Error                 : <nil>
	//Value  [now()]  : now()
	//Error           : <nil>
	//Value  [drop table]  : <nil>
	//Error                : SQL injection embedded in string [drop table] : drop table
}

func ExampleFmtAttr() {

	s, err := FmtAttr(Attr{})
	fmt.Printf("Name  [\"\"]  : %v\n", NilEmpty(s))
	fmt.Printf("Error       : %v\n", err)

	s, err = FmtAttr(Attr{Name: "attr_name_1"})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Name: "attr_name_2", Val: 1234})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Name: "attr_name_3", Val: false})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	//s, err = FmtAttr(util.Attr{Name: "attr_name_4", Val: time.Now()})
	//fmt.Println("default format:", time.Now())
	//fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	//fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Name: "attr_name_5", Val: "value string"})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Name: "attr_name_6", Val: Function("now()")})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	//Output:
	//Name  [""]  : <nil>
	//Error       : invalid attribute argument, attribute name is empty
	//Name  [attr_name]  : attr_name_1 = NULL
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_2 = 1234
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_3 = false
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_5 = 'value string'
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_6 = now()
	//Error              : <nil>

}
