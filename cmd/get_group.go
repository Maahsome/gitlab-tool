package cmd

import (
	"fmt"
	"strings"

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Get a listing of gitlab group/groups",
	Long: `EXAMPLE:
List the group info of the group directory you are in, or the parent group of the
project that you are in.

> gitlab-tool get group

ID  	GROUP                                   	BASH
1887	falkor/gitlab-automation/child-pipelines	<bash:gitlab-tool get project -g 1887>

EXAMPLE:
List all of the groups in the gitlab instance

> gitlab-tool get group -a

ID  	GROUP                                   	BASH
620 	winston                                 	<bash:gitlab-tool get project -g 620>
145 	www                                     	<bash:gitlab-tool get project -g 145>
`,
	Run: func(cmd *cobra.Command, args []string) {
		glUser, _ := cmd.Flags().GetString("user")
		glGroup, _ := cmd.Flags().GetInt("group")
		showAll, _ := cmd.Flags().GetBool("all")

		if glGroup > 0 && cwdGroupID > 0 && glGroup != cwdGroupID {
			logrus.Warn(fmt.Sprintf("The groupID provided via --group (-g) doesn't match %d", cwdGroupID))
		}
		// Default to --project-id (-p) passed in
		if glGroup == 0 && cwdGroupID > 0 {
			glGroup = cwdGroupID
		}

		if showAll {
			getGroups(glUser)
		} else {
			getGroup(glGroup, glUser)
		}
	},
}

func getGroup(group int, user string) {

	gr, err := gitClient.GetGroup(group)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to fetch group list")
	}

	fmt.Println(gDataToString(gr, fmt.Sprintf("%#v", gr), user))

}

func getGroups(user string) {

	gr, err := gitClient.GetGroups("")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to fetch group list")
	}

	fmt.Println(grDataToString(gr, fmt.Sprintf("%#v", gr), user))
}

func gDataToString(gData gl.Group, raw string, user string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return gData.ToJSON()
	case "gron":
		return gData.ToGRON()
	case "yaml":
		return gData.ToYAML()
	case "text", "table":
		return gData.ToTEXT(c.NoHeaders, user)
	default:
		return gData.ToTEXT(c.NoHeaders, user)
	}
}

func grDataToString(grData gl.GroupList, raw string, user string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return grData.ToJSON()
	case "gron":
		return grData.ToGRON()
	case "yaml":
		return grData.ToYAML()
	case "text", "table":
		return grData.ToTEXT(c.NoHeaders, user)
	default:
		return grData.ToTEXT(c.NoHeaders, user)
	}
}

func init() {
	getCmd.AddCommand(groupCmd)

	groupCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
	groupCmd.Flags().IntP("group", "g", 0, "Specify a single group to fetch")
	groupCmd.Flags().BoolP("all", "a", false, "Show ALL Groups, normally only show parent group")
}
