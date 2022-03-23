/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteReleaseCmd represents the release command
var deleteReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Delete a release from the repository",
	Long: `EXAMPLE:
> git-tool delete release --project-id --release-tag v0.0.7
`,
	Run: func(cmd *cobra.Command, args []string) {

		prID, _ := cmd.Flags().GetInt("project-id")
		releaseTag, _ := cmd.Flags().GetString("release-tag")

		uri := fmt.Sprintf("/projects/%d/releases/%s", prID, releaseTag)

		gitlabClient = gl.New(glHost, "", glToken)

		releaseInfo, derr := gitlabClient.Delete(uri)
		if derr != nil {
			logrus.Fatal("Failed to delete the Release", releaseTag)
		}
		fmt.Println(releaseInfo)

	},
}

func init() {
	deleteCmd.AddCommand(deleteReleaseCmd)
	deleteReleaseCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	deleteReleaseCmd.Flags().StringP("release-tag", "t", "", "Specify the Release TAG to remove")

	deleteReleaseCmd.MarkFlagRequired("project-id")
}
