package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/objects"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mrCmd represents the mr command
var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "Get a list of MR for a ProjectID",
	Long: `EXAMPLE:
You are in a project directory AND have a configuration for the directory you are in, this will
return the MRs for your current project.

> gitlab-tool get mr

IID	TITLE                                  	STATE 	AUTHOR      	CREATED            	DIFF
2  	Configure WhiteSource for GitLab Server	opened	@whitesource	2022-03-03 06:37:31	<bash:gitlab-tool get diff -p 6122 -m 2>

EXAMPLE:

You want to get an MR list for a specific project, using the ProjectID, from another gitlab project directory

> gitlab-tool get mr -p 6122
WARN[0001] The projectID provided via --project-id (-p) doesn't match 6123
IID	TITLE                                  	STATE 	AUTHOR      	CREATED            	DIFF
2  	Configure WhiteSource for GitLab Server	opened	@whitesource	2022-03-03 06:37:31	<bash:gitlab-tool get diff -p 6122 -m 2>
`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		showAll, _ := cmd.Flags().GetBool("all")

		if prID > 0 && cwdProjectID > 0 && prID != cwdProjectID {
			logrus.Warn(fmt.Sprintf("The projectID provided via --project-id (-p) doesn't match %d", cwdProjectID))
		}
		// Default to --project-id (-p) passed in
		if prID == 0 && cwdProjectID > 0 {
			prID = cwdProjectID
		}
		err := getMergeRequest(prID, showAll)
		if err != nil {
			logrus.WithError(err).Error("Bad, bad programmer")
		}
	},
}

func getMergeRequest(id int, all bool) error {

	var uri string
	if all {
		uri = fmt.Sprintf("/projects/%d/merge_requests?state=all&per_page=50", id)
	} else {
		uri = fmt.Sprintf("/projects/%d/merge_requests?state=opened&per_page=50", id)
	}

	gitdata, err := gitClient.Get(uri)
	if err != nil {
		logrus.WithError(err).Error("Bad fetch from gitlab")
	}

	var mr objects.MergeRequest
	marshErr := json.Unmarshal([]byte(gitdata), &mr)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Pipeline", marshErr)
	}

	fmt.Println(mrDataToString(mr, gitdata))

	return nil
}

func mrDataToString(mrData objects.MergeRequest, raw string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return mrData.ToJSON()
	case "gron":
		return mrData.ToGRON()
	case "yaml":
		return mrData.ToYAML()
	case "text", "table":
		return mrData.ToTEXT(c.NoHeaders)
	default:
		return mrData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	getCmd.AddCommand(mrCmd)

	mrCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	mrCmd.Flags().BoolP("all", "a", false, "Show ALL MRs, normally only show 'opened'")
}
