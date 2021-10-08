package main

import (
	"fmt"
	"go.uber.org/multierr"
	stdLog "log"
	"os"
)

var lgr = stdLog.New(os.Stderr, "", 0)

func lvlf(lvl string, format string, v ...interface{}) {
	lgr.Printf("%s: %s\n", lvl, fmt.Sprintf(format, v...))
}

func lvl(lvl string, v interface{}) {
	lgr.Println(lvl+":", v)
}

func errorf(format string, v ...interface{}) {
	lvlf("Error", format, v...)
}

func errorr(v interface{}) {
	lvl("Error", v)
}

func errors(errs []error) {
	for _, err := range errs {
		errorr(err)
	}
}

func multiErr(err error) {
	errs := multierr.Errors(err)
	errors(errs)
}

func fatalfc(code int, format string, v ...interface{}) {
	errorf(format, v...)
	os.Exit(code)
}

func fatalf(format string, v ...interface{}) {
	fatalfc(1, format, v...)
}

func fatalc(code int, v interface{}) {
	errorr(v)
	os.Exit(code)
}

func fatal(v interface{}) {
	fatalc(1, v)
}

func fatalErrs(errs []error) {
	errors(errs)
	os.Exit(1)
}

func fatalMultiErr(err error) {
	multiErr(err)
	os.Exit(1)
}
