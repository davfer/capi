package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newColCmd)
}

var newColCmd = &cobra.Command{
	Use:   "new",
	Short: "Api resource",
	Long:  `aaa`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NEW")
	},
}
