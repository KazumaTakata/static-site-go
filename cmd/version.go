package cmd

import (
	"fmt"

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

		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
