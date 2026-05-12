package awx

import (
	"context"

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
			"label_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The resolved AWX label ID used internally for efficient deletes.",
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

	return syncLabelAssociationCreateState(d, "workflow_job_template_id", label)
}

func resourceWorkflowJobTemplateLabelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	wjtID := d.Get("workflow_job_template_id").(int)
	name := d.Get("name").(string)
	organizationID := d.Get("organization_id").(int)

	labels, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplateLabels(wjtID)
	if err != nil {
		return utils.DiagNotFound(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	label := findAssociatedLabel(labels, name, organizationID)
	if label != nil {
		return syncLabelAssociationState(d, "workflow_job_template_id", label)
	}

	// Label is no longer associated with this workflow job template — remove from state.
	d.SetId("")
	return nil
}

func resourceWorkflowJobTemplateLabelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	wjtID := d.Get("workflow_job_template_id").(int)

	if labelID, ok := storedLabelAssociationLabelID(d); ok {
		if err := client.WorkflowJobTemplateService.DisAssociateLabel(wjtID, labelID); err != nil {
			return utils.DiagDelete(diagWorkflowJobTemplateLabelTitle, wjtID, err)
		}

		d.SetId("")
		return nil
	}

	name := d.Get("name").(string)
	organizationID := d.Get("organization_id").(int)

	labels, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplateLabels(wjtID)
	if err != nil {
		return utils.DiagDelete(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	label := findAssociatedLabel(labels, name, organizationID)
	if label == nil {
		d.SetId("")
		return nil
	}

	if err := client.WorkflowJobTemplateService.DisAssociateLabel(wjtID, label.ID); err != nil {
		return utils.DiagDelete(diagWorkflowJobTemplateLabelTitle, wjtID, err)
	}

	d.SetId("")
	return nil
}
