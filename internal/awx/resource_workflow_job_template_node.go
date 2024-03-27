package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceWorkflowJobTemplateNode() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource WorkflowJobTemplateNode manages the workflow job template node in AWX.",
		CreateContext: resourceWorkflowJobTemplateNodeCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,

		Schema: map[string]*schema.Schema{

			"extra_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   normalizeJsonYaml,
				Description: "Extra data for the workflow job template node.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory.",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "SCM branch to use for the job template.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "Type of job to run.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags to use for the job template.",
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags to skip for the job template.",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Limit the job template to a specific host or group.",
			},
			"diff_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable diff mode for the job template.",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Verbosity level for the job template. One of 0, 1, 2, 3, 4 or 5.",
			},
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Workflow job template ID to use for the workflow job template node.",
			},
			"unified_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unified job template ID to use for the workflow job template node.",
			},
			//"success_nodes": &schema.Schema{
			//	Type: schema.TypeList,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeInt,
			//	},
			//	Optional: true,
			//},
			//"failure_nodes": &schema.Schema{
			//	Type: schema.TypeList,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeInt,
			//	},
			//	Optional: true,
			//},
			//"always_nodes": &schema.Schema{
			//	Type: schema.TypeList,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeInt,
			//	},
			//	Optional: true,
			//},

			"all_parents_must_converge": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for the workflow job template node.",
			},
		},
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		//
		//Timeouts: &schema.ResourceTimeout{
		//	Create: schema.DefaultTimeout(1 * time.Minute),
		//	Update: schema.DefaultTimeout(1 * time.Minute),
		//	Delete: schema.DefaultTimeout(1 * time.Minute),
		//},
	}
}

func resourceWorkflowJobTemplateNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService

	result, err := awxService.CreateWorkflowJobTemplateNode(map[string]interface{}{
		"extra_data":            d.Get("extra_data").(string),
		"inventory":             d.Get("inventory_id").(int),
		"scm_branch":            d.Get("scm_branch").(string),
		"skip_tags":             d.Get("skip_tags").(string),
		"job_type":              d.Get("job_type").(string),
		"job_tags":              d.Get("job_tags").(string),
		"limit":                 d.Get("limit").(string),
		"diff_mode":             d.Get("diff_mode").(bool),
		"verbosity":             d.Get("verbosity").(int),
		"workflow_job_template": d.Get("workflow_job_template_id").(int),
		"unified_job_template":  d.Get("unified_job_template_id").(int),
		//"failure_nodes":         d.Get("failure_nodes").([]interface{}),
		//"success_nodes": d.Get("success_nodes").([]interface{}),
		//"always_nodes":          d.Get("always_nodes").([]interface{}),

		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplateNode",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNode with JobTemplateID %d and WorkflowID: %d failed to create %s", d.Get("unified_job_template_id").(int), d.Get("workflow_job_template_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNummeric("Update WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetWorkflowJobTemplateNodeByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("workflow job template node", id, err)
	}

	_, err = awxService.UpdateWorkflowJobTemplateNode(id, map[string]interface{}{
		"extra_data":            d.Get("extra_data").(string),
		"inventory":             d.Get("inventory_id").(int),
		"scm_branch":            d.Get("scm_branch").(string),
		"skip_tags":             d.Get("skip_tags").(string),
		"job_type":              d.Get("job_type").(string),
		"job_tags":              d.Get("job_tags").(string),
		"limit":                 d.Get("limit").(string),
		"diff_mode":             d.Get("diff_mode").(bool),
		"verbosity":             d.Get("verbosity").(int),
		"workflow_job_template": d.Get("workflow_job_template_id").(int),
		"unified_job_template":  d.Get("unified_job_template_id").(int),
		//"failure_nodes":             d.Get("failure_nodes").([]interface{}),
		//"success_nodes": d.Get("success_nodes").([]interface{}),
		//"always_nodes":              d.Get("always_nodes").([]interface{}),
		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update WorkflowJobTemplateNode",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNode with name %s in the project id %d failed to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNummeric("Read WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetWorkflowJobTemplateNodeByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("workflow job template node", id, err)

	}
	d = setWorkflowJobTemplateNodeResourceData(d, res)
	return nil
}

func resourceWorkflowJobTemplateNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteWorkflowJobTemplateNode(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setWorkflowJobTemplateNodeResourceData(d *schema.ResourceData, r *awx.WorkflowJobTemplateNode) *schema.ResourceData {

	if err := d.Set("extra_data", normalizeJsonYaml(r.ExtraData)); err != nil {
		return d
	}
	if err := d.Set("inventory_id", strconv.Itoa(r.Inventory)); err != nil {
		return d
	}
	if err := d.Set("scm_branch", r.ScmBranch); err != nil {
		return d
	}
	if err := d.Set("job_type", r.JobType); err != nil {
		return d
	}
	if err := d.Set("job_tags", r.JobTags); err != nil {
		return d
	}
	if err := d.Set("skip_tags", r.SkipTags); err != nil {
		return d
	}
	if err := d.Set("limit", r.Limit); err != nil {
		return d
	}
	if err := d.Set("diff_mode", r.DiffMode); err != nil {
		return d
	}
	if err := d.Set("verbosity", r.Verbosity); err != nil {
		return d
	}
	//d.Set("failure_nodes", r.FailureNodes)
	//d.Set("success_nodes", r.SuccessNodes)
	//d.Set("always_nodes", r.AlwaysNodes)

	if err := d.Set("workflow_job_template_id", strconv.Itoa(r.WorkflowJobTemplate)); err != nil {
		return d
	}
	if err := d.Set("unified_job_template_id", strconv.Itoa(r.UnifiedJobTemplate)); err != nil {
		return d
	}
	if err := d.Set("all_parents_must_converge", r.AllParentsMustConverge); err != nil {
		return d
	}
	if err := d.Set("identifier", r.Identifier); err != nil {
		return d
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
