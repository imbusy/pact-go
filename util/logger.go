package util

type Logger interface {
	Println(args ...interface{})
	Printf(format string, args ...interface{})
}
