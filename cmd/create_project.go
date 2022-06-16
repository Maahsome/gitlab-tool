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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createProjectCmd represents the project command
var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new project in the group associated with the current directory",
	Long: `EXAMPLE:
Let say you are in 'gitlab-automation' directory, which is associated with the
'gitlab-automation' group.  And you want to create a new project named
'promotion-job' in that group.

> gitlab-tool create project -n "promotion-job" -v "internal"

* project created *

EXAMPLE:
Perhaps you want to override the group in which you are creating the new project.

> gitlab-tool create project -n "promotion-job" -v "internal" -g 5309

* project created *

`,
	Run: func(cmd *cobra.Command, args []string) {
		glGroup, _ := cmd.Flags().GetInt("group")
		projName, _ := cmd.Flags().GetString("name")
		projVisibility, _ := cmd.Flags().GetString("visibility")

		if glGroup > 0 && cwdGroupID > 0 && glGroup != cwdGroupID {
			logrus.Warn(fmt.Sprintf("The groupID provided via --group (-g) doesn't match %d", cwdGroupID))
		}
		// Default to --project-id (-p) passed in
		if glGroup == 0 && cwdGroupID > 0 {
			glGroup = cwdGroupID
		}

		createProject(projName, projVisibility, glGroup)
	},
}

func createProject(path string, visibility string, group int) {

	newProject, err := gitClient.CreateProject(group, path, visibility)
	if err != nil {
		logrus.WithError(err).Error("Failed to create project")
	}
	fmt.Printf("New project %s (%d) has been created", newProject.Path, newProject.ID)
}

func init() {
	createCmd.AddCommand(createProjectCmd)

	createProjectCmd.Flags().IntP("group", "g", 0, "Specify the group to create the project in")
	createProjectCmd.Flags().StringP("name", "n", "", "Specify the name/path of the project")
	createProjectCmd.Flags().StringP("visibility", "v", "internal", "Specify the visibility of the project")
	createProjectCmd.MarkFlagRequired("name")
	createProjectCmd.MarkFlagRequired("visibility")
}
