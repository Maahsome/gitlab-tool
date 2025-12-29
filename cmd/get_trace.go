package cmd

import (
	"fmt"

	"github.com/acarl005/stripansi"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// traceCmd represents the trace command
var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		prID, _ := cmd.Flags().GetInt("project-id")
		pjID, _ := cmd.Flags().GetInt("job-id")
		suppressTimestamp, _ := cmd.Flags().GetBool("suppress-timestamp")
		if c.StripTimestamp {
			// Simply do the opposite of the flag value
			// This way you can default to true, and use -s to switch to false
			suppressTimestamp = !suppressTimestamp
		}
		getJobTrace(prID, pjID, suppressTimestamp)
	},
}

func getJobTrace(pr int, pj int, suppressTimestamp bool) error {
	restClient := resty.New()

	uri := fmt.Sprintf("https://%s/api/v4/projects/%d/jobs/%d/trace", glHost, pr, pj)

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	plainText := stripansi.Strip(string(resp.Body()[:]))
	if suppressTimestamp {
		plainText = stripTimestamps(plainText)
	}
	fmt.Println(plainText)

	return nil
}

func stripTimestamps(input string) string {
	var output string
	lines := splitLines(input)
	// 2025-12-29T11:03:09.104840Z 00O Running with gitlab-runner 18.7.1 (cc7f9277)
	// Check the first line, where "Running with" starts, that is how many characters to drop for each line
	var dropChars int
	if len(lines) > 0 {
		for i, char := range lines[0] {
			if char == 'R' {
				dropChars = i
				break
			}
		}
	}

	for _, line := range lines {
		if len(line) > dropChars {
			output += line[dropChars:] + "\n"
		} else {
			output += line + "\n"
		}
	}
	return output
}

func splitLines(input string) []string {
	var lines []string
	currentLine := ""
	for _, char := range input {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func init() {
	getCmd.AddCommand(traceCmd)

	traceCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	traceCmd.Flags().IntP("job-id", "j", 0, "Specify the JobID")
	traceCmd.Flags().BoolP("suppress-timestamp", "s", false, "Suppress timestamps in the output")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
