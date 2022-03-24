package cmd

import (
	"fmt"

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AddMembersCmd represents the members command
var addMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	gitClient := gl.New(glHost, "", glToken)
	if group > 0 {
		memberdata, err := gitClient.AddGroupMember(group, user, access)
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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
