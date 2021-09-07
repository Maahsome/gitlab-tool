package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	gl "github.com/maahsome/gitlab-tool/cmd/gitlab"
	"github.com/maahsome/gitlab-tool/cmd/objects"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mrCmd represents the mr command
var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "Get a list of MR for a ProjectID",
	Long: `EXAMPLE:
> gitlab-tool get mr -p 2342355`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		showAll, _ := cmd.Flags().GetBool("all")
		err := getMergeRequest(prID, showAll)
		if err != nil {
			logrus.WithError(err).Error("Bad, bad programmer")
		}
	},
}

func getMergeRequest(id int, all bool) error {
	gitClient := gl.New(glHost, "", glToken)

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
	// mrCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
