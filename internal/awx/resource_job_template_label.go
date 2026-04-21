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

	d.SetId(strconv.Itoa(label.ID))
	return nil
}

func resourceJobTemplateLabelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)

	labelID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: invalid resource ID %q: %w", diagJobTemplateLabelTitle, d.Id(), err))
	}

	labels, err := client.JobTemplateService.ListJobTemplateLabels(jobTemplateID)
	if err != nil {
		return utils.DiagNotFound(diagJobTemplateLabelTitle, jobTemplateID, err)
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

	// Label is no longer associated with this job template — remove from state.
	d.SetId("")
	return nil
}

func resourceJobTemplateLabelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)

	labelID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: invalid resource ID %q: %w", diagJobTemplateLabelTitle, d.Id(), err))
	}

	if err := client.JobTemplateService.DisAssociateLabel(jobTemplateID, labelID); err != nil {
		return utils.DiagDelete(diagJobTemplateLabelTitle, jobTemplateID, err)
	}

	d.SetId("")
	return nil
}
