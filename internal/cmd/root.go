package cmd

import (
	"bs/internal"
	"bs/internal/cmd/get"
	g "bs/internal/cmd/global"
	"bs/internal/cmd/new"
	ow "bs/internal/cmd/ow"
	rm "bs/internal/cmd/rm"
	"bs/internal/logs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bs",
	Short: "manage a simple book store's data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tr, err := internal.NewDir(g.Dir)
		if err != nil {
			logs.Fatal(err)
		} else {
			g.R = *tr
		}
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&g.Dir, "directory", "d", ".", "data directory")
	rootCmd.AddCommand(new.Cmd, get.Cmd, ow.Cmd, rm.Cmd)
}
