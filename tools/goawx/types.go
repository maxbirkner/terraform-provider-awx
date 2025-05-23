package awx

import (
	"encoding/json"
	"time"
)

// Common types definition here
// For common usage, we made `Related` and `Summary` as two common field,
// it maybe happened that some structs don't have some fields in `Related` or `Summary`.

// Pagination represents the awx api pagination params.
type Pagination struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
}

// PaginationRequest : Paged response from AWX.
type PaginationRequest struct {
	AllPages *bool
}

// Application represents the awx api application.
type Application struct {
	Name                   string   `json:"name"`
	ID                     int      `json:"id"`
	Type                   string   `json:"type"`
	URL                    string   `json:"url"`
	Related                *Related `json:"related"`
	Description            string   `json:"description"`
	ClientID               string   `json:"client_id"`
	ClientSecret           string   `json:"client_secret"`
	ClientType             string   `json:"client_type"`
	RedirectURIs           string   `json:"redirect_uris"`
	AuthorizationGrantType string   `json:"authorization_grant_type"`
	SkipAuthorization      bool     `json:"skip_authorization"`
	OrganizationID         int      `json:"organization"`
}

// ProjectUpdateCancel represents the awx project update cancel api response.
type ProjectUpdateCancel struct {
	CanCancel bool `json:"can_cancel"`
}

// Related represents the awx api related field.
type Related struct {
	NamedURL                     string `json:"named_url"`
	CreatedBy                    string `json:"created_by"`
	ModifiedBy                   string `json:"modified_by"`
	JobTemplates                 string `json:"job_templates"`
	VariableData                 string `json:"variable_data"`
	RootGroups                   string `json:"root_groups"`
	ObjectRoles                  string `json:"object_roles"`
	AdHocCommands                string `json:"ad_hoc_commands"`
	Script                       string `json:"script"`
	Tree                         string `json:"tree"`
	AccessList                   string `json:"access_list"`
	ActivityStream               string `json:"activity_stream"`
	InstanceGroups               string `json:"instance_groups"`
	Hosts                        string `json:"hosts"`
	Job                          string `json:"job"`
	Host                         string `json:"host"`
	Groups                       string `json:"groups"`
	Copy                         string `json:"copy"`
	UpdateInventorySources       string `json:"update_inventory_sources"`
	InventorySources             string `json:"inventory_sources"`
	FactVersions                 string `json:"fact_versions"`
	SmartInventories             string `json:"smart_inventories"`
	Insights                     string `json:"insights"`
	Organization                 string `json:"organization"`
	Labels                       string `json:"labels"`
	Inventory                    string `json:"inventory"`
	Project                      string `json:"project"`
	Credential                   string `json:"credential"`
	ExtraCredentials             string `json:"extra_credentials"`
	Credentials                  string `json:"credentials"`
	NotificationTemplatesError   string `json:"notification_templates_error"`
	NotificationTemplatesSuccess string `json:"notification_templates_success"`
	Jobs                         string `json:"jobs"`
	NotificationTemplatesAny     string `json:"notification_templates_any"`
	Launch                       string `json:"launch"`
	Schedules                    string `json:"schedules"`
	SurveySpec                   string `json:"survey_spec"`
	UnifiedJobTemplate           string `json:"unified_job_template"`
	Stdout                       string `json:"stdout"`
	Notifications                string `json:"notifications"`
	JobHostSummaries             string `json:"job_host_summaries"`
	JobEvents                    string `json:"job_events"`
	JobTemplate                  string `json:"job_template"`
	Cancel                       string `json:"cancel"`
	ProjectUpdate                string `json:"project_update"`
	CreateSchedule               string `json:"create_schedule"`
	Relaunch                     string `json:"relaunch"`
	AdminOfOrganizations         string `json:"admin_of_organizations"`
	Organizations                string `json:"organizations"`
	Roles                        string `json:"roles"`
	Teams                        string `json:"teams"`
	Projects                     string `json:"projects"`
	PotentialChildren            string `json:"potential_children"`
	AllHosts                     string `json:"all_hosts"`
	AllGroups                    string `json:"all_groups"`
	AdHocCommandEvents           string `json:"ad_hoc_command_events"`
	Children                     string `json:"children"`
	AnsibleFacts                 string `json:"ansible_facts"`
	ExecutionEnvironment         string `json:"execution_environment"`
	ExecutionEnvironments        string `json:"execution_environments"`
}

