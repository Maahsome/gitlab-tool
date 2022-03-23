package cmd

import (
	"fmt"

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getMembersCmd represents the members command
var getMembersCmd = &cobra.Command{
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

		getMembers(groupID, projectID)
	},
}

func getMembers(group, project int) {

	gitClient := gl.New(glHost, "", glToken)
	if group > 0 {
		groupdata, err := gitClient.GetGroupMembers(group)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}
		fmt.Println(groupdata)
	}
	// if project > 0 {
	// groupdata, err := gitClient.GetGroupMembers(group)
	// if err != nil {
	// 	logrus.WithError(err).Error("Bad fetch from gitlab")
	// }
	// }

	// var mr objects.MergeRequest
	// marshErr := json.Unmarshal([]byte(gitdata), &mr)
}

func init() {
	getCmd.AddCommand(getMembersCmd)

	getMembersCmd.Flags().IntP("group-id", "g", 0, "Specify the GroupID")
	getMembersCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
