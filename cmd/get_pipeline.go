package cmd

import (
	"fmt"
	"strings"
	"time"

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Pipeline []struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Sha       string    `json:"sha"`
	Ref       string    `json:"ref"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	WebURL    string    `json:"web_url"`
}

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pipelines"},
	Short:   "Display a list of pipelines",
	Long: `EXAMPLE:
Get the pipelines for the current working directory project.  The output will
provide command lines that can be made to be clickable in iTerm2 in order to
list JOBS, and them follow on to list TRACES.

> gitlab-tool get pipelines

ID     	PROJECT ID	STATUS  	JOBS
1389624	6609      	success 	<bash:gitlab-tool get jobs -p 6609 -l 1389624>
1366864	6609      	success 	<bash:gitlab-tool get jobs -p 6609 -l 1366864>
1366850	6609      	success 	<bash:gitlab-tool get jobs -p 6609 -l 1366850>
1366833	6609      	failed  	<bash:gitlab-tool get jobs -p 6609 -l 1366833>
1366705	6609      	failed  	<bash:gitlab-tool get jobs -p 6609 -l 1366705>
1366182	6609      	success 	<bash:gitlab-tool get jobs -p 6609 -l 1366182>
...
...
...
`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		glUser, _ := cmd.Flags().GetString("user")

		if prID > 0 && cwdProjectID > 0 && prID != cwdProjectID {
			logrus.Warn(fmt.Sprintf("The projectID provided via --project-id (-p) doesn't match %d", cwdProjectID))
		}
		// Default to --project-id (-p) passed in
		if prID == 0 && cwdProjectID > 0 {
			prID = cwdProjectID
		}

		err := getPipeline(prID, glUser)
		if err != nil {
			logrus.WithError(err).Error("Bad, bad programmer, failed to fetch pipeline list")
		}
	},
}

func getPipeline(id int, user string) error {

	pipelines, err := gitClient.GetPipelines(id, user)
	if err != nil {
		// logrus.WithError(err).Error("Failed to fetch pipeline list")
		return err
	}

	fmt.Println(plDataToString(pipelines, fmt.Sprintf("%#v", pipelines)))

	return nil
}

func plDataToString(plData gl.Pipelines, raw string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return plData.ToJSON()
	case "gron":
		return plData.ToGRON()
	case "yaml":
		return plData.ToYAML()
	case "text", "table":
		return plData.ToTEXT(c.NoHeaders)
	default:
		return plData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	getCmd.AddCommand(pipelineCmd)

	pipelineCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	pipelineCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
}