// ExecutionEnvironmentSummary represents the awx api instance group summary fields.
type ExecutionEnvironmentSummary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// OrganizationSummary represents the awx api organization summary fields.
type OrganizationSummary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ByUserSummary represents the awx api user summary fields.
type ByUserSummary struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// JobTemplateSummary represents the awx api job template summary fields.
type JobTemplateSummary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// InstanceGroupSummary represents the awx api instance group summary fields.
type InstanceGroupSummary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ApplyRole represents the awx api apply role.
type ApplyRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ObjectRoles represents the awx api object roles.
type ObjectRoles struct {
	AdhocRole                    *ApplyRole `json:"adhoc_role"`
	AdminRole                    *ApplyRole `json:"admin_role"`
	ApprovalRole                 *ApplyRole `json:"approval_role"`
	AuditorRole                  *ApplyRole `json:"auditor_role"`
	CredentialAdminRole          *ApplyRole `json:"credential_admin_role"`
	ExecuteRole                  *ApplyRole `json:"execute_role"`
	ExecuteEnvironmentsAdminRole *ApplyRole `json:"execution_environment_admin_role"`
	InventoryAdminRole           *ApplyRole `json:"inventory_admin_role"`
	JobTemplateAdminRole         *ApplyRole `json:"job_template_admin_role"`
	MemberRole                   *ApplyRole `json:"member_role"`
	NotificationAdminRole        *ApplyRole `json:"notification_admin_role"`
	ProjectAdminRole             *ApplyRole `json:"project_admin_role"`
	ReadRole                     *ApplyRole `json:"read_role"`
	UpdateRole                   *ApplyRole `json:"update_role"`
	UseRole                      *ApplyRole `json:"use_role"`
	WorkflowAdminRole            *ApplyRole `json:"workflow_admin_role"`
}

// UserCapabilities represents the awx api user capabilities.
type UserCapabilities struct {
	Edit     bool `json:"edit"`
	Start    bool `json:"start"`
	Schedule bool `json:"schedule"`
	Copy     bool `json:"copy"`
	Adhoc    bool `json:"adhoc"`
	Delete   bool `json:"delete"`
}

// Labels represents the awx api labels.
type Labels struct {
	Count   int           `json:"count"`
	Results []interface{} `json:"results"`
}

// Summary represents the awx api summary fields.
type Summary struct {
	InstanceGroup               *InstanceGroupSummary        `json:"instance_group"`
	Organization                *OrganizationSummary         `json:"organization"`
	CreatedBy                   *ByUserSummary               `json:"created_by"`
	ModifiedBy                  *ByUserSummary               `json:"modified_by"`
	ObjectRoles                 *ObjectRoles                 `json:"object_roles"`
	UserCapabilities            *UserCapabilities            `json:"user_capabilities"`
	Project                     *Project                     `json:"project"`
	LastJob                     map[string]interface{}       `json:"last_job"`
	CurrentJob                  map[string]interface{}       `json:"current_job"`
	LastUpdate                  map[string]interface{}       `json:"last_update"`
	Inventory                   *Inventory                   `json:"inventory"`
	RecentJobs                  []interface{}                `json:"recent_jobs"`
	Groups                      *Groups                      `json:"groups"`
	Credentials                 []Credential                 `json:"credentials"`
	Credential                  *Credential                  `json:"credential"`
	Labels                      *Labels                      `json:"labels"`
	JobTemplate                 *JobTemplateSummary          `json:"job_template"`
	UnifiedJobTemplate          *UnifiedJobTemplate          `json:"unified_job_template"`
	ExtraCredentials            []interface{}                `json:"extra_credentials"`
	ProjectUpdate               *ProjectUpdate               `json:"project_update"`
	ExecutionEnvironmentSummary *ExecutionEnvironmentSummary `json:"execution_environment"`
}

// ProjectUpdate represents the awx api project update.
type ProjectUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Failed      bool   `json:"failed"`
}

