package cmd

import (
	"fmt"
	"strings"

	gl "github.com/maahsome/gitlab-go"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getVariablesCmd represents the mr command
var getVariablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "Get the CICD Variables for a ProjectID",
	Long: `EXAMPLE:
You are in a project directory AND have a configuration for the directory you are in, this will
return the MRs for your current project.

> gitlab-tool get variables

EXAMPLE:

You want to get an MR list for a specific project, using the ProjectID, from another gitlab project directory

> gitlab-tool get variables -p 6122
WARN[0001] The projectID provided via --project-id (-p) doesn't match 6123
IID	TITLE                                  	STATE 	AUTHOR      	CREATED            	DIFF
`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		showAll, _ := cmd.Flags().GetBool("all")
		includeProject, _ := cmd.Flags().GetBool("project-vars")

		if prID > 0 && cwdProjectID > 0 && prID != cwdProjectID {
			logrus.Warn(fmt.Sprintf("The projectID provided via --project-id (-p) doesn't match %d", cwdProjectID))
		}
		// Default to --project-id (-p) passed in
		if prID == 0 && cwdProjectID > 0 {
			prID = cwdProjectID
		}
		err := getVariables(prID, showAll, includeProject)
		if err != nil {
			logrus.WithError(err).Error("Bad, bad programmer")
		}
	},
}

func getVariables(id int, all bool, projects bool) error {

	if all {
		if projects {
			logrus.Warn("You have chosen to include PROJECT level variables, this will take awhile...")
		}
		// TODO: detect the TOP level group based on where you are at in the directory tree
		topGroupID, err := gitClient.GetGroupID(topGroupName)
		if err != nil {
			logrus.WithError(err).Fatal("Could not get the groupID for the top level group")
		}
		variables, err := gitClient.GetCicdVariablesFromGroup(topGroupID, projects)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}

		fmt.Println(variableDataToString(variables, fmt.Sprintf("%#v", variables)))
	} else {
		variables, err := gitClient.GetCicdVariables(id)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}

		fmt.Println(variableDataToString(variables, fmt.Sprintf("%#v", variables)))

	}
	return nil
}

func variableDataToString(varData gl.Variables, raw string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return varData.ToJSON()
	case "gron":
		return varData.ToGRON()
	case "yaml":
		return varData.ToYAML()
	case "text", "table":
		return varData.ToTEXT(c.NoHeaders)
	default:
		return varData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	getCmd.AddCommand(getVariablesCmd)

	getVariablesCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	getVariablesCmd.Flags().BoolP("all", "a", false, "Show ALL CICD Variables starting at the TOP level group")
	getVariablesCmd.Flags().Bool("project-vars", false, "Return PROJECT level CICD Variables as well (assumes -a) LONG RUN TIME")
}
