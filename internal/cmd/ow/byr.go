package new

import (
	g "github.com/agcom/bs/internal/cmd/global"
	"github.com/agcom/bs/internal/cmd/utils"
	"github.com/agcom/bs/internal/logs"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var cmdByr = &cobra.Command{
	Use:     "buyer",
	Short:   "overwrite a buyer",
	Aliases: []string{"byr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logs.Fatalf("accepts 2 args, received %d", len(args))
		}

		name := args[0]
		if strings.ContainsRune(name, filepath.Separator) {
			logs.Fatalf("invalid name; contains file path separator '%c'", filepath.Separator)
		}

		s := args[1]
		j, err := utils.StrToJsn(s)
		if err != nil {
			logs.Fatalf("invalid json string; %v.", err)
		}

		err = g.R.Byrs.Ow(name, j)
		if err != nil {
			logs.Fatal(err)
		}
	},
}
