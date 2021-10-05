package logs

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func Lvlf(lvl string, format string, v ...interface{}) {
	log.Printf("%s: %s\n", lvl, fmt.Sprintf(format, v...))
}

func Lvl(lvl string, v interface{}) {
	log.Printf("%s: %v\n", lvl, v)
}

func Errorf(format string, v ...interface{}) {
	Lvlf("ERROR", format, v...)
}

func Error(v interface{}) {
	Errorf("%v", v)
}

func Warnf(format string, v ...interface{}) {
	Lvlf("WARN", format, v...)
}

func Warn(v interface{}) {
	Warnf("%v", v)
}

func Fatalf(format string, v ...interface{}) {
	Errorf(format, v...)
	os.Exit(1)
}

func Fatal(v interface{}) {
	Fatalf("%v", v)
}
