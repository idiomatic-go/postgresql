package pgxsql

import "fmt"

type Vendor interface {
	Lookup(key string) string
}

type Lookup = func(key string) string

type Directory interface {
	Vendor
	Lookup
}

//type Lookup interface {
//	Directory
//	//Map
//}

func ExampleTypeSets() {
	//var d Directory
	var fn Lookup = Implementation

	fmt.Printf("FindValue : %v\n", Find("test", fn))

	//Output:
	// fail
}

func Find[T Lookup](key string, d T) string {
	if d == nil {
		return key
	}
	return d(key)
}

func Implementation(key string) string {
	return "implemented"
}
