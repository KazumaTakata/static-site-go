package cmd

import (
	"github.com/KazumaTakata/static-site-go/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "server",
	Short: "run dev server",
	Long:  `run dev server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.RunDevServer()
	},
}
