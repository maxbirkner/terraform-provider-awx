package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

//nolint:funlen
func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_project` manages projects within an organization.",
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
		UpdateContext: resourceProjectUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this project",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this project.",
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.",
			},

			"scm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of \"\" (manual), git, hg, svn",
			},

			"scm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "",
			},
			"scm_credential_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Numeric ID of the scm used credential",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specific branch, tag or commit to checkout.",
			},
			"scm_clean": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scm_delete_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric ID of the project organization",
			},
			"scm_update_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scm_update_cache_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"allow_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow SCM branch override",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	orgID := d.Get("organization_id").(int)
	projectName := d.Get("name").(string)
	_, res, err := client.ProjectService.ListProjects(map[string]string{
		"name":         projectName,
		"organization": strconv.Itoa(orgID),
	},
	)
	if err != nil {
		return utils.DiagFetch(diagProjectTitle, orgID, err)
	}
	if len(res.Results) >= 1 {
		return utils.Diagf("Create: Always exist", "Project with name %s  already exists in the Organization ID %v", projectName, orgID)
	}
	credentials := ""
	if d.Get("scm_credential_id").(int) > 0 {
		credentials = strconv.Itoa(d.Get("scm_credential_id").(int))
	}
	result, err := client.ProjectService.CreateProject(map[string]interface{}{
		"name":                 projectName,
		"description":          d.Get("description").(string),
		"local_path":           d.Get("local_path").(string),
		"scm_type":             d.Get("scm_type").(string),
		"scm_url":              d.Get("scm_url").(string),
		"scm_branch":           d.Get("scm_branch").(string),
		"scm_clean":            d.Get("scm_clean").(bool),
		"scm_delete_on_update": d.Get("scm_delete_on_update").(bool),
		"organization":         d.Get("organization_id").(int),
		"credential":           credentials,

		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
		"allow_override":           d.Get("allow_override").(bool),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagProjectTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update Project", d)
	if diags.HasError() {
		return diags
	}
	credentials := ""
	if d.Get("scm_credential_id").(int) > 0 {
		credentials = strconv.Itoa(d.Get("scm_credential_id").(int))
	}

	data := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"scm_type":                 d.Get("scm_type").(string),
		"scm_url":                  d.Get("scm_url").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"scm_clean":                d.Get("scm_clean").(bool),
		"scm_delete_on_update":     d.Get("scm_delete_on_update").(bool),
		"credential":               credentials,
		"organization":             d.Get("organization_id").(int),
		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
		"allow_override":           d.Get("allow_override").(bool),
	}

	// Cannot change local_path for git-based projects
	if d.Get("local_path").(string) != "" && d.Get("scm_type").(string) != "git" {
		data["local_path"] = d.Get("local_path").(string)
	}

	if _, err := client.ProjectService.UpdateProject(id, data, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagProjectTitle, id, err)
	}
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read Project", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.ProjectService.GetProjectByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagProjectTitle, id, err)
	}
	d = setProjectResourceData(d, res)
	return diags
}

func resourceProjectDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	var jobID int
	var finished time.Time
	id, diags := utils.StateIDToInt("Delete Project", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.ProjectService.GetProjectByID(id, make(map[string]string))
	if err != nil {
		d.SetId("")
		return utils.DiagNotFound(diagProjectTitle, id, err)
	}

	if res.SummaryFields.CurrentJob["id"] != nil {
		jobID = int(res.SummaryFields.CurrentJob["id"].(float64))
	} else if res.SummaryFields.LastJob["id"] != nil {
		jobID = int(res.SummaryFields.LastJob["id"].(float64))
	}
	if jobID != 0 {
		if _, err = client.ProjectUpdatesService.ProjectUpdateCancel(jobID); err != nil {
			return utils.Diagf(
				"Delete: Fail to canel Job",
				"Fail to canel the Job %v for Project with ID %v, got %s",
				jobID, id, err,
			)
		}
	}
	// check if finished is 0
	for finished.IsZero() {
		prj, _ := client.ProjectUpdatesService.ProjectUpdateGet(jobID)
		finished = prj.Finished
		time.Sleep(1 * time.Second)
	}

	if _, err = client.ProjectService.DeleteProject(id); err != nil {
		return utils.DiagDelete(diagProjectTitle, id, err)
	}
	d.SetId("")
	return diags
}

func setProjectResourceData(d *schema.ResourceData, r *awx.Project) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("scm_type", r.ScmType); err != nil {
		fmt.Println("Error setting scm_type", err)
	}
	if err := d.Set("scm_url", r.ScmURL); err != nil {
		fmt.Println("Error setting scm_url", err)
	}
	if err := d.Set("scm_branch", r.ScmBranch); err != nil {
		fmt.Println("Error setting scm_branch", err)
	}
	if err := d.Set("scm_clean", r.ScmClean); err != nil {
		fmt.Println("Error setting scm_clean", err)
	}
	if err := d.Set("scm_delete_on_update", r.ScmDeleteOnUpdate); err != nil {
		fmt.Println("Error setting scm_delete_on_update", err)
	}
	if err := d.Set("organization_id", r.Organization); err != nil {
		fmt.Println("Error setting organization_id", err)
	}
	if err := d.Set("scm_credential_id", r.Credential); err != nil {
		fmt.Println("Error setting scm_credential_id", err)
	}
	if err := d.Set("scm_update_on_launch", r.ScmUpdateOnLaunch); err != nil {
		fmt.Println("Error setting scm_update_on_launch", err)
	}
	if err := d.Set("scm_update_cache_timeout", r.ScmUpdateCacheTimeout); err != nil {
		fmt.Println("Error setting scm_update_cache_timeout", err)
	}
	if err := d.Set("allow_override", r.AllowOverride); err != nil {
		fmt.Println("Error setting allow_override", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
