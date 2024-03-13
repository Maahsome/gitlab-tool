package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	giturls "github.com/whilp/git-urls"
	"gopkg.in/yaml.v2"
)

// MergeRequest represents a GitLab merge request.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/merge_requests.html
type MergeRequest struct {
	ID               int        `json:"id"`
	IID              int        `json:"iid"`
	TargetBranch     string     `json:"target_branch"`
	SourceBranch     string     `json:"source_branch"`
	ProjectID        int        `json:"project_id"`
	Title            string     `json:"title"`
	State            string     `json:"state"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Description      string     `json:"description"`
	MergeStatus      string     `json:"merge_status"`
	MergeError       string     `json:"merge_error"`
	MergedAt         *time.Time `json:"merged_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	HasConflicts     bool       `json:"has_conflicts"`
	WebURL           string     `json:"web_url"`
	Error            string     `json:"error"`
	ErrorDescription string     `json:"error_description"`
	Scope            string     `json:"scope"`
}

type WorktreeRef struct {
	Gitdir string `yaml:"gitdir"`
}

var (
	gitlabClient gl.GitlabClient
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
		mrDescription, _ := cmd.Flags().GetString("description")
		mrSource, _ := cmd.Flags().GetString("source")
		mrTarget, _ := cmd.Flags().GetString("base")
		nosquash, _ := cmd.Flags().GetBool("no-sqash-on-merge")
		removeBranch, _ := cmd.Flags().GetBool("remove-source-branch")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// We want to squash by default, so to override would be adding
		// --no-squash-on-merge, which sets to true, so we want the opposite
		squash := !nosquash

		workDir, err := os.Getwd()
		if err != nil {
			logrus.Fatal("Failed to get the current working directory?  That is odd.")
		}

		if stat, err := os.Stat(fmt.Sprintf("%s/.git", workDir)); err == nil && !stat.IsDir() {
			if len(mrSource) == 0 {
				logrus.Fatal(fmt.Sprintf("The working directory is a worktree: %s\nPlease use --source to provide the source branch name\n", workDir))
			}
			mainDir, mderr := ReadWorkTreeRef(fmt.Sprintf("%s/.git", workDir))
			if mderr != nil {
				logrus.Fatal("Error getting main git directory")
			}
			workDir = mainDir
		}

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

		// fmt.Printf("Source Branch: %s\n", mrSource)
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
		if dryRun {
			fmt.Printf("Create an MR on:\n\tProjectID: %d\n\tTitle: %s\n\tSource Branch: %s\n\tTarget Branch: %s\n", projectID, mrTitle, mrSource, mrTarget)
		} else {
			mrURL, mrerr := CreateMR(projectID, mrTitle, mrSource, mrTarget, mrDescription, squash, removeBranch)
			if mrerr != nil {
				logrus.WithError(mrerr).Fatal("Failed to create MR")
			}
			fmt.Printf("%s", mrURL)
		}
	},
}

// func CreateMR(projectID int, title string, body string, src string, dst string) (string, error) {
func CreateMR(projectID int, title string, src string, dst string, description string, squash bool, removeSource bool) (string, error) {

	resp, rerr := gitlabClient.CreateMergeRequest(projectID, title, src, dst, description, squash, removeSource)
	if rerr != nil {
		logrus.Fatal("Failed to create the MR")
	}
	var mr MergeRequest
	marshErr := json.Unmarshal([]byte(resp), &mr)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall MR", marshErr)
	}

	if len(mr.Error) > 0 {
		logrus.Fatal(fmt.Sprintf("There was an error with the API call.\n\tError: %s\n\tDescription: %s\n\tScope: %s", mr.Error, mr.ErrorDescription, mr.Scope))
	}
	return mr.WebURL, nil

}

func init() {
	createCmd.AddCommand(createMrCmd)

	createMrCmd.Flags().StringP("title", "t", "", "Specify the MR title")
	createMrCmd.Flags().StringP("description", "d", "", "Specify the MR description")
	createMrCmd.Flags().StringP("source", "s", "", "Specify the MR source branch")
	createMrCmd.Flags().StringP("base", "b", "", "Specify the MR base/target branch")
	createMrCmd.Flags().BoolP("no-squash-on-merge", "q", false, "Set the MR to squash-on-merge")
	createMrCmd.Flags().BoolP("remove-source-branch", "r", true, "Set the MR to remove-source-branch")
	createMrCmd.Flags().Bool("dry-run", false, "No actions performed, just output what will happen")
	createMrCmd.MarkFlagRequired("title")
	createMrCmd.MarkFlagRequired("base")
}

func ReadWorkTreeRef(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	gwt := WorktreeRef{}
	if err := yaml.Unmarshal(bytes, &gwt); err != nil {
		return "", err
	}

	mainTree := strings.Split(gwt.Gitdir, ".git")
	return mainTree[0], nil
}
