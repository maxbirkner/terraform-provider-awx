package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagWorkflowJobTemplateLabelTitle = "Workflow Job Template Label"

func resourceWorkflowJobTemplateLabel() *schema.Resource {
	return &schema.Resource{
		Description: "Resource `awx_workflow_job_template_label` creates a label and associates it with a workflow job template. " +
			"AWX will reuse an existing label if one with the same name already exists in the given organization.",
		CreateContext: resourceWorkflowJobTemplateLabelCreate,
		ReadContext:   resourceWorkflowJobTemplateLabelRead,
		DeleteContext: resourceWorkflowJobTemplateLabelDelete,

		Schema: map[string]*schema.Schema{
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the workflow job template to associate the label with.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the label. AWX will find or create a label with this name inside the given organization.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the organization that owns the label.",
			},
		},
	}
}

func resourceWorkflowJobTemplateLabelCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	wjtID := d.Get("workflow_job_template_id").(int)

	if _, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(wjtID, make(map[string]string)); err != nil {
		return utils.DiagNotFound(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	label, err := client.WorkflowJobTemplateService.AssociateLabel(wjtID, map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(int),
	})
	if err != nil {
		return utils.DiagCreate(diagWorkflowJobTemplateLabelTitle, err)
	}

	d.SetId(strconv.Itoa(label.ID))
	return nil
}

func resourceWorkflowJobTemplateLabelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	wjtID := d.Get("workflow_job_template_id").(int)

	labelID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: invalid resource ID %q: %w", diagWorkflowJobTemplateLabelTitle, d.Id(), err))
	}

	labels, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplateLabels(wjtID)
	if err != nil {
		return utils.DiagNotFound(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	for _, label := range labels {
		if label.ID == labelID {
			if err := d.Set("name", label.Name); err != nil {
				return diag.FromErr(fmt.Errorf("error setting name: %w", err))
			}
			if err := d.Set("organization_id", label.Organization); err != nil {
				return diag.FromErr(fmt.Errorf("error setting organization_id: %w", err))
			}
			return nil
		}
	}

	// Label is no longer associated with this workflow job template — remove from state.
	d.SetId("")
	return nil
}

func resourceWorkflowJobTemplateLabelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	wjtID := d.Get("workflow_job_template_id").(int)

	labelID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: invalid resource ID %q: %w", diagWorkflowJobTemplateLabelTitle, d.Id(), err))
	}

	if err := client.WorkflowJobTemplateService.DisAssociateLabel(wjtID, labelID); err != nil {
		return utils.DiagDelete(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	d.SetId("")
	return nil
}
