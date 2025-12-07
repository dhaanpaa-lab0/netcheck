package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies for netcheck",
	Long: `Install dependencies required for netcheck functionality.

This command helps set up required dependencies like Python for running
custom check scripts.`,
}

func init() {
	rootCmd.AddCommand(installCmd)
}
