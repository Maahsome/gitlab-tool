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

type Project []struct {
	ID                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         time.Time     `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	Topics            []interface{} `json:"topics"`
	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string        `json:"http_url_to_repo"`
	WebURL            string        `json:"web_url"`
	ReadmeURL         string        `json:"readme_url"`
	AvatarURL         interface{}   `json:"avatar_url"`
	ForksCount        int           `json:"forks_count"`
	StarCount         int           `json:"star_count"`
	LastActivityAt    time.Time     `json:"last_activity_at"`
	Namespace         struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Path      string      `json:"path"`
		Kind      string      `json:"kind"`
		FullPath  string      `json:"full_path"`
		ParentID  int         `json:"parent_id"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"namespace"`
	ContainerRegistryImagePrefix string `json:"container_registry_image_prefix"`
	Links                        struct {
		Self          string `json:"self"`
		Issues        string `json:"issues"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Labels        string `json:"labels"`
		Events        string `json:"events"`
		Members       string `json:"members"`
	} `json:"_links"`
	PackagesEnabled                bool   `json:"packages_enabled"`
	EmptyRepo                      bool   `json:"empty_repo"`
	Archived                       bool   `json:"archived"`
	Visibility                     string `json:"visibility"`
	ResolveOutdatedDiffDiscussions bool   `json:"resolve_outdated_diff_discussions"`
	ContainerExpirationPolicy      struct {
		Cadence       string      `json:"cadence"`
		Enabled       bool        `json:"enabled"`
		KeepN         int         `json:"keep_n"`
		OlderThan     string      `json:"older_than"`
		NameRegex     string      `json:"name_regex"`
		NameRegexKeep interface{} `json:"name_regex_keep"`
		NextRunAt     time.Time   `json:"next_run_at"`
	} `json:"container_expiration_policy"`
	IssuesEnabled                             bool          `json:"issues_enabled"`
	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled"`
	WikiEnabled                               bool          `json:"wiki_enabled"`
	JobsEnabled                               bool          `json:"jobs_enabled"`
	SnippetsEnabled                           bool          `json:"snippets_enabled"`
	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled"`
	ServiceDeskEnabled                        bool          `json:"service_desk_enabled"`
	ServiceDeskAddress                        string        `json:"service_desk_address"`
	CanCreateMergeRequestIn                   bool          `json:"can_create_merge_request_in"`
	IssuesAccessLevel                         string        `json:"issues_access_level"`
	RepositoryAccessLevel                     string        `json:"repository_access_level"`
	MergeRequestsAccessLevel                  string        `json:"merge_requests_access_level"`
	ForkingAccessLevel                        string        `json:"forking_access_level"`
	WikiAccessLevel                           string        `json:"wiki_access_level"`
	BuildsAccessLevel                         string        `json:"builds_access_level"`
	SnippetsAccessLevel                       string        `json:"snippets_access_level"`
	PagesAccessLevel                          string        `json:"pages_access_level"`
	OperationsAccessLevel                     string        `json:"operations_access_level"`
	AnalyticsAccessLevel                      string        `json:"analytics_access_level"`
	ContainerRegistryAccessLevel              string        `json:"container_registry_access_level"`
	EmailsDisabled                            interface{}   `json:"emails_disabled"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	LfsEnabled                                bool          `json:"lfs_enabled"`
	CreatorID                                 int           `json:"creator_id"`
	ImportStatus                              string        `json:"import_status"`
	OpenIssuesCount                           int           `json:"open_issues_count"`
	CiDefaultGitDepth                         int           `json:"ci_default_git_depth"`
	CiForwardDeploymentEnabled                bool          `json:"ci_forward_deployment_enabled"`
	CiJobTokenScopeEnabled                    bool          `json:"ci_job_token_scope_enabled"`
	PublicJobs                                bool          `json:"public_jobs"`
	BuildTimeout                              int           `json:"build_timeout"`
	AutoCancelPendingPipelines                string        `json:"auto_cancel_pending_pipelines"`
	BuildCoverageRegex                        interface{}   `json:"build_coverage_regex"`
	CiConfigPath                              string        `json:"ci_config_path"`
	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	AllowMergeOnSkippedPipeline               interface{}   `json:"allow_merge_on_skipped_pipeline"`
	RestrictUserDefinedVariables              bool          `json:"restrict_user_defined_variables"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge"`
	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
	MergeMethod                               string        `json:"merge_method"`
	SquashOption                              string        `json:"squash_option"`
	SuggestionCommitMessage                   interface{}   `json:"suggestion_commit_message"`
	AutoDevopsEnabled                         bool          `json:"auto_devops_enabled"`
	AutoDevopsDeployStrategy                  string        `json:"auto_devops_deploy_strategy"`
	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues"`
	KeepLatestArtifact                        bool          `json:"keep_latest_artifact"`
	ExternalAuthorizationClassificationLabel  string        `json:"external_authorization_classification_label"`
	RequirementsEnabled                       bool          `json:"requirements_enabled"`
	SecurityAndComplianceEnabled              bool          `json:"security_and_compliance_enabled"`
	ComplianceFrameworks                      []interface{} `json:"compliance_frameworks"`
}

