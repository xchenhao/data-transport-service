package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func Fatalln(v ...interface{}) {
	Println(v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{})  {
	Printf(format, v...)
	os.Exit(1)
}

func Println(v ...interface{}) {
	file, line := getFileANDLineNumber()
	content := fmt.Sprintln(v...)
	log.Printf("[%s:%d] %s", file, line, content)
}

func Printf(format string, v ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, line := getFileANDLineNumber()
	content := fmt.Sprintf(format, v...)
	log.Printf("[%s:%d] %s", file, line, content)
}

func getFileANDLineNumber(skip ...int) (string, int) {
	var _skip int = 3
	if len(skip) > 0 {
		_skip = skip[0]
	}
	_, file, line, _ := runtime.Caller(_skip)
	return file, line
}

func GetFileLineStr() string {
	f, l := getFileANDLineNumber(2)
	return fmt.Sprintf("%s:%d", f, l)
}
