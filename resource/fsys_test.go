package resource

import "fmt"

func Example_FS() {
	buf, err := ReadFile("fs/pgxsql/config_dev.txt")
	fmt.Printf("test: readFile() -> %v %v/n", err, string(buf))

	//Output:
	//fail
}
