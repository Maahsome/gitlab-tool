package cmd

import (
	"github.com/spf13/cobra"
)

// deleteMrCmd represents the mr command
var deleteMrCmd = &cobra.Command{
	Use:   "mr",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// mrID, _ := cmd.Flags().GetInt("mr")

		// uri := "/projects/5784/releases/v0.0.7"

	},
}

func init() {
	deleteCmd.AddCommand(deleteMrCmd)
	deleteMrCmd.Flags().Int("mr", 0, "Specify the MR number")
}
