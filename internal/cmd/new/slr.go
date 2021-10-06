package new

import (
	g "bs/internal/cmd/global"
	"bs/internal/cmd/utils"
	"bs/internal/logs"
	"bs/jsns"
	"fmt"
	"github.com/spf13/cobra"
)

var cmdSlr = &cobra.Command{
	Use:     "seller",
	Short:   "create a new seller",
	Aliases: []string{"slr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logs.Fatalf("accepts 1 arg, received %d", len(args))
		}

		s := args[0]
		j, err := utils.StrToJsn(s)
		if err != nil {
			logs.Fatalf("invalid json string; %v.", err)
		}

		name, err := jsns.NewJsnGenName(g.R.Slrs, j)
		fmt.Println(name)
	},
}
