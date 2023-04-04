package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newResCmd)
}

var newResCmd = &cobra.Command{
	Use:   "[resource] new",
	Short: "Api resource",
	Long:  `aaa`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NEW")
	},
}