// type Project []struct {
// 	ID                int           `json:"id"`
// 	Description       string        `json:"description"`
// 	Name              string        `json:"name"`
// 	NameWithNamespace string        `json:"name_with_namespace"`
// 	Path              string        `json:"path"`
// 	PathWithNamespace string        `json:"path_with_namespace"`
// 	CreatedAt         time.Time     `json:"created_at"`
// 	DefaultBranch     string        `json:"default_branch"`
// 	TagList           []interface{} `json:"tag_list"`
// 	Topics            []interface{} `json:"topics"`
// 	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
// 	HTTPURLToRepo     string        `json:"http_url_to_repo"`
// 	WebURL            string        `json:"web_url"`
// 	ReadmeURL         interface{}   `json:"readme_url"`
// 	AvatarURL         interface{}   `json:"avatar_url"`
// 	ForksCount        int           `json:"forks_count"`
// 	StarCount         int           `json:"star_count"`
// 	LastActivityAt    time.Time     `json:"last_activity_at"`
// 	Namespace         struct {
// 		ID        int         `json:"id"`
// 		Name      string      `json:"name"`
// 		Path      string      `json:"path"`
// 		Kind      string      `json:"kind"`
// 		FullPath  string      `json:"full_path"`
// 		ParentID  int         `json:"parent_id"`
// 		AvatarURL interface{} `json:"avatar_url"`
// 		WebURL    string      `json:"web_url"`
// 	} `json:"namespace"`
// 	ContainerRegistryImagePrefix string `json:"container_registry_image_prefix"`
// 	Links                        struct {
// 		Self          string `json:"self"`
// 		Issues        string `json:"issues"`
// 		MergeRequests string `json:"merge_requests"`
// 		RepoBranches  string `json:"repo_branches"`
// 		Labels        string `json:"labels"`
// 		Events        string `json:"events"`
// 		Members       string `json:"members"`
// 	} `json:"_links"`
// 	PackagesEnabled                bool   `json:"packages_enabled"`
// 	EmptyRepo                      bool   `json:"empty_repo"`
// 	Archived                       bool   `json:"archived"`
// 	Visibility                     string `json:"visibility"`
// 	ResolveOutdatedDiffDiscussions bool   `json:"resolve_outdated_diff_discussions"`
// 	ContainerExpirationPolicy      struct {
// 		Cadence       string      `json:"cadence"`
// 		Enabled       bool        `json:"enabled"`
// 		KeepN         int         `json:"keep_n"`
// 		OlderThan     string      `json:"older_than"`
// 		NameRegex     string      `json:"name_regex"`
// 		NameRegexKeep interface{} `json:"name_regex_keep"`
// 		NextRunAt     time.Time   `json:"next_run_at"`
// 	} `json:"container_expiration_policy"`
// 	IssuesEnabled                             bool          `json:"issues_enabled"`
// 	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled"`
// 	WikiEnabled                               bool          `json:"wiki_enabled"`
// 	JobsEnabled                               bool          `json:"jobs_enabled"`
// 	SnippetsEnabled                           bool          `json:"snippets_enabled"`
// 	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled"`
// 	ServiceDeskEnabled                        bool          `json:"service_desk_enabled"`
// 	ServiceDeskAddress                        string        `json:"service_desk_address"`
// 	CanCreateMergeRequestIn                   bool          `json:"can_create_merge_request_in"`
// 	IssuesAccessLevel                         string        `json:"issues_access_level"`
// 	RepositoryAccessLevel                     string        `json:"repository_access_level"`
// 	MergeRequestsAccessLevel                  string        `json:"merge_requests_access_level"`
// 	ForkingAccessLevel                        string        `json:"forking_access_level"`
// 	WikiAccessLevel                           string        `json:"wiki_access_level"`
// 	BuildsAccessLevel                         string        `json:"builds_access_level"`
// 	SnippetsAccessLevel                       string        `json:"snippets_access_level"`
// 	PagesAccessLevel                          string        `json:"pages_access_level"`
// 	OperationsAccessLevel                     string        `json:"operations_access_level"`
// 	AnalyticsAccessLevel                      string        `json:"analytics_access_level"`
// 	ContainerRegistryAccessLevel              string        `json:"container_registry_access_level"`
// 	EmailsDisabled                            interface{}   `json:"emails_disabled"`
// 	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
// 	LfsEnabled                                bool          `json:"lfs_enabled"`
// 	CreatorID                                 int           `json:"creator_id"`
// 	ImportStatus                              string        `json:"import_status"`
// 	OpenIssuesCount                           int           `json:"open_issues_count"`
// 	CiDefaultGitDepth                         int           `json:"ci_default_git_depth"`
// 	CiForwardDeploymentEnabled                bool          `json:"ci_forward_deployment_enabled"`
// 	CiJobTokenScopeEnabled                    bool          `json:"ci_job_token_scope_enabled"`
// 	PublicJobs                                bool          `json:"public_jobs"`
// 	BuildTimeout                              int           `json:"build_timeout"`
// 	AutoCancelPendingPipelines                string        `json:"auto_cancel_pending_pipelines"`
// 	BuildCoverageRegex                        interface{}   `json:"build_coverage_regex"`
// 	CiConfigPath                              string        `json:"ci_config_path"`
// 	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
// 	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
// 	AllowMergeOnSkippedPipeline               interface{}   `json:"allow_merge_on_skipped_pipeline"`
// 	RestrictUserDefinedVariables              bool          `json:"restrict_user_defined_variables"`
// 	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
// 	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
// 	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge"`
// 	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
// 	MergeMethod                               string        `json:"merge_method"`
// 	SquashOption                              string        `json:"squash_option"`
// 	SuggestionCommitMessage                   interface{}   `json:"suggestion_commit_message"`
// 	AutoDevopsEnabled                         bool          `json:"auto_devops_enabled"`
// 	AutoDevopsDeployStrategy                  string        `json:"auto_devops_deploy_strategy"`
// 	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues"`
// 	KeepLatestArtifact                        bool          `json:"keep_latest_artifact"`
// 	ApprovalsBeforeMerge                      int           `json:"approvals_before_merge"`
// 	Mirror                                    bool          `json:"mirror"`
// 	ExternalAuthorizationClassificationLabel  string        `json:"external_authorization_classification_label"`
// 	MarkedForDeletionAt                       interface{}   `json:"marked_for_deletion_at"`
// 	MarkedForDeletionOn                       interface{}   `json:"marked_for_deletion_on"`
// 	RequirementsEnabled                       bool          `json:"requirements_enabled"`
// 	SecurityAndComplianceEnabled              bool          `json:"security_and_compliance_enabled"`
// 	ComplianceFrameworks                      []interface{} `json:"compliance_frameworks"`
// 	IssuesTemplate                            interface{}   `json:"issues_template"`
// 	MergeRequestsTemplate                     interface{}   `json:"merge_requests_template"`
// 	MergePipelinesEnabled                     bool          `json:"merge_pipelines_enabled"`
// 	MergeTrainsEnabled                        bool          `json:"merge_trains_enabled"`
// 	Permissions                               struct {
// 		ProjectAccess interface{} `json:"project_access"`
// 		GroupAccess   interface{} `json:"group_access"`
// 	} `json:"permissions"`
// 	Owner struct {
// 		ID        int    `json:"id"`
// 		Name      string `json:"name"`
// 		Username  string `json:"username"`
// 		State     string `json:"state"`
// 		AvatarURL string `json:"avatar_url"`
// 		WebURL    string `json:"web_url"`
// 	} `json:"owner,omitempty"`
// 	ForkedFromProject struct {
// 		ID                int           `json:"id"`
// 		Description       string        `json:"description"`
// 		Name              string        `json:"name"`
// 		NameWithNamespace string        `json:"name_with_namespace"`
// 		Path              string        `json:"path"`
// 		PathWithNamespace string        `json:"path_with_namespace"`
// 		CreatedAt         time.Time     `json:"created_at"`
// 		DefaultBranch     string        `json:"default_branch"`
// 		TagList           []interface{} `json:"tag_list"`
// 		Topics            []interface{} `json:"topics"`
// 		SSHURLToRepo      string        `json:"ssh_url_to_repo"`
// 		HTTPURLToRepo     string        `json:"http_url_to_repo"`
// 		WebURL            string        `json:"web_url"`
// 		ReadmeURL         string        `json:"readme_url"`
// 		AvatarURL         string        `json:"avatar_url"`
// 		ForksCount        int           `json:"forks_count"`
// 		StarCount         int           `json:"star_count"`
// 		LastActivityAt    time.Time     `json:"last_activity_at"`
// 		Namespace         struct {
// 			ID        int         `json:"id"`
// 			Name      string      `json:"name"`
// 			Path      string      `json:"path"`
// 			Kind      string      `json:"kind"`
// 			FullPath  string      `json:"full_path"`
// 			ParentID  interface{} `json:"parent_id"`
// 			AvatarURL interface{} `json:"avatar_url"`
// 			WebURL    string      `json:"web_url"`
// 		} `json:"namespace"`
// 	} `json:"forked_from_project,omitempty"`
// }

