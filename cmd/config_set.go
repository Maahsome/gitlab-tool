package cmd

import (
	"github.com/spf13/cobra"
)

// setConfigCmd represents the set command
var setConfigCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{""},
	Short:   "Set/Add a host configuration",
	Long: `EXAMPLE:
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	configCmd.AddCommand(setConfigCmd)
}
