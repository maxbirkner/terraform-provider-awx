package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

//nolint:funlen
func resourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_job_template` manages job templates within AWX.",
		CreateContext: resourceJobTemplateCreate,
		ReadContext:   resourceJobTemplateRead,
		UpdateContext: resourceJobTemplateUpdate,
		DeleteContext: resourceJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the job template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the job template.",
			},
			// Run, Check, Scan
			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Can be one of: `run`, `check`, or `scan`",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The inventory ID to associate with the job template.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The project ID to associate with the job template.",
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The playbook to associate with the job template.",
			},
			"forks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of forks to associate with the job template.",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The limit to apply to filter hosts that run on this job template.",
			},
			//0,1,2,3,4,5
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The extra variables to associate with the job template.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The job tags to associate with the job template.",
			},
			"force_handlers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force handlers to run on the job template.",
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The tags to skip on the job template.",
			},
			"start_at_task": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The task to start at on the job template.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The timeout to associate with the job template. Default is 0",
			},
			"use_fact_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use the fact cache on the job template.",
			},
			"host_config_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_diff_mode_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_verbosity_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_credential_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"become_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_job_type_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"execution_environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The selected execution environment that this playbook will be run in.",
			},
		},
	}
}

func resourceJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService

	result, err := awxService.CreateJobTemplate(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"project":                  d.Get("project_id").(int),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        utils.AtoiDefault(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":    utils.AtoiDefault(d.Get("execution_environment").(string), nil),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d, failed to create %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Update JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetJobTemplateByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)
	}

	_, err = awxService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"project":                  d.Get("project_id").(int),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        utils.AtoiDefault(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":    utils.AtoiDefault(d.Get("execution_environment").(string), nil),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d failed to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Read JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)

	}
	d = setJobTemplateResourceData(d, res)
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awx.JobTemplate) *schema.ResourceData {
	if err := d.Set("allow_simultaneous", r.AllowSimultaneous); err != nil {
		return d
	}
	if err := d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch); err != nil {
		return d
	}
	if err := d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch); err != nil {
		return d
	}
	if err := d.Set("ask_limit_on_launch", r.AskLimitOnLaunch); err != nil {
		return d
	}
	if err := d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch); err != nil {
		return d
	}
	if err := d.Set("ask_tags_on_launch", r.AskTagsOnLaunch); err != nil {
		return d
	}
	if err := d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch); err != nil {
		return d
	}
	if err := d.Set("description", r.Description); err != nil {
		return d
	}
	if err := d.Set("extra_vars", normalizeJsonYaml(r.ExtraVars)); err != nil {
		return d
	}
	if err := d.Set("force_handlers", r.ForceHandlers); err != nil {
		return d
	}
	if err := d.Set("forks", r.Forks); err != nil {
		return d
	}
	if err := d.Set("host_config_key", r.HostConfigKey); err != nil {
		return d
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		return d
	}
	if err := d.Set("job_tags", r.JobTags); err != nil {
		return d
	}
	if err := d.Set("job_type", r.JobType); err != nil {
		return d
	}
	if err := d.Set("diff_mode", r.DiffMode); err != nil {
		return d
	}
	if err := d.Set("custom_virtualenv", r.CustomVirtualenv); err != nil {
		return d
	}
	if err := d.Set("limit", r.Limit); err != nil {
		return d
	}
	if err := d.Set("name", r.Name); err != nil {
		return d
	}
	if err := d.Set("become_enabled", r.BecomeEnabled); err != nil {
		return d
	}
	if err := d.Set("use_fact_cache", r.UseFactCache); err != nil {
		return d
	}
	if err := d.Set("playbook", r.Playbook); err != nil {
		return d
	}
	if err := d.Set("project_id", r.Project); err != nil {
		return d
	}
	if err := d.Set("skip_tags", r.SkipTags); err != nil {
		return d
	}
	if err := d.Set("start_at_task", r.StartAtTask); err != nil {
		return d
	}
	if err := d.Set("survey_enabled", r.SurveyEnabled); err != nil {
		return d
	}
	if err := d.Set("verbosity", r.Verbosity); err != nil {
		return d
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
