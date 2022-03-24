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
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		glUser, _ := cmd.Flags().GetString("user")
		glGroup, _ := cmd.Flags().GetInt("group")
		if glGroup > 0 {
			getGroup(glGroup, glUser)
		} else {
			getGroups(glUser)
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
}
