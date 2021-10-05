package get

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "remove",
	Short:   "remove a book, seller, or buyer",
	Aliases: []string{"rm"},
}

func init() {
	Cmd.AddCommand(cmdBk, cmdSlr, cmdByr)
}
