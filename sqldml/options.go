package sqldml

var test bool

func SetTestEnv(enabled bool) {
	test = enabled
}

func IsTestEnv() bool {
	return test
}
