package get

import (
	"fmt"
	g "github.com/agcom/bs/internal/cmd/global"
	"github.com/agcom/bs/internal/cmd/utils"
	"github.com/agcom/bs/internal/logs"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var cmdSlr = &cobra.Command{
	Use:     "seller",
	Short:   "read a seller",
	Aliases: []string{"slr"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logs.Fatalf("accepts 1 arg, received %d", len(args))
		}

		name := args[0]
		if strings.ContainsRune(name, filepath.Separator) {
			logs.Fatalf("invalid name; contains file path separator '%c'", filepath.Separator)
		}

		j, err := g.R.Slrs.Get(name)
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
