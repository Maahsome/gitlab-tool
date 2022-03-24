package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete gitlab resources",
	Long: `EXAMPLE:
> gitlab-tool delete mr --id 5
`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("delete called")
	// },
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
