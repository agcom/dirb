package main

import "fmt"

func main() {
	err := parseArgs()
	if err != nil {
		fatalErr(fmt.Errorf("failed to parse the command line arguments; %w", err))
	}

	cmd()
}