// Project represents the awx api project.
//
//nolint:maligned
type Project struct {
	ID                    int       `json:"id"`
	Type                  string    `json:"type"`
	URL                   string    `json:"url"`
	Related               *Related  `json:"related"`
	SummaryFields         *Summary  `json:"summary_fields"`
	Created               time.Time `json:"created"`
	Modified              time.Time `json:"modified"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	LocalPath             string    `json:"local_path"`
	ScmType               string    `json:"scm_type"`
	ScmURL                string    `json:"scm_url"`
	ScmBranch             string    `json:"scm_branch"`
	ScmClean              bool      `json:"scm_clean"`
	ScmDeleteOnUpdate     bool      `json:"scm_delete_on_update"`
	Credential            int       `json:"credential"`
	Timeout               int       `json:"timeout"`
	LastJobRun            time.Time `json:"last_job_run"`
	LastJobFailed         bool      `json:"last_job_failed"`
	NextJobRun            time.Time `json:"next_job_run"`
	Status                string    `json:"status"`
	Organization          int       `json:"organization"`
	ScmDeleteOnNextUpdate bool      `json:"scm_delete_on_next_update"`
	ScmUpdateOnLaunch     bool      `json:"scm_update_on_launch"`
	ScmUpdateCacheTimeout int       `json:"scm_update_cache_timeout"`
	AllowOverride         bool      `json:"allow_override"`
	ScmRevision           string    `json:"scm_revision"`
	LastUpdateFailed      bool      `json:"last_update_failed"`
	LastUpdated           time.Time `json:"last_updated"`
}

// Inventory represents the awx api inventory.
//
//nolint:maligned
type Inventory struct {
	ID                           int         `json:"id"`
	Type                         string      `json:"type"`
	URL                          string      `json:"url"`
	Related                      *Related    `json:"related"`
	SummaryFields                *Summary    `json:"summary_fields"`
	Created                      time.Time   `json:"created"`
	Modified                     time.Time   `json:"modified"`
	Name                         string      `json:"name"`
	Description                  string      `json:"description"`
	Organization                 int         `json:"organization"`
	OrganizationID               int         `json:"organization_id"`
	Kind                         string      `json:"kind"`
	HostFilter                   interface{} `json:"host_filter"`
	Variables                    string      `json:"variables"`
	HasActiveFailures            bool        `json:"has_active_failures"`
	TotalHosts                   int         `json:"total_hosts"`
	HostsWithActiveFailures      int         `json:"hosts_with_active_failures"`
	TotalGroups                  int         `json:"total_groups"`
	GroupsWithActiveFailures     int         `json:"groups_with_active_failures"`
	HasInventorySources          bool        `json:"has_inventory_sources"`
	TotalInventorySources        int         `json:"total_inventory_sources"`
	InventorySourcesWithFailures int         `json:"inventory_sources_with_failures"`
	InsightsCredential           interface{} `json:"insights_credential"`
	PendingDeletion              bool        `json:"pending_deletion"`
}

// Credential represents the awx api credential.
type Credential struct {
	Description      string                 `json:"description"`
	ID               int                    `json:"id"`
	Kind             string                 `json:"kind"`
	Name             string                 `json:"name"`
	OrganizationID   int                    `json:"organization"`
	CredentialTypeID int                    `json:"credential_type"`
	Inputs           map[string]interface{} `json:"inputs"`
	SummaryFields    *Summary               `json:"summary_fields"`
}

// CredentialType represents the awx api credential type.
type CredentialType struct {
	ID          int         `json:"ID"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Kind        string      `json:"kind"`
	Inputs      interface{} `json:"inputs"`
	Injectors   interface{} `json:"injectors"`
}

