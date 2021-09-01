package cmd

import (
	"fmt"

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
		getJobTrace(prID, pjID)
	},
}

func getJobTrace(pr int, pj int) error {
	restClient := resty.New()

	uri := fmt.Sprintf("https://%s/api/v4/projects/%d/jobs/%d/trace", glHost, pr, pj)

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	fmt.Println(string(resp.Body()[:]))

	return nil
}

func init() {
	getCmd.AddCommand(traceCmd)

	traceCmd.Flags().IntP("project-id", "p", 0, "Specify the ProjectID")
	traceCmd.Flags().IntP("job-id", "j", 0, "Specify the JobID")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