// getProjectCmd represents the projects command
var getProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		showOpts := 0
		glUser, _ := cmd.Flags().GetString("user")
		grID, _ := cmd.Flags().GetInt("group-id")
		showPipeline, _ := cmd.Flags().GetBool("pipeline")
		showMergeRequest, _ := cmd.Flags().GetBool("merge-request")
		if !showPipeline && !showMergeRequest {
			showOpts = 1
		} else {
			if showPipeline {
				showOpts = showOpts + 1
			}
			if showMergeRequest {
				showOpts = showOpts + 2
			}
		}

		getProject(grID, glUser, showOpts)
	},
}

func getProject(id int, user string, opts int) error {
	restClient := resty.New()

	uri := fmt.Sprintf("https://%s/api/v4/groups/%d/projects", glHost, id)

	resp, resperr := restClient.R().
		SetHeader("PRIVATE-TOKEN", glToken).
		Get(uri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
	}

	var pr Project
	marshErr := json.Unmarshal(resp.Body(), &pr)
	if marshErr != nil {
		logrus.Fatal("Cannot marshall Project", marshErr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	switch opts {
	case 1:
		table.SetHeader([]string{"ID", "NAME", "PATH", "PIPELINES"})
	case 2:
		table.SetHeader([]string{"ID", "NAME", "PATH", "MR"})
	case 3:
		table.SetHeader([]string{"ID", "NAME", "PATH", "PIPELINE / MR"})
	}
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
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

	for _, v := range pr {

		var cmdLine string
		var mrLine string
		if len(user) > 0 {
			cmdLine = fmt.Sprintf("<bash:gitlab-tool get pipeline -p %d -u %s>", v.ID, user)
		} else {
			cmdLine = fmt.Sprintf("<bash:gitlab-tool get pipeline -p %d>", v.ID)
		}
		mrLine = fmt.Sprintf("<bash:gitlab-tool get mr -p %d>", v.ID)

		row := []string{}
		switch opts {
		case 1:
			row = []string{
				fmt.Sprintf("%d", v.ID),
				v.Name,
				v.Path,
				cmdLine,
			}
		case 2:
			row = []string{
				fmt.Sprintf("%d", v.ID),
				v.Name,
				v.Path,
				mrLine,
			}
		case 3:
			row = []string{
				fmt.Sprintf("%d", v.ID),
				v.Name,
				v.Path,
				fmt.Sprintf("%s\n%s", cmdLine, mrLine),
			}
		}
		table.Append(row)
	}
	table.Render()

	// fmt.Println(string(resp.Body()[:]))

	return nil
}

func init() {
	getCmd.AddCommand(getProjectCmd)

	getProjectCmd.Flags().StringP("user", "u", "", "Specify the gitlab User")
	getProjectCmd.Flags().IntP("group-id", "g", 0, "Specify the GroupID")
	getProjectCmd.Flags().BoolP("pipeline", "p", false, "Show PipeLine Links")
	getProjectCmd.Flags().BoolP("merge-request", "m", false, "Show Merge Request Links")
}
