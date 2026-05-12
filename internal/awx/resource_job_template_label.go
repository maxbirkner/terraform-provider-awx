package awx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagJobTemplateLabelTitle = "Job Template Label"

func resourceJobTemplateLabel() *schema.Resource {
	return &schema.Resource{
		Description: "Resource `awx_job_template_label` creates a label and associates it with a job template. " +
			"AWX will reuse an existing label if one with the same name already exists in the given organization.",
		CreateContext: resourceJobTemplateLabelCreate,
		ReadContext:   resourceJobTemplateLabelRead,
		DeleteContext: resourceJobTemplateLabelDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the job template to associate the label with.",
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

func resourceJobTemplateLabelCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)

	if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
		return utils.DiagNotFound(diagJobTemplateLabelTitle, jobTemplateID, err)
	}

	label, err := client.JobTemplateService.AssociateLabel(jobTemplateID, map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(int),
	})
	if err != nil {
		return utils.DiagCreate(diagJobTemplateLabelTitle, err)
	}

	return syncLabelAssociationCreateState(d, "job_template_id", label)
}

func resourceJobTemplateLabelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)

	name := d.Get("name").(string)
	organizationID := d.Get("organization_id").(int)

	labels, err := client.JobTemplateService.ListJobTemplateLabels(jobTemplateID)
	if err != nil {
		return utils.DiagNotFound(diagJobTemplateLabelTitle, jobTemplateID, err)
	}

	label := findAssociatedLabel(labels, name, organizationID)
	if label != nil {
		return syncLabelAssociationState(d, "job_template_id", label)
	}

	// Label is no longer associated with this job template — remove from state.
	d.SetId("")
	return nil
}

func resourceJobTemplateLabelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)

	if labelID, ok := storedLabelAssociationLabelID(d); ok {
		if err := client.JobTemplateService.DisAssociateLabel(jobTemplateID, labelID); err != nil {
			return utils.DiagDelete(diagJobTemplateLabelTitle, jobTemplateID, err)
		}

		d.SetId("")
		return nil
	}

	name := d.Get("name").(string)
	organizationID := d.Get("organization_id").(int)

	labels, err := client.JobTemplateService.ListJobTemplateLabels(jobTemplateID)
	if err != nil {
		return utils.DiagDelete(diagJobTemplateLabelTitle, jobTemplateID, err)
	}

	label := findAssociatedLabel(labels, name, organizationID)
	if label == nil {
		d.SetId("")
		return nil
	}

	if err := client.JobTemplateService.DisAssociateLabel(jobTemplateID, label.ID); err != nil {
		return utils.DiagDelete(diagJobTemplateLabelTitle, jobTemplateID, err)
	}

	d.SetId("")
	return nil
}