// CredentialInputSource represents the awx api input source.
type CredentialInputSource struct {
	ID               int                    `json:"id"`
	Description      string                 `json:"description"`
	TargetCredential int                    `json:"target_credential"`
	SourceCredential int                    `json:"source_credential"`
	InputFieldName   string                 `json:"input_field_name"`
	SummaryFields    *Summary               `json:"summary_fields"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// UnifiedJobTemplate represents the awx api unified job template.
type UnifiedJobTemplate struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	UnifiedJobType string `json:"unified_job_type"`
}

// InstanceGroup represents the awx api instance group.
type InstanceGroup struct {
	ID               int    `json:"id"`
	Capacity         int    `json:"capacity"`
	CredentialID     int    `json:"credential"` //nolint:golint,stylecheck
	Name             string `json:"name"`
	IsContainerGroup bool   `json:"is_container_group"`
	PodSpecOverride  string `json:"pod_spec_override"`
}

// Result data type.
type Result struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Groups represents the awx api hosts group list.
type Groups struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

// Instance represents the awx api instance.
type Instance struct {
	Node      string    `json:"node"`
	Heartbeat time.Time `json:"heartbeat"`
	Version   string    `json:"version"`
	Capacity  int       `json:"capacity"`
}

// Ping represents the awx api ping.
type Ping struct {
	Instances      []Instance      `json:"instances"`
	InstanceGroups []InstanceGroup `json:"instance_groups"`
	Ha             bool            `json:"ha"`
	Version        string          `json:"version"`
	ActiveNode     string          `json:"active_node"`
}

// JobTemplate represents the awx api job template.
//
//nolint:maligned
type JobTemplate struct {
	ID                              int         `json:"id"`
	Type                            string      `json:"type"`
	URL                             string      `json:"url"`
	Related                         *Related    `json:"related"`
	SummaryFields                   *Summary    `json:"summary_fields"`
	Created                         time.Time   `json:"created"`
	Modified                        time.Time   `json:"modified"`
	Name                            string      `json:"name"`
	Description                     string      `json:"description"`
	JobType                         string      `json:"job_type"`
	Inventory                       int         `json:"inventory"`
	Project                         int         `json:"project"`
	Playbook                        string      `json:"playbook"`
	Forks                           int         `json:"forks"`
	JobSliceCount                   int         `json:"job_slice_count"`
	Limit                           string      `json:"limit"`
	Verbosity                       int         `json:"verbosity"`
	ExtraVars                       string      `json:"extra_vars"`
	JobTags                         string      `json:"job_tags"`
	ForceHandlers                   bool        `json:"force_handlers"`
	SkipTags                        string      `json:"skip_tags"`
	StartAtTask                     string      `json:"start_at_task"`
	Timeout                         int         `json:"timeout"`
	UseFactCache                    bool        `json:"use_fact_cache"`
	LastJobRun                      interface{} `json:"last_job_run"`
	LastJobFailed                   bool        `json:"last_job_failed"`
	NextJobRun                      interface{} `json:"next_job_run"`
	Status                          string      `json:"status"`
	HostConfigKey                   string      `json:"host_config_key"`
	AskScmBranchOnLaunch            bool        `json:"ask_scm_branch_on_launch"`
	AskDiffModeOnLaunch             bool        `json:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch            bool        `json:"ask_variables_on_launch"`
	AskLimitOnLaunch                bool        `json:"ask_limit_on_launch"`
	AskTagsOnLaunch                 bool        `json:"ask_tags_on_launch"`
	AskSkipTagsOnLaunch             bool        `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch              bool        `json:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch            bool        `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch            bool        `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch           bool        `json:"ask_credential_on_launch"`
	AskExecutionEnvironmentOnLaunch bool        `json:"ask_execution_environment_on_launch"`
	AskLabelsOnLaunch               bool        `json:"ask_labels_on_launch"`
	AskForksOnLaunch                bool        `json:"ask_forks_on_launch"`
	AskJobSliceCountOnLaunch        bool        `json:"ask_job_slice_count_on_launch"`
	AskTimeoutOnLaunch              bool        `json:"ask_timeout_on_launch"`
	AskInstanceGroupsOnLaunch       bool        `json:"ask_instance_groups_on_launch"`
	SurveyEnabled                   bool        `json:"survey_enabled"`
	BecomeEnabled                   bool        `json:"become_enabled"`
	DiffMode                        bool        `json:"diff_mode"`
	AllowSimultaneous               bool        `json:"allow_simultaneous"`
	CustomVirtualenv                interface{} `json:"custom_virtualenv"`
	Credential                      int         `json:"credential"`
	VaultCredential                 interface{} `json:"vault_credential"`
	ExecutionEnvironment            int         `json:"execution_environment"`
}

