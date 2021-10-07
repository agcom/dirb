package main

import (
	"fmt"
	"go.uber.org/multierr"
	stdLog "log"
	"os"
)

var lgr = stdLog.New(os.Stderr, "", 0)

func logf(format string, v ...interface{}) {
	lgr.Printf("%s\n", fmt.Sprintf(format, v...))
}

func log(v interface{}) {
	lgr.Println(v)
}

func fatalfc(code int, format string, v ...interface{}) {
	logf(format, v...)
	os.Exit(code)
}

func fatalf(format string, v ...interface{}) {
	fatalfc(1, format, v...)
}

func fatalc(code int, v interface{}) {
	log(v)
	os.Exit(code)
}

func fatal(v interface{}) {
	fatalc(1, v)
}

func fatalErrs(errs []error) {
	for _, err := range errs {
		log(err)
	}
	os.Exit(1)
}

func fatalMultiErr(err error) {
	errs := multierr.Errors(err)
	fatalErrs(errs)
}

func fatalErr(err error) {
	fatalMultiErr(err)
}
