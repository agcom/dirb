package get

import (
	g "bs/internal/cmd/global"
	"bs/internal/cmd/utils"
	"bs/internal/logs"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var cmdByr = &cobra.Command{
	Use:     "buyer",
	Short:   "read a buyer",
	Aliases: []string{"byr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logs.Fatalf("accepts 1 arg, received %d", len(args))
		}

		name := args[0]
		if strings.ContainsRune(name, filepath.Separator) {
			logs.Fatalf("invalid name; contains file path separator '%c'", filepath.Separator)
		}

		j, err := g.R.Byrs.Get(name)
		if err != nil {
			logs.Fatal(err)
		}

		s, err := utils.JsnToStr(j, pretty)
		if err != nil {
			logs.Fatal(err)
		}

		if s[len(s)-1] == '\n' {
			fmt.Print(s)
		} else {
			fmt.Println(s)
		}
	},
}