// JobLaunch represents the awx api job launch.
//
//nolint:maligned
type JobLaunch struct {
	Job                     int               `json:"job"`
	IgnoredFields           map[string]string `json:"ignored_fields"`
	ID                      int               `json:"id"`
	Type                    string            `json:"type"`
	URL                     string            `json:"url"`
	Related                 *Related          `json:"related"`
	SummaryFields           *Summary          `json:"summary_fields"`
	Created                 time.Time         `json:"created"`
	Modified                time.Time         `json:"modified"`
	Name                    string            `json:"name"`
	Description             string            `json:"description"`
	JobType                 string            `json:"job_type"`
	Inventory               int               `json:"inventory"`
	Project                 int               `json:"project"`
	Playbook                string            `json:"playbook"`
	Forks                   int               `json:"forks"`
	Limit                   string            `json:"limit"`
	Verbosity               int               `json:"verbosity"`
	ExtraVars               string            `json:"extra_vars"`
	JobTags                 string            `json:"job_tags"`
	ForceHandlers           bool              `json:"force_handlers"`
	SkipTags                string            `json:"skip_tags"`
	StartAtTask             string            `json:"start_at_task"`
	Timeout                 int               `json:"timeout"`
	UseFactCache            bool              `json:"use_fact_cache"`
	UnifiedJobTemplate      int               `json:"unified_job_template"`
	LaunchType              string            `json:"launch_type"`
	Status                  string            `json:"status"`
	Failed                  bool              `json:"failed"`
	Started                 interface{}       `json:"started"`
	Finished                interface{}       `json:"finished"`
	Elapsed                 float64           `json:"elapsed"`
	JobArgs                 string            `json:"job_args"`
	JobCwd                  string            `json:"job_cwd"`
	JobEnv                  map[string]string `json:"job_env"`
	JobExplanation          string            `json:"job_explanation"`
	ExecutionNode           string            `json:"execution_node"`
	ResultTraceback         string            `json:"result_traceback"`
	EventProcessingFinished bool              `json:"event_processing_finished"`
	JobTemplate             int               `json:"job_template"`
	PasswordsNeededToStart  []interface{}     `json:"passwords_needed_to_start"`
	AskDiffModeOnLaunch     bool              `json:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch    bool              `json:"ask_variables_on_launch"`
	AskLimitOnLaunch        bool              `json:"ask_limit_on_launch"`
	AskTagsOnLaunch         bool              `json:"ask_tags_on_launch"`
	AskScmBranchOnLaunch    bool              `json:"ask_scm_branch_on_launch"`
	AskSkipTagsOnLaunch     bool              `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch      bool              `json:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch    bool              `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch    bool              `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch   bool              `json:"ask_credential_on_launch"`
	AllowSimultaneous       bool              `json:"allow_simultaneous"`
	Artifacts               map[string]string `json:"artifacts"`
	ScmRevision             string            `json:"scm_revision"`
	InstanceGroup           interface{}       `json:"instance_group"`
	DiffMode                bool              `json:"diff_mode"`
	Credential              int               `json:"credential"`
	VaultCredential         interface{}       `json:"vault_credential"`
}

// Job represents the awx api job.
//
//nolint:maligned
type Job struct {
	ID                      int               `json:"id"`
	Type                    string            `json:"type"`
	URL                     string            `json:"url"`
	Related                 *Related          `json:"related"`
	SummaryFields           *Summary          `json:"summary_fields"`
	Created                 time.Time         `json:"created"`
	Modified                time.Time         `json:"modified"`
	Name                    string            `json:"name"`
	Description             string            `json:"description"`
	JobType                 string            `json:"job_type"`
	Inventory               int               `json:"inventory"`
	Project                 int               `json:"project"`
	Playbook                string            `json:"playbook"`
	Forks                   int               `json:"forks"`
	Limit                   string            `json:"limit"`
	Verbosity               int               `json:"verbosity"`
	ExtraVars               string            `json:"extra_vars"`
	JobTags                 string            `json:"job_tags"`
	ForceHandlers           bool              `json:"force_handlers"`
	SkipTags                string            `json:"skip_tags"`
	StartAtTask             string            `json:"start_at_task"`
	Timeout                 int               `json:"timeout"`
	UseFactCache            bool              `json:"use_fact_cache"`
	UnifiedJobTemplate      int               `json:"unified_job_template"`
	LaunchType              string            `json:"launch_type"`
	Status                  string            `json:"status"`
	Failed                  bool              `json:"failed"`
	Started                 time.Time         `json:"started"`
	Finished                time.Time         `json:"finished"`
	Elapsed                 float64           `json:"elapsed"`
	JobArgs                 string            `json:"job_args"`
	JobCwd                  string            `json:"job_cwd"`
	JobEnv                  map[string]string `json:"job_env"`
	JobExplanation          string            `json:"job_explanation"`
	ExecutionNode           string            `json:"execution_node"`
	ResultTraceback         string            `json:"result_traceback"`
	EventProcessingFinished bool              `json:"event_processing_finished"`
	JobTemplate             int               `json:"job_template"`
	PasswordsNeededToStart  []interface{}     `json:"passwords_needed_to_start"`
	AskDiffModeOnLaunch     bool              `json:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch    bool              `json:"ask_variables_on_launch"`
	AskLimitOnLaunch        bool              `json:"ask_limit_on_launch"`
	AskTagsOnLaunch         bool              `json:"ask_tags_on_launch"`
	AskScmBranchOnLaunch    bool              `json:"ask_scm_branch_on_launch"`
	AskSkipTagsOnLaunch     bool              `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch      bool              `json:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch    bool              `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch    bool              `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch   bool              `json:"ask_credential_on_launch"`
	AllowSimultaneous       bool              `json:"allow_simultaneous"`
	Artifacts               map[string]string `json:"artifacts"`
	ScmRevision             string            `json:"scm_revision"`
	InstanceGroup           int               `json:"instance_group"`
	DiffMode                bool              `json:"diff_mode"`
	Credential              int               `json:"credential"`
	VaultCredential         interface{}       `json:"vault_credential"`
}

// HostSummaryHost represents the awx api host summary host fields.
type HostSummaryHost struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	HasActiveFailures   bool   `json:"has_active_failures"`
	HasInventorySources bool   `json:"has_inventory_sources"`
}

// HostSummaryJob represents the awx api host summary job fields.
type HostSummaryJob struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	Failed          bool    `json:"failed"`
	Elapsed         float64 `json:"elapsed"`
	JobTemplateID   int     `json:"job_template_id"`
	JobTemplateName string  `json:"job_template_name"`
}

// HostSummaryFields represents the awx api host summary fields.
type HostSummaryFields struct {
	Role map[string]string `json:"role"`
	Host *HostSummaryHost  `json:"host"`
	Job  *HostSummaryJob   `json:"job"`
}

