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
> gitlab-tool get mr -p 2342355`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		showAll, _ := cmd.Flags().GetBool("all")

		// workDir, werr := os.Getwd()
		// if werr != nil {
		// 	logrus.Fatal("Failed to get the current working directory?  That is odd.")
		// }

		// gitDir := fmt.Sprintf("%s/.git", workDir)
		// if stat, err := os.Stat(gitDir); err == nil && !stat.IsDir() {
		// 	realDir, rerr := os.ReadFile(gitDir)
		// 	if rerr != nil {
		// 		logrus.Fatal("Failed to read the worktree gitdir...")
		// 	}
		// 	workDir = strings.Split(strings.TrimSpace(strings.TrimPrefix(string(realDir[:]), "gitdir: ")), ".git")[0]
		// }

		// repo, rerr := git.PlainOpen(workDir)
		// if rerr != nil {
		// 	logrus.Fatal("Error retrieving git info")
		// }
		// repoConfig, rcerr := repo.Config()
		// if rcerr != nil {
		// 	logrus.Fatal("Error getting Config")
		// }
		// // fmt.Printf("%#v\n", repoConfig)
		// pURLs, _ := giturls.Parse(repoConfig.Remotes["origin"].URLs[0])
		// glSlug := strings.TrimPrefix(strings.TrimSuffix(pURLs.EscapedPath(), ".git"), "/")
		// glSlug = url.PathEscape(glSlug)
		// gitlabClient = gl.New(glHost, "", glToken)
		// projectID, pierr := gitlabClient.GetProjectID(glSlug)
		// if pierr != nil {
		// 	logrus.Fatal("Could not get ProjectID from Slug", glSlug)
		// }

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
	// gitClient := gl.New(glHost, "", glToken)

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
	// mrCmd.MarkFlagRequired("project-id")
}
