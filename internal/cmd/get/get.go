package get

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "read",
	Short:   "read a book, seller, or buyer",
	Aliases: []string{"get"},
}

func init() {
	Cmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "pretty json output")
	Cmd.AddCommand(cmdBk, cmdSlr, cmdByr)
}
