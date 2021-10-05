package new

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "overwrite",
	Short:   "overwrite a book, seller, or buyer",
	Aliases: []string{"ow", "replace"},
}

func init() {
	Cmd.AddCommand(cmdBk, cmdSlr, cmdByr)
}
