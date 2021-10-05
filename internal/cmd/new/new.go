package new

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new book, seller, or buyer",
	Aliases: []string{"new", "add"},
}

func init() {
	Cmd.AddCommand(cmdBk, cmdSlr, cmdByr)
}
