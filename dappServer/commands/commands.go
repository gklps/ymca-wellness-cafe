package commands

import "github.com/spf13/cobra"

// Create the root command
var RootCmd = &cobra.Command{
	Use:   "rubix-interactive-cli",
	Short: "rubix-interactive-cli is a CLI application which will help you accustomed with various Rubix commands and APIs",
	Long:  `rubix-interactive-cli is a CLI application which will help you accustomed with various Rubix commands and APIs`,
}
