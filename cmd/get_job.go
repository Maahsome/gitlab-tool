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

type PipelineJob []struct {
	ID             int         `json:"id"`
	Status         string      `json:"status"`
	Stage          string      `json:"stage"`
	Name           string      `json:"name"`
	Ref            string      `json:"ref"`
	Tag            bool        `json:"tag"`
	Coverage       interface{} `json:"coverage"`
	AllowFailure   bool        `json:"allow_failure"`
	CreatedAt      time.Time   `json:"created_at"`
	StartedAt      time.Time   `json:"started_at"`
	FinishedAt     time.Time   `json:"finished_at"`
	Duration       float64     `json:"duration"`
	QueuedDuration float64     `json:"queued_duration"`
	User           struct {
		ID              int         `json:"id"`
		Name            string      `json:"name"`
		Username        string      `json:"username"`
		State           string      `json:"state"`
		AvatarURL       string      `json:"avatar_url"`
		WebURL          string      `json:"web_url"`
		CreatedAt       time.Time   `json:"created_at"`
		Bio             string      `json:"bio"`
		BioHTML         string      `json:"bio_html"`
		Location        string      `json:"location"`
		PublicEmail     string      `json:"public_email"`
		Skype           string      `json:"skype"`
		Linkedin        string      `json:"linkedin"`
		Twitter         string      `json:"twitter"`
		WebsiteURL      string      `json:"website_url"`
		Organization    string      `json:"organization"`
		JobTitle        string      `json:"job_title"`
		Pronouns        string      `json:"pronouns"`
		Bot             bool        `json:"bot"`
		WorkInformation interface{} `json:"work_information"`
		Followers       int         `json:"followers"`
		Following       int         `json:"following"`
	} `json:"user"`
	Commit struct {
		ID             string    `json:"id"`
		ShortID        string    `json:"short_id"`
		CreatedAt      time.Time `json:"created_at"`
		ParentIds      []string  `json:"parent_ids"`
		Title          string    `json:"title"`
		Message        string    `json:"message"`
		AuthorName     string    `json:"author_name"`
		AuthorEmail    string    `json:"author_email"`
		AuthoredDate   time.Time `json:"authored_date"`
		CommitterName  string    `json:"committer_name"`
		CommitterEmail string    `json:"committer_email"`
		CommittedDate  time.Time `json:"committed_date"`
		Trailers       struct {
		} `json:"trailers"`
		WebURL string `json:"web_url"`
	} `json:"commit"`
	Pipeline struct {
		ID        int       `json:"id"`
		ProjectID int       `json:"project_id"`
		Sha       string    `json:"sha"`
		Ref       string    `json:"ref"`
		Status    string    `json:"status"`
		Source    string    `json:"source"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		WebURL    string    `json:"web_url"`
	} `json:"pipeline"`
	WebURL    string `json:"web_url"`
	Artifacts []struct {
		FileType   string      `json:"file_type"`
		Size       int         `json:"size"`
		Filename   string      `json:"filename"`
		FileFormat interface{} `json:"file_format"`
	} `json:"artifacts"`
	Runner struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		IPAddress   string `json:"ip_address"`
		Active      bool   `json:"active"`
		IsShared    bool   `json:"is_shared"`
		RunnerType  string `json:"runner_type"`
		Name        string `json:"name"`
		Online      bool   `json:"online"`
		Status      string `json:"status"`
	} `json:"runner"`
	ArtifactsExpireAt interface{} `json:"artifacts_expire_at"`
	TagList           []string    `json:"tag_list"`
}

// jobCmd represents the jobs command
var jobCmd = &cobra.Command{
	Use:     "job",
	Aliases: []string{"jobs"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		plID, _ := cmd.Flags().GetInt("pipeline-id")
		getPipelineJobs(prID, plID)
	},
}

func getPipelineJobs(pr int, pl int) error {
	restClient := resty.New()

	uri := fmt.Sprintf("https://%s/api/v4/projects/%d/pipelines/%d/jobs", glHost, pr, pl)

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	// fmt.Println(string(resp.Body()[:]))
	var pj PipelineJob
	marshErr := json.Unmarshal(resp.Body(), &pj)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Pipeline", marshErr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "STATUS", "NAME", "TRACE"})
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

	for _, v := range pj {

		row := []string{
			fmt.Sprintf("%d", v.ID),
			v.Status,
			v.Name,
			fmt.Sprintf("<bash:gitlab-tool get trace -p %d -j %d>", v.Pipeline.ProjectID, v.ID),
		}
		table.Append(row)
	}
	table.Render()

	return nil
}

func init() {
	getCmd.AddCommand(jobCmd)

	jobCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	jobCmd.Flags().IntP("pipeline-id", "l", 0, "Specify the PipelineID")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
