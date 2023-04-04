package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Api resource",
	Long:  `aaa`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("IMPORT")
	},
}
