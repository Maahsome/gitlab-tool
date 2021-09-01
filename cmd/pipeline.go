package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
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
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		glUser, _ := cmd.Flags().GetString("user")
		err := getPipeline(prID, glUser)
		if err != nil {
			logrus.WithError(err).Error("Bad, bad programmer")
		}
	},
}

func getPipeline(id int, user string) error {
	restClient := resty.New()

	var uri string
	if len(user) > 0 {
		uri = fmt.Sprintf("https://%s/api/v4/projects/%d/pipelines?username=%s", glHost, id, user)
	} else {
		uri = fmt.Sprintf("https://%s/api/v4/projects/%d/pipelines", glHost, id)
	}

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	var pl Pipeline
	marshErr := json.Unmarshal(resp.Body(), &pl)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Pipeline", marshErr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "PROJECT_ID", "STATUS", "JOBS"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, v := range pl {

		row := []string{
			fmt.Sprintf("%d", v.ID),
			fmt.Sprintf("%d", v.ProjectID),
			v.Status,
			fmt.Sprintf("bash:gitlab-tool get jobs -p %d -l %d", v.ProjectID, v.ID),
		}
		table.Append(row)
	}
	table.Render()

	// fmt.Println(string(resp.Body()[:]))

	return nil
}

func init() {
	getCmd.AddCommand(pipelineCmd)

	pipelineCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	pipelineCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
