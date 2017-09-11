package main

import (
	"fmt"

	"github.com/restic/restic/internal/options"

	"github.com/spf13/cobra"
)

var optionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Print list of extended options",
	Long: `
The "options" command prints a list of extended options.
`,
	Hidden:            true,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("All Extended Options:\n")
		for _, opt := range options.List() {
			fmt.Printf("  %-15s   %s\n", opt.Namespace+"."+opt.Name, opt.Text)
		}
	},
}

func init() {
	cmdRoot.AddCommand(optionsCmd)
}
