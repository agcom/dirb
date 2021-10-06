package get

import (
	g "github.com/agcom/bs/internal/cmd/global"
	"github.com/agcom/bs/internal/logs"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var cmdByr = &cobra.Command{
	Use:     "buyer",
	Short:   "remove a buyer",
	Aliases: []string{"byr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logs.Fatalf("accepts 1 arg, received %d", len(args))
		}

		name := args[0]
		if strings.ContainsRune(name, filepath.Separator) {
			logs.Fatalf("invalid name; contains file path separator '%c'", filepath.Separator)
		}

		err := g.R.Byrs.Rm(name)
		if err != nil {
			logs.Fatal(err)
		}

	},
}