// HostSummary represents the awx api host summary.
type HostSummary struct {
	ID            int                `json:"id"`
	Type          string             `json:"type"`
	URL           string             `json:"url"`
	Related       *Related           `json:"related"`
	SummaryFields *HostSummaryFields `json:"summary_fields"`
	Created       time.Time          `json:"created"`
	Modified      time.Time          `json:"modified"`
	Job           int                `json:"job"`
	Host          int                `json:"host"`
	HostName      string             `json:"host_name"`
	Changed       int                `json:"changed"`
	Dark          int                `json:"dark"`
	Failures      int                `json:"failures"`
	Ok            int                `json:"ok"`
	Processed     int                `json:"processed"`
	Skipped       int                `json:"skipped"`
	Failed        bool               `json:"failed"`
}

// EventModuleArgs represents the awx api event module args.
//
//nolint:maligned
type EventModuleArgs struct {
	Creates    interface{} `json:"creates"`
	Executable interface{} `json:"executable"`
	UsesShell  bool        `json:"_uses_shell"`
	RawParams  string      `json:"_raw_params"`
	Removes    interface{} `json:"removes"`
	Warn       bool        `json:"warn"`
	Chdir      string      `json:"chdir"`
	Stdin      interface{} `json:"stdin"`
}

// EventInvocation represents the awx api event invocation.
type EventInvocation struct {
	ModuleArgs *EventModuleArgs `json:"module_args"`
}

// EventRes represents the awx api event response.
//
//nolint:maligned
type EventRes struct {
	AnsibleParsed bool             `json:"_ansible_parsed"`
	StderrLines   []string         `json:"stderr_lines"`
	Changed       bool             `json:"changed"`
	End           string           `json:"end"`
	AnsibleNoLog  bool             `json:"_ansible_no_log"`
	Stdout        string           `json:"stdout"`
	Cmd           string           `json:"cmd"`
	Start         string           `json:"start"`
	Delta         string           `json:"delta"`
	Stderr        string           `json:"stderr"`
	Rc            int              `json:"rc"`
	Invocation    *EventInvocation `json:"invocation"`
	StdoutLines   []string         `json:"stdout_lines"`
	Warnings      []string         `json:"warnings"`
}

// EventData represents the awx api event data.
type EventData struct {
	PlayPattern  string      `json:"play_pattern"`
	Play         string      `json:"play"`
	EventLoop    interface{} `json:"event_loop"`
	TaskArgs     string      `json:"task_args"`
	RemoteAddr   string      `json:"remote_addr"`
	Res          *EventRes   `json:"res"`
	Pid          int         `json:"pid"`
	PlayUUID     string      `json:"play_uuid"`
	TaskUUID     string      `json:"task_uuid"`
	Task         string      `json:"task"`
	PlaybookUUID string      `json:"playbook_uuid"`
	Playbook     string      `json:"playbook"`
	TaskAction   string      `json:"task_action"`
	Host         string      `json:"host"`
	Role         string      `json:"role"`
	TaskPath     string      `json:"task_path"`
}

// JobEvent represents the awx api job event.
type JobEvent struct {
	ID            int                `json:"id"`
	Type          string             `json:"type"`
	URL           string             `json:"url"`
	Related       *Related           `json:"related"`
	SummaryFields *HostSummaryFields `json:"summary_fields"`
	Created       time.Time          `json:"created"`
	Modified      time.Time          `json:"modified"`
	Job           int                `json:"job"`
	Event         string             `json:"event"`
	Counter       int                `json:"counter"`
	EventDisplay  string             `json:"event_display"`
	EventData     *EventData         `json:"event_data"`
	EventLevel    int                `json:"event_level"`
	Failed        bool               `json:"failed"`
	Changed       bool               `json:"changed"`
	UUID          string             `json:"uuid"`
	ParentUUID    string             `json:"parent_uuid"`

	// Host : Host description
	Host interface{} `json:"host"`

	HostName  string      `json:"host_name"`
	Parent    interface{} `json:"parent"`
	Playbook  string      `json:"playbook"`
	Play      string      `json:"play"`
	Task      string      `json:"task"`
	Role      string      `json:"role"`
	Stdout    string      `json:"stdout"`
	StartLine int         `json:"start_line"`
	EndLine   int         `json:"end_line"`
	Verbosity int         `json:"verbosity"`
}

// User represents an user.
type User struct {
	ID              int         `json:"id"`
	Type            string      `json:"type"`
	URL             string      `json:"url"`
	Related         *Related    `json:"related"`
	SummaryFields   *Summary    `json:"summary_fields"`
	Created         time.Time   `json:"created"`
	Username        string      `json:"username"`
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	Email           string      `json:"email"`
	IsSuperUser     bool        `json:"is_superuser"`
	IsSystemAuditor bool        `json:"is_system_auditor"`
	Password        string      `json:"password"`
	LdapDn          string      `json:"ldap_dn"`
	ExternalAccount interface{} `json:"external_account"`
}

