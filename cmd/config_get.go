package cmd

import (
	"fmt"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getConfigCmd represents the getConfig command
var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a host configuration",
	Long: `EXAMPLE:
`,
	Run: func(cmd *cobra.Command, args []string) {
		configName, _ := cmd.Flags().GetString("name")
		var hostList objects.HostList
		err := viper.UnmarshalKey("hosts", &hostList)
		if err != nil {
			logrus.Fatal("Error unmarshalling...")
		}
		for _, v := range hostList {
			if strings.EqualFold(v.Name, configName) {
				fmt.Printf("%s\t%s\t%s\n", v.Name, v.Host, v.EnvVar)
			}
		}

	},
}

func init() {
	configCmd.AddCommand(getConfigCmd)
	getConfigCmd.Flags().StringP("name", "n", "", "Specify the config name")
	getConfigCmd.MarkFlagRequired("name")
}
