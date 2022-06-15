package cmd

import (
	"fmt"

	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configGroupCmd represents the group command
var configGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Configure a group to a local directory",
	Long: `EXAMPLE:

> gitlab-tool config group --gitlab-host git.alteryx.com --group futurama --directory ~/dev/futurama --env-var GITLAB_TOKEN_ALTERYX

EXAMPLE:

> gitlab-tool config group --gitlab-host gitlab.com --group alteryx_futurama --directory ~/src/alteryx_futurama --env-var GITLAB_TOKEN_PUBLIC
`,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabHost, _ := cmd.Flags().GetString("gitlab-host")
		topGroup, _ := cmd.Flags().GetString("group")
		directory, _ := cmd.Flags().GetString("directory")
		envvar, _ := cmd.Flags().GetString("env-var")

		if err := configureTool(gitlabHost, topGroup, directory, envvar); err != nil {
			logrus.WithError(err).Error("Failed to configure, please see previous output")
		}

	},
}

func configureTool(host string, group string, dir string, ev string) error {

	logrus.Info(fmt.Sprintf("gitlab-host: %s", host))
	logrus.Info(fmt.Sprintf("group: %s", group))
	logrus.Info(fmt.Sprintf("directory: %s", dir))
	logrus.Info(fmt.Sprintf("env-var: %s", ev))

	var configList objects.ConfigList
	err := viper.UnmarshalKey("configs", &configList)
	if err != nil {
		logrus.Fatal("Error unmarshalling...")
	}

	newConfig := objects.DirectoryConfig{
		Directory: dir,
		EnvVar:    ev,
		Group:     group,
		Host:      host,
	}

	configList = append(configList, newConfig)

	viper.Set("configs", configList)
	verr := viper.WriteConfig()
	if verr != nil {
		logrus.WithError(verr).Info("Failed to write config")
	} else {
		logrus.Info(fmt.Sprintf("Successfully saved gitlab-host (%s) and directory (%s) to config.yaml\n", host, dir))
	}
	return nil
}

func init() {
	configCmd.AddCommand(configGroupCmd)

	configGroupCmd.Flags().String("gitlab-host", "gitlab.com", "Specify the gitlab host")
	configGroupCmd.Flags().StringP("group", "g", "", "Specify the TOP level group name")
	configGroupCmd.Flags().StringP("directory", "d", "", "Specify the directory where the group sources will be housed eg. ~/src/<GROUP NAME>/")
	configGroupCmd.Flags().StringP("env-var", "e", "", "Environment VAR containing the access token to authenticate")

	configGroupCmd.MarkFlagRequired("group")
	configGroupCmd.MarkFlagRequired("directory")
	configGroupCmd.MarkFlagRequired("env-var")
}
