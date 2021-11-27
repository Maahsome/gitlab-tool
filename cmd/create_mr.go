package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	gl "github.com/maahsome/gitlab-tool/cmd/gitlab"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	giturls "github.com/whilp/git-urls"
)

// MergeRequest represents a GitLab merge request.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/merge_requests.html
type MergeRequest struct {
	ID           int        `json:"id"`
	IID          int        `json:"iid"`
	TargetBranch string     `json:"target_branch"`
	SourceBranch string     `json:"source_branch"`
	ProjectID    int        `json:"project_id"`
	Title        string     `json:"title"`
	State        string     `json:"state"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Description  string     `json:"description"`
	MergeStatus  string     `json:"merge_status"`
	MergeError   string     `json:"merge_error"`
	MergedAt     *time.Time `json:"merged_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	HasConflicts bool       `json:"has_conflicts"`
	WebURL       string     `json:"web_url"`
}

var (
	gitlabClient *gl.Gitlab
)

// createMrCmd represents the mr command
var createMrCmd = &cobra.Command{
	Use:   "mr",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// gh pr create --title "RELEASE ${RELEASE_VERSION}" --body "Release ${RELEASE_VERSION}" --base main

		mrTitle, _ := cmd.Flags().GetString("title")
		// mrDescription, _ := cmd.Flags().GetString("description")
		mrSource, _ := cmd.Flags().GetString("source")
		mrTarget, _ := cmd.Flags().GetString("base")

		workDir, err := os.Getwd()
		if err != nil {
			logrus.Fatal("Failed to get the current working directory?  That is odd.")
		}

		workDir = "/Users/christopher.maahs/dev/falkor/tools/alterjira"
		repo, err := git.PlainOpen(workDir)
		if err != nil {
			logrus.Fatal("Error retrieving git info")
		}
		if len(mrSource) == 0 {
			head, err := repo.Head()
			if err != nil {
				logrus.Fatal("Error getting branch")
			}
			mrSource = strings.Replace(string(head.Name()), "refs/heads/", "", -1)
			if len(mrSource) == 0 {
				logrus.Fatal("Error grabbing local branch name")
			}
		}

		fmt.Printf("Source Branch: %s\n", mrSource)
		repoConfig, rcerr := repo.Config()
		if rcerr != nil {
			logrus.Fatal("Error getting Config")
		}
		pURLs, _ := giturls.Parse(repoConfig.Remotes["origin"].URLs[0])
		glSlug := strings.TrimPrefix(strings.TrimSuffix(pURLs.EscapedPath(), ".git"), "/")
		glSlug = url.PathEscape(glSlug)
		gitlabClient = gl.New(glHost, "", glToken)
		projectID, pierr := gitlabClient.GetProjectID(glSlug)
		if pierr != nil {
			logrus.Fatal("Could not get ProjectID from Slug", glSlug)
		}

		// mrURL, mrerr := CreateMR(projectID, mrTitle, mrDescription, mrSource, mrTarget)
		mrURL, mrerr := CreateMR(projectID, mrTitle, mrSource, mrTarget)
		if mrerr != nil {
			logrus.WithError(mrerr).Fatal("Failed to create MR")
		}
		fmt.Printf("Merge Request: %s\n", mrURL)
	},
}

// func CreateMR(projectID int, title string, body string, src string, dst string) (string, error) {
func CreateMR(projectID int, title string, src string, dst string) (string, error) {

	resp, rerr := gitlabClient.CreateMergeRequest(projectID, title, src, dst)
	if rerr != nil {
		logrus.Fatal("Failed to create the MR")
	}
	var mr MergeRequest
	marshErr := json.Unmarshal([]byte(resp), &mr)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall MR", marshErr)
	}

	return mr.WebURL, nil

}

func init() {
	createCmd.AddCommand(createMrCmd)

	createMrCmd.Flags().StringP("title", "t", "", "Specify the MR title")
	// createMrCmd.Flags().StringP("description", "d", "", "Specify the MR description")
	createMrCmd.Flags().StringP("source", "s", "", "Specify the MR source branch")
	createMrCmd.Flags().StringP("base", "b", "", "Specify the MR base/target branch")
	createMrCmd.MarkFlagRequired("title")
	createMrCmd.MarkFlagRequired("base")
}