// Group represents a group
//
//nolint:maligned
type Group struct {
	ID                       int       `json:"id"`
	Type                     string    `json:"type"`
	URL                      string    `json:"url"`
	Related                  *Related  `json:"related"`
	SummaryFields            *Summary  `json:"summary_fields"`
	Created                  time.Time `json:"created"`
	Modified                 time.Time `json:"modified"`
	Name                     string    `json:"name"`
	Description              string    `json:"description"`
	Inventory                int       `json:"inventory"`
	Variables                string    `json:"variables"`
	HasActiveFailures        bool      `json:"has_active_failures"`
	TotalHosts               int       `json:"total_hosts"`
	HostsWithActiveFailures  int       `json:"hosts_with_active_failures"`
	TotalGroups              int       `json:"total_groups"`
	GroupsWithActiveFailures int       `json:"groups_with_active_failures"`
	HasInventorySources      bool      `json:"has_inventory_sources"`
}

// Host represents a host
//
//nolint:maligned
type Host struct {
	ID                   int         `json:"id"`
	Type                 string      `json:"type"`
	URL                  string      `json:"url"`
	Related              *Related    `json:"related"`
	SummaryFields        *Summary    `json:"summary_fields"`
	Created              time.Time   `json:"created"`
	Modified             time.Time   `json:"modified"`
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	Inventory            int         `json:"inventory"`
	Enabled              bool        `json:"enabled"`
	InstanceID           string      `json:"instance_id"`
	Variables            string      `json:"variables"`
	HasActiveFailures    bool        `json:"has_active_failures"`
	HasInventorySources  bool        `json:"has_inventory_sources"`
	LastJob              int         `json:"last_job"`
	LastJobHostSummary   int         `json:"last_job_host_summary"`
	InsightsSystemID     string      `json:"insights_system_id"`
	AnsibleFactsModified interface{} `json:"ansible_facts_modified"`
}

