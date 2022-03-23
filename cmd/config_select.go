package cmd

import (
	"fmt"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// selectConfigCmd represents the selectConfig command
var selectConfigCmd = &cobra.Command{
	Use:   "select",
	Short: "Select the active configuration",
	Long: `EXAMPLE:
`,
	Run: func(cmd *cobra.Command, args []string) {
		configName, _ := cmd.Flags().GetString("name")

		var hostList objects.HostList
		err := viper.UnmarshalKey("hosts", &hostList)
		if err != nil {
			logrus.Fatal("Error unmarshalling...")
		}
		found := false
		for _, v := range hostList {
			if strings.EqualFold(v.Name, configName) {
				found = true
				viper.Set("currenthost", v.Name)
				verr := viper.WriteConfig()
				if verr != nil {
					logrus.WithError(verr).Info("Failed to write config")
				} else {
					logrus.Info(fmt.Sprintf("Successfully saved gitlab-host (%s) to config.yaml\n", v.Name))
				}
			}
		}
		if !found {
			fmt.Println("No match was found for host, please use 'gitlab-tool config list' to get a list of hosts")
		}
	},
}

func init() {
	configCmd.AddCommand(selectConfigCmd)

	selectConfigCmd.Flags().StringP("name", "n", "", "Specify the config name")
	selectConfigCmd.MarkFlagRequired("name")

}
