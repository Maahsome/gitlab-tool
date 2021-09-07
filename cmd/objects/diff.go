package objects

import "time"

type MergeRequestDiff []struct {
	ID             int         `json:"id"`
	Iid            int         `json:"iid"`
	ProjectID      int         `json:"project_id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	State          string      `json:"state"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	MergedBy       interface{} `json:"merged_by"`
	MergedAt       interface{} `json:"merged_at"`
	ClosedBy       interface{} `json:"closed_by"`
	ClosedAt       interface{} `json:"closed_at"`
	TargetBranch   string      `json:"target_branch"`
	SourceBranch   string      `json:"source_branch"`
	UserNotesCount int         `json:"user_notes_count"`
	Upvotes        int         `json:"upvotes"`
	Downvotes      int         `json:"downvotes"`
	Author         struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	Assignees []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignees"`
	Assignee struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignee"`
	Reviewers                 []interface{} `json:"reviewers"`
	SourceProjectID           int           `json:"source_project_id"`
	TargetProjectID           int           `json:"target_project_id"`
	Labels                    []interface{} `json:"labels"`
	Draft                     bool          `json:"draft"`
	WorkInProgress            bool          `json:"work_in_progress"`
	Milestone                 interface{}   `json:"milestone"`
	MergeWhenPipelineSucceeds bool          `json:"merge_when_pipeline_succeeds"`
	MergeStatus               string        `json:"merge_status"`
	Sha                       string        `json:"sha"`
	MergeCommitSha            interface{}   `json:"merge_commit_sha"`
	SquashCommitSha           interface{}   `json:"squash_commit_sha"`
	DiscussionLocked          interface{}   `json:"discussion_locked"`
	ShouldRemoveSourceBranch  interface{}   `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch   bool          `json:"force_remove_source_branch"`
	Reference                 string        `json:"reference"`
	References                struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	WebURL    string `json:"web_url"`
	TimeStats struct {
		TimeEstimate        int         `json:"time_estimate"`
		TotalTimeSpent      int         `json:"total_time_spent"`
		HumanTimeEstimate   interface{} `json:"human_time_estimate"`
		HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Squash               bool `json:"squash"`
	TaskCompletionStatus struct {
		Count          int `json:"count"`
		CompletedCount int `json:"completed_count"`
	} `json:"task_completion_status"`
	HasConflicts                bool        `json:"has_conflicts"`
	BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        interface{} `json:"approvals_before_merge"`
	Subscribed                  bool        `json:"subscribed"`
	ChangesCount                string      `json:"changes_count"`
	LatestBuildStartedAt        time.Time   `json:"latest_build_started_at"`
	LatestBuildFinishedAt       time.Time   `json:"latest_build_finished_at"`
	FirstDeployedToProductionAt interface{} `json:"first_deployed_to_production_at"`
	Pipeline                    struct {
		ID        int       `json:"id"`
		ProjectID int       `json:"project_id"`
		Sha       string    `json:"sha"`
		Ref       string    `json:"ref"`
		Status    string    `json:"status"`
		Source    string    `json:"source"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		WebURL    string    `json:"web_url"`
	} `json:"pipeline"`
	HeadPipeline struct {
		ID         int         `json:"id"`
		ProjectID  int         `json:"project_id"`
		Sha        string      `json:"sha"`
		Ref        string      `json:"ref"`
		Status     string      `json:"status"`
		Source     string      `json:"source"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		WebURL     string      `json:"web_url"`
		BeforeSha  string      `json:"before_sha"`
		Tag        bool        `json:"tag"`
		YamlErrors interface{} `json:"yaml_errors"`
		User       struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		} `json:"user"`
		StartedAt      time.Time   `json:"started_at"`
		FinishedAt     time.Time   `json:"finished_at"`
		CommittedAt    interface{} `json:"committed_at"`
		Duration       int         `json:"duration"`
		QueuedDuration int         `json:"queued_duration"`
		Coverage       interface{} `json:"coverage"`
		DetailedStatus struct {
			Icon         string      `json:"icon"`
			Text         string      `json:"text"`
			Label        string      `json:"label"`
			Group        string      `json:"group"`
			Tooltip      string      `json:"tooltip"`
			HasDetails   bool        `json:"has_details"`
			DetailsPath  string      `json:"details_path"`
			Illustration interface{} `json:"illustration"`
			Favicon      string      `json:"favicon"`
		} `json:"detailed_status"`
	} `json:"head_pipeline"`
	DiffRefs struct {
		BaseSha  string `json:"base_sha"`
		HeadSha  string `json:"head_sha"`
		StartSha string `json:"start_sha"`
	} `json:"diff_refs"`
	MergeError interface{} `json:"merge_error"`
	User       struct {
		CanMerge bool `json:"can_merge"`
	} `json:"user"`
	Changes []struct {
		OldPath     string `json:"old_path"`
		NewPath     string `json:"new_path"`
		AMode       string `json:"a_mode"`
		BMode       string `json:"b_mode"`
		NewFile     bool   `json:"new_file"`
		RenamedFile bool   `json:"renamed_file"`
		DeletedFile bool   `json:"deleted_file"`
		Diff        string `json:"diff"`
	} `json:"changes"`
	Overflow bool `json:"overflow"`
}
