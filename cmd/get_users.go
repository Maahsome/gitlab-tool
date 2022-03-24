package cmd

import (
	"fmt"

	gl "github.com/maahsome/gitlab-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get a list of gitlab users",
	Long: `EXAMPLE:
> gitlab-tool get users

`,
	Run: func(cmd *cobra.Command, args []string) {
		userSearch, _ := cmd.Flags().GetString("search")
		getUsers(userSearch)
	},
}

func getUsers(search string) {

	gitClient := gl.New(glHost, "", glToken)

	gitdata, err := gitClient.GetUsers(search)
	if err != nil {
		logrus.WithError(err).Error("Bad fetch from gitlab")
	}

	fmt.Println(gitdata)
	// var mr objects.MergeRequest
	// marshErr := json.Unmarshal([]byte(gitdata), &mr)
	// if marshErr != nil {
	// 	logrus.Fatal("Cannot marshall Pipeline", marshErr)
	// }

	// fmt.Println(mrDataToString(mr, gitdata))

	// return nil
}

func init() {
	getCmd.AddCommand(usersCmd)

	usersCmd.Flags().StringP("search", "s", "", "Search for user")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
