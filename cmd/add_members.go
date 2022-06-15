package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AddMembersCmd represents the members command
var addMembersCmd = &cobra.Command{
	Use:     "members",
	Aliases: []string{"member"},
	Short:   "Add a member to a group or project",
	Long: `EXAMPLE:
Add a user to group falkor/gitlab-automation-sandbox with default access level of 40

> gitlab-tool add members -g 2134 -u 441

EXAMPLE:
Add a user to project falkor/gitlab-automation-sandbox/gitlab-jobs with an access level of 50

> gitlab-tool add members -p 6104 -u 441 -a 50

`,
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetInt("group-id")
		projectID, _ := cmd.Flags().GetInt("project-id")
		userID, _ := cmd.Flags().GetInt("user-id")
		accessLevel, _ := cmd.Flags().GetInt("access-level")

		if groupID == 0 && projectID == 0 {
			logrus.Fatal("Must specify at least one of --group-id or --project-id")
		}
		if userID == 0 {
			logrus.Fatal("Must specify the --user-id to add to the member list")
		}
		addMembers(groupID, projectID, userID, accessLevel)
	},
}

func addMembers(group, project, user, access int) {

	if group > 0 {
		memberdata, err := gitClient.AddGroupMember(group, user, access)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}
		fmt.Println(memberdata)
	}

	if project > 0 {
		memberdata, err := gitClient.AddProjectMember(project, user, access)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}
		fmt.Println(memberdata)
	}

}

func init() {
	addCmd.AddCommand(addMembersCmd)

	addMembersCmd.Flags().IntP("group-id", "g", 0, "Specify the GroupID")
	addMembersCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	addMembersCmd.Flags().IntP("user-id", "u", 0, "Specify the UserID")
	addMembersCmd.Flags().IntP("access-level", "a", 30, "Specify the Access Level")
}