// Organization represents an organization in AWX.
type Organization struct {
	Created          time.Time `json:"created"`
	CustomVirtualenv string    `json:"custom_virtualenv"`
	Description      string    `json:"description"`
	ID               int       `json:"id"`
	MaxHosts         int       `json:"max_hosts"`
	Modified         time.Time `json:"modified"`
	Name             string    `json:"name"`
	Related          *Related  `json:"related"`
	SummaryFields    *Summary  `json:"summary_fields"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
}

// SettingSummary represents a summary of settings in AWX.
type SettingSummary struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

// Setting abstracted settings pulled from AWX.
type Setting map[string]json.RawMessage

// Team represents a team in AWX.
type Team struct {
	Created       time.Time `json:"created"`
	Description   string    `json:"description"`
	ID            int       `json:"id"`
	Modified      time.Time `json:"modified"`
	Name          string    `json:"name"`
	Organization  int       `json:"organization"`
	Related       *Related  `json:"related"`
	SummaryFields *Summary  `json:"summary_fields"`
	Type          string    `json:"type"`
	URL           string    `json:"url"`
}

// InventorySource represents the awx api inventory source.
//
//nolint:maligned
type InventorySource struct {
	Created               time.Time   `json:"created"`
	Credential            interface{} `json:"credential"`
	CustomVirtualenv      interface{} `json:"custom_virtualenv"`
	Description           string      `json:"description"`
	GroupBy               string      `json:"group_by"`
	ID                    int         `json:"id"`
	EnabledVar            string      `json:"enabled_var"`
	EnabledValue          string      `json:"enabled_value"`
	ExecutionEnvironment  int         `json:"execution_environment"`
	InstanceFilters       string      `json:"instance_filters"`
	HostFilter            string      `json:"host_filter"`
	Inventory             int         `json:"inventory"`
	LastJobFailed         bool        `json:"last_job_failed"`
	LastJobRun            interface{} `json:"last_job_run"`
	LastUpdateFailed      bool        `json:"last_update_failed"`
	LastUpdated           interface{} `json:"last_updated"`
	Modified              time.Time   `json:"modified"`
	Name                  string      `json:"name"`
	NextJobRun            interface{} `json:"next_job_run"`
	Overwrite             bool        `json:"overwrite"`
	OverwriteVars         bool        `json:"overwrite_vars"`
	Related               *Related    `json:"related"`
	Source                string      `json:"source"`
	SourcePath            string      `json:"source_path"`
	SourceProject         int         `json:"source_project"`
	SourceRegions         string      `json:"source_regions"`
	SourceScript          interface{} `json:"source_script"`
	SourceVars            string      `json:"source_vars"`
	Status                string      `json:"status"`
	SummaryFields         *Summary    `json:"summary_fields"`
	Timeout               int         `json:"timeout"`
	Type                  string      `json:"type"`
	UpdateCacheTimeout    int         `json:"update_cache_timeout"`
	UpdateOnLaunch        bool        `json:"update_on_launch"`
	UpdateOnProjectUpdate bool        `json:"update_on_project_update"`
	URL                   string      `json:"url"`
	Verbosity             int         `json:"verbosity"`
}

// WorkflowJobTemplate : represents a workflow job template in AWX.
//
//nolint:maligned
type WorkflowJobTemplate struct {
	ID                   int         `json:"id"`
	Type                 string      `json:"type"`
	URL                  string      `json:"url"`
	Related              *Related    `json:"related"`
	SummaryFields        *Summary    `json:"summary_fields"`
	Created              time.Time   `json:"created"`
	Modified             time.Time   `json:"modified"`
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	LastJobRun           interface{} `json:"last_job_run"`
	LastJobFailed        bool        `json:"last_job_failed"`
	NextJobRun           interface{} `json:"next_job_run"`
	Status               string      `json:"status"`
	ExtraVars            string      `json:"extra_vars"`
	Organization         int         `json:"organization"`
	SurveyEnabled        bool        `json:"survey_enabled"`
	AllowSimultaneous    bool        `json:"allow_simultaneous"`
	AskVariablesOnLaunch bool        `json:"ask_variables_on_launch"`
	Inventory            *int        `json:"inventory"`
	Limit                interface{} `json:"limit"`
	ScmBranch            interface{} `json:"scm_branch"`
	AskInventoryOnLaunch bool        `json:"ask_inventory_on_launch"`
	AskScmBranchOnLaunch bool        `json:"ask_scm_branch_on_launch"`
	AskLimitOnLaunch     bool        `json:"ask_limit_on_launch"`
	WebhookService       string      `json:"webhook_service"`
	WebhookCredential    interface{} `json:"webhook_credential"`
}

// WorkflowJobTemplateNode represents the awx api workflow job template node.
type WorkflowJobTemplateNode struct {
	ID                     int                    `json:"id"`
	Type                   string                 `json:"type"`
	URL                    string                 `json:"url"`
	Related                *Related               `json:"related"`
	SummaryFields          *Summary               `json:"summary_fields"`
	Created                time.Time              `json:"created"`
	Modified               time.Time              `json:"modified"`
	ExtraData              map[string]interface{} `json:"extra_data"`
	Inventory              int                    `json:"inventory"`
	ScmBranch              string                 `json:"scm_branch"`
	JobType                string                 `json:"job_type"`
	JobTags                string                 `json:"job_tags"`
	SkipTags               string                 `json:"skip_tags"`
	Limit                  string                 `json:"limit"`
	DiffMode               bool                   `json:"diff_mode"`
	Verbosity              int                    `json:"verbosity"`
	WorkflowJobTemplate    int                    `json:"workflow_job_template"`
	UnifiedJobTemplate     int                    `json:"unified_job_template"`
	SuccessNodes           []int                  `json:"success_nodes"`
	FailureNodes           []int                  `json:"failure_nodes"`
	AlwaysNodes            []int                  `json:"always_nodes"`
	AllParentsMustConverge bool                   `json:"all_parents_must_converge"`
	Identifier             string                 `json:"identifier"`
}

// Schedule : represents the awx api schedule.
type Schedule struct {
	ID                 int                    `json:"id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Rrule              string                 `json:"rrule"`
	Enabled            bool                   `json:"enabled"`
	UnifiedJobTemplate int                    `json:"unified_job_template"`
	Inventory          int                    `json:"inventory"`
	ExtraData          map[string]interface{} `json:"extra_data"`
}

// NotificationTemplate : represents the awx api notification template.
type NotificationTemplate struct {
	ID                        int                    `json:"id"`
	Name                      string                 `json:"name"`
	Description               string                 `json:"description"`
	Organization              int                    `json:"organization"`
	NotificationType          string                 `json:"notification_type"`
	NotificationConfiguration map[string]interface{} `json:"notification_configuration"`
	Messages                  interface{}            `json:"messages"`
}

// ExecutionEnvironment represents the awx api execution environment summary fields.
type ExecutionEnvironment struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	URL           string    `json:"url"`
	Related       *Related  `json:"related"`
	SummaryFields *Summary  `json:"summary_fields"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Organization  int       `json:"organization"`
	Image         string    `json:"image"`
	Managed       bool      `json:"managed"`
	Credential    int       `json:"credential"`
	Pull          string    `json:"pull"`
}

type SurveySpec struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Spec        []*Spec `json:"spec"`
}

type Spec struct {
	Type                string   `json:"type"`
	QuestionName        string   `json:"question_name"`
	QuestionDescription string   `json:"question_description"`
	Variable            string   `json:"variable"`
	Required            bool     `json:"required"`
	Default             string   `json:"default"`
	Min                 int      `json:"min"`
	Max                 int      `json:"max"`
	Choices             []string `json:"choices"`
}
