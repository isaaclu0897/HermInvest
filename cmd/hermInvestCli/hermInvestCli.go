package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// TODO: there is a lot of thing to do
// 1. Add testcases for 'add', 'delete', ...:
// * Test foreign key and date; complete test documentation.
// 2. Separate 'add', 'delete', and other functionalities into individual files; complete the missing parts.
// 3. Refactor 'version' functionality into a separate file.
// 4. Complete comprehensive documentation for CLI instructions.
// 5. Update the documentation.
// 6. Revisit CLI interface design for enhancements.

const (
	// hermInvestCli version
	version = "v0.6.4"
)

// version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hermInvestCli is %s\n", version)
	},
}

// stock
var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Stock management",
	Long:  `Manage the stock inventory via HermInvestCli.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

// root
var rootCmd = &cobra.Command{
	Use:  "hermInvestCli",
	Long: "Operate on the stock inventory for detailed management.",
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	rootCmd.Execute()
}
