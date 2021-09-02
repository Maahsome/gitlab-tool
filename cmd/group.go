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

type Group []struct {
	ID                             int         `json:"id"`
	WebURL                         string      `json:"web_url"`
	Name                           string      `json:"name"`
	Path                           string      `json:"path"`
	Description                    string      `json:"description"`
	Visibility                     string      `json:"visibility"`
	ShareWithGroupLock             bool        `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool        `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int         `json:"two_factor_grace_period"`
	ProjectCreationLevel           string      `json:"project_creation_level"`
	AutoDevopsEnabled              interface{} `json:"auto_devops_enabled"`
	SubgroupCreationLevel          string      `json:"subgroup_creation_level"`
	EmailsDisabled                 interface{} `json:"emails_disabled"`
	MentionsDisabled               interface{} `json:"mentions_disabled"`
	LfsEnabled                     bool        `json:"lfs_enabled"`
	DefaultBranchProtection        int         `json:"default_branch_protection"`
	AvatarURL                      interface{} `json:"avatar_url"`
	RequestAccessEnabled           bool        `json:"request_access_enabled"`
	FullName                       string      `json:"full_name"`
	FullPath                       string      `json:"full_path"`
	CreatedAt                      time.Time   `json:"created_at"`
	ParentID                       int         `json:"parent_id"`
	LdapCn                         interface{} `json:"ldap_cn"`
	LdapAccess                     interface{} `json:"ldap_access"`
}

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		glUser, _ := cmd.Flags().GetString("user")
		getGroup(glUser)
	},
}

func getGroup(user string) error {
	restClient := resty.New()

	uri := fmt.Sprintf("https://%s/api/v4/groups?per_page=100", glHost)

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	var gr Group
	marshErr := json.Unmarshal(resp.Body(), &gr)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Project", marshErr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "PATH", "PROJECTS"})
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

	for _, v := range gr {

		var cmdLine string
		if len(user) > 0 {
			cmdLine = fmt.Sprintf("<bash:gitlab-tool get project -g %d -u %s>", v.ID, user)
		} else {
			cmdLine = fmt.Sprintf("<bash:gitlab-tool get project -g %d>", v.ID)
		}
		row := []string{
			fmt.Sprintf("%d", v.ID),
			v.FullPath,
			cmdLine,
		}
		table.Append(row)
	}
	table.Render()

	// fmt.Println(string(resp.Body()[:]))

	return nil
}

func init() {
	getCmd.AddCommand(groupCmd)

	groupCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
