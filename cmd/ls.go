package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "Api resource",
	Long:  `aaa`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LS")
	},
}
