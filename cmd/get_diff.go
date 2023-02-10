package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/maahsome/gitlab-tool/cmd/objects"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Retrieve a diff for an MR",
	Long: `EXAMPLE:
Get a diff for the current directory project

> gitlab-tool get diff -m 1

* displays a diff format *

EXAMPLE:
Get a diff specifying the project ID

> gitlab-tool get diff -p 28247395 -m 291

* displays diff format *

EXAMPLE:
Format the diff inside of a YAML block intended to submit to a Jira comment

> gitlabtool get diff -p 1234 -m 128 --j

* display the diff inside a YAML block
`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		mrID, _ := cmd.Flags().GetInt("mr-id")
		jira, _ := cmd.Flags().GetBool("jira")

		if prID > 0 && cwdProjectID > 0 && prID != cwdProjectID {
			logrus.Warn(fmt.Sprintf("The projectID provided via --project-id (-p) doesn't match %d", cwdProjectID))
		}
		// Default to --project-id (-p) passed in
		if prID == 0 && cwdProjectID > 0 {
			prID = cwdProjectID
		}

		getMergeRequestDiff(prID, mrID, jira)
	},
}

func getMergeRequestDiff(id int, iid int, jira bool) error {

	// var uri string
	uri := fmt.Sprintf("/projects/%d/merge_requests/%d/changes?access_raw_diffs=true", id, iid)

	gitdata, err := gitClient.Get(uri)
	if err != nil {
		logrus.WithError(err).Error("Bad fetch from gitlab")
	}

	var mrd objects.MergeRequestDiff
	marshErr := json.Unmarshal([]byte(gitdata), &mrd)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Pipeline", marshErr)
	}

	output := ""
	indentSpaces := ""
	if jira {
		output = `- type: paragraph
  data: |
    git diff
- type: codeBlock
  data: |
`
		indentSpaces = "    "
	}

	for _, c := range mrd[0].Changes {
		// diff --git \(.old_path) \(.new_path)\n--- \(.old_path)\n+++ \(.new_path)\n\(.diff)"' | delta
		output += fmt.Sprintf("%sdiff --git %s %s\n", indentSpaces, c.OldPath, c.NewPath)
		output += fmt.Sprintf("%s--- %s\n", indentSpaces, c.OldPath)
		output += fmt.Sprintf("%s+++ %s\n", indentSpaces, c.NewPath)
		lines := strings.Split(c.Diff, "\n")
		for _, l := range lines {
			output += fmt.Sprintf("%s%s\n", indentSpaces, l)
		}
	}

	if jira {
		fmt.Println(output)
	} else {
		// Could read $PAGER rather than hardcoding the path.
		customPager := os.Getenv("GITLAB_TOOL_PAGER")
		if len(customPager) == 0 {
			customPager = os.Getenv("PAGER")
			if len(customPager) == 0 {
				customPager = "more"
			}
		}
		cmd := exec.Command(customPager)
		// cmd := exec.Command("/usr/local/bin/delta")

		// Feed it with the string you want to display.
		cmd.Stdin = strings.NewReader(output)

		// This is crucial - otherwise it will write to a null device.
		cmd.Stdout = os.Stdout

		// Fork off a process and wait for it to terminate.
		pageerr := cmd.Run()
		if pageerr != nil {
			logrus.WithError(pageerr).Error("Error calling to pager")
		}
	}

	return nil
}

func init() {
	getCmd.AddCommand(diffCmd)

	diffCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	diffCmd.Flags().IntP("mr-id", "m", 0, "Specify the Merge Request IID")
	diffCmd.Flags().BoolP("jira", "j", false, "Output the diff in a JIRA comment update YAML block")
}
