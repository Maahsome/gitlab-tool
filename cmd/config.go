package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set gitlab-tool configuration",
	Long: `EXAMPLE:
> gitlab-tool config set --name alteryx-private --host git.alteryx.com --env-var GLA_TOKEN

EXAMPLE:
> gitlab-tool config get --name alteryx-private

EXAMPLE:
> gitlab-tool config list

EXAMPLE:
> gitlab-tool config select --name alteryx-private
`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
