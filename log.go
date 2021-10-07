package main

import (
	"fmt"
	"go.uber.org/multierr"
	"log"
	"os"
)

var lgr = log.New(os.Stderr, "", 0)

func lvlf(lvl string, format string, v ...interface{}) {
	lgr.Printf("%s: %s\n", lvl, fmt.Sprintf(format, v...))
}

func lvl(lvl string, v interface{}) {
	lvlf(lvl, "%v", v)
}

func errorf(format string, v ...interface{}) {
	lvlf("Error", format, v...)
}

func errorr(v interface{}) {
	errorf("%v", v)
}

func warnf(format string, v ...interface{}) {
	lvlf("Warning", format, v...)
}

func warn(v interface{}) {
	warnf("%v", v)
}

func fatalfCode(code int, format string, v ...interface{}) {
	errorf(format, v...)
	os.Exit(code)
}

func fatalf(format string, v ...interface{}) {
	fatalfCode(1, format, v...)
}

func fatalCode(code int, v interface{}) {
	fatalfCode(code, "%v", v)
}

func fatal(v interface{}) {
	fatalf("%v", v)
}

func fatalErrs(errs []error) {
	for _, err := range errs {
		errorr(err)
	}
	os.Exit(1)
}

func fatalMultierr(err error) {
	errs := multierr.Errors(err)
	fatalErrs(errs)
}

func fatalErr(err error) {
	fatalMultierr(err)
}
