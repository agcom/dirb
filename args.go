package main

import (
	"fmt"
	"github.com/agcom/dirb/sflag"
	"os"
)

var flagNoNxtArgVal = []string{"p", "pretty", "l", "left-operand-is-field-reference", "r", "right-operand-is-field-reference"}

var arg0 = os.Args[0]
var aArgs = os.Args[1:] // All arguments
var remArgs []string    // Remaining arguments
var flags []*sflag.Flag

func parseArgs() error {
	fs, pa, err := sflag.ParseNoNxtArgVal(aArgs, flagNoNxtArgVal)
	flags = fs
	remArgs = pa

	if err != nil {
		return fmt.Errorf("failed to parse the command line arguments; %w", err)
	}

	return nil
}
