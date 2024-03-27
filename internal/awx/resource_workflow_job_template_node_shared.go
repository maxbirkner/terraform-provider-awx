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

var workflowJobNodeSchema = map[string]*schema.Schema{

	"extra_data": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "",
		StateFunc:   normalizeJsonYaml,
	},
	"workflow_job_template_node_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The workflow_job_template_node id from with the new node will start",
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
		Description: "The SCM branch to use for the job template.",
	},
	"job_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "run",
		Description: "The type of job to run.",
	},
	"job_tags": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A list of job tags to use for the job template.",
	},
	"skip_tags": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A list of job tags to skip for the job template.",
	},
	"limit": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A host pattern to limit the job template to.",
	},
	"diff_mode": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Whether to enable diff mode for the job template.",
	},
	"verbosity": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "The verbosity level for the job template. Can be one of 0, 1, 2, 3, or 4.",
	},
	"workflow_job_template_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The workflow job template id to which the node belongs",
	},
	"unified_job_template_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The unified job template id to which the node belongs",
	},
	"all_parents_must_converge": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Whether all parents must converge before this node can start",
	},
	"identifier": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The identifier for the node",
	},
}

func createNodeForWorkflowJob(awxService *awx.WorkflowJobTemplateNodeStepService, ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	templateNodeID := d.Get("workflow_job_template_node_id").(int)
	result, err := awxService.CreateWorkflowJobTemplateNodeStep(templateNodeID, map[string]interface{}{
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
		//"success_nodes":         d.Get("success_nodes").([]interface{}),
		//"always_nodes":          d.Get("always_nodes").([]interface{}),

		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplateNodeSuccess",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNodeSuccess with JobTemplateID %d failed to create %s", d.Get("unified_job_template_id").(int), err.Error()),
		})
		return diags
	}
	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}
