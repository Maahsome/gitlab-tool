/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listConfigCmd represents the listConfig command
var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "List host configurations",
	Long: `EXAMPLE:
`,
	Run: func(cmd *cobra.Command, args []string) {
		var configList objects.ConfigList
		err := viper.UnmarshalKey("configs", &configList)
		if err != nil {
			logrus.Fatal("Error unmarshalling...")
		}

		fmt.Println(clDataToString(configList, fmt.Sprintf("%#v", configList)))

	},
}

func clDataToString(clData objects.ConfigList, raw string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return clData.ToJSON()
	case "gron":
		return clData.ToGRON()
	case "yaml":
		return clData.ToYAML()
	case "text", "table":
		return clData.ToTEXT(c.NoHeaders)
	default:
		return clData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	configCmd.AddCommand(listConfigCmd)
}
