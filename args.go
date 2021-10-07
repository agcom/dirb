package main

import (
	"fmt"
	"github.com/agcom/bs/sflag"
	"os"
)

var flagNoNxtArgVal = []string{"p", "pretty"}

var arg0 = os.Args[0]
var args = os.Args[1:]
var pArgs []string // Positional arguments
var flags []*sflag.Flag

func parseArgs() error {
	fs, pa, err := sflag.ParseNoNxtArgVal(args, flagNoNxtArgVal)
	flags = fs
	pArgs = pa

	if err != nil {
		return fmt.Errorf("failed to parse the command line arguments; %w", err)
	}

	return nil
}
