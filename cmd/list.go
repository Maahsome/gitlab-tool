package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Produce a list of gitlab group/projects for the current directory",
	Long: `EXAMPLE:
> gitlab-tool list

ge-w- bender
gc--- farnsworth
ge+-- hermes
pe+-5 test-project

-     : group or project (g/p)
 -    : exists on disk (e exists, c created/group, c created during clone/project)
  -   : dirty (- clean, + dirty) (project only)
   -  : has worktrees (-/w) (project only)
    - : has # of MRs, + more than 9
`,
	Run: func(cmd *cobra.Command, args []string) {
		if inProject {
			fmt.Println("You are currently in a PROJECT directory, nothing to list")
		} else {
			gitlabListing()
		}
	},
}

func buildStatBlock(itemType string, path string, objID int) string {
	sb := "-----"

	// Set the Type of the entry
	sb = itemType + sb[1:]
	if exists(path) {
		sb = sb[:1] + "e" + sb[2:]
	} else {
		if itemType == "g" {
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				logrus.WithError(err).Error(fmt.Sprintf("Failed to create directory for group: %s", path))
			}
		}
		sb = sb[:1] + "c" + sb[2:]
	}

	if itemType == "p" {

		uri := fmt.Sprintf("/projects/%d/merge_requests?state=opened&per_page=50", objID)

		gitdata, err := gitClient.Get(uri)
		if err != nil {
			logrus.WithError(err).Error("Bad fetch from gitlab")
		}

		var mr objects.MergeRequestList
		marshErr := json.Unmarshal([]byte(gitdata), &mr)
		if marshErr != nil {
			logrus.Fatal("Cannot marshall Pipeline", marshErr)
		}

		if exists(fmt.Sprintf("%s/%s/.git", currentWorkDir, path)) {
			repo, rerr := git.PlainOpen(fmt.Sprintf("%s/%s", currentWorkDir, path))
			if rerr != nil {
				logrus.Fatal("Error retrieving git info")
			}
			worktree, wtErr := repo.Worktree()
			if wtErr != nil {
				logrus.WithError(err).Debug("Error creating worktree")
			}
			status, serr := worktree.Status()
			if serr != nil {
				logrus.WithError(serr).Debug("Failed to get a git status")
			}

			// Set the Dirty/Status
			if status.IsClean() {
				sb = sb[:2] + "-" + sb[3:]
			} else {
				sb = sb[:2] + "+" + sb[3:]
			}

			// Set the WorkTree Count
			// We will just do this manually, by examining .git/worktree/*
			if exists(fmt.Sprintf("%s/%s/.git/worktrees", currentWorkDir, path)) {
				worktreeDirs, err := os.ReadDir(fmt.Sprintf("%s/%s/.git/worktrees", currentWorkDir, path))
				if err != nil {
					logrus.WithError(err).Error("Failed to read the .git/worktree directory list")
					sb = sb[:3] + "-" + sb[4:]
				}
				if len(worktreeDirs) > 0 {
					sb = sb[:3] + "w" + sb[4:]
				}
			} else {
				sb = sb[:3] + "-" + sb[4:]
			}
		} else {
			// no .git directory, so these are empty
			sb = sb[:2] + "-" + sb[3:]
			sb = sb[:3] + "-" + sb[4:]
		}

		// Set the MR Count (field 5)
		mrCount := "-"
		if len(mr) > 9 {
			mrCount = "+"
		} else {
			mrCount = strconv.Itoa(len(mr))
		}
		sb = sb[:4] + mrCount + sb[5:]
	}

	return sb
}

func gitlabListing() {

	var list objects.GitListing

	// fmt.Printf("\tCWD: %s\n\tHost: %s\n\tGroupID: %d\n\tProjectID: %d\n", currentWorkDir, cwdGitlabHost, cwdGroupID, cwdProjectID)

	// use the current group, collect all the other GROUPS at this level
	subGroups, gerr := gitClient.GetSubGroups(cwdGroupID)
	if gerr != nil {
		logrus.WithError(gerr).Error("Failed to fetch sub-groups")
	}
	for _, v := range subGroups {
		list = append(list, objects.GitListItem{
			StatBlock: buildStatBlock("g", v.Path, v.ID),
			Path:      v.Path,
		})
		// fmt.Printf("%s -> %s\n", v.Path, v.FullPath)
	}

	// use the current group, collec t all the PROJECTS at this level
	projects, gerr := gitClient.GetGroupProjects(cwdGroupID)
	if gerr != nil {
		logrus.WithError(gerr).Error("Failed to fetch sub-groups")
	}
	for _, v := range projects {
		if v.Namespace.ID == cwdGroupID {
			list = append(list, objects.GitListItem{
				StatBlock: buildStatBlock("p", v.Path, v.ID),
				Path:      v.Path,
			})
		}
		// fmt.Printf("%s -> %s\n", v.Path, v.FullPath)
	}

	// output the data

	fmt.Println(glDataToString(list, ""))
}

func glDataToString(glData objects.GitListing, raw string) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return raw
	case "json":
		return glData.ToJSON()
	case "gron":
		return glData.ToGRON()
	case "yaml":
		return glData.ToYAML()
	case "text", "table":
		return glData.ToTEXT(c.NoHeaders)
	default:
		return glData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
