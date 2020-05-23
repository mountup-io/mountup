package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print mountup version",
	Long:  `print mountup version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mountup version 0.0.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
