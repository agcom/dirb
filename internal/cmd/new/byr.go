package new

import (
	"fmt"
	g "github.com/agcom/bs/internal/cmd/global"
	"github.com/agcom/bs/internal/cmd/utils"
	"github.com/agcom/bs/internal/logs"
	"github.com/agcom/bs/jsns"
	"github.com/spf13/cobra"
)

var cmdByr = &cobra.Command{
	Use:     "buyer",
	Short:   "create a new buyer",
	Aliases: []string{"byr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logs.Fatalf("accepts 1 arg, received %d", len(args))
		}

		s := args[0]
		j, err := utils.StrToJsn(s)
		if err != nil {
			logs.Fatalf("invalid json string; %v.", err)
		}

		name, err := jsns.NewJsnGenName(g.R.Byrs, j)
		fmt.Println(name)
	},
}
