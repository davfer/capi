package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CApi",
	Long:  `All software has versions. This is CApi's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NOT IMPLEMENTED (_embed)")
	},
}
