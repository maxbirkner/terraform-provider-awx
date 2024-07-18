package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceJobTemplateInstanceGroups() *schema.Resource {
	return &schema.Resource{
		Description:   "Associates an instance group to a job template",
		CreateContext: resourceJobTemplateInstanceGroupsCreate,
		DeleteContext: resourceJobTemplateInstanceGroupsDelete,
		ReadContext:   resourceJobTemplateInstanceGroupsRead,

		Schema: map[string]*schema.Schema{

			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the job template to associate the instance group with",
			},
			"instance_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance group to associate with the job template",
			},
		},
	}
}

func resourceJobTemplateInstanceGroupsCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)
	if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
		return utils.DiagNotFound("JobTemplate Credential", jobTemplateID, err)
	}

	result, err := client.JobTemplateService.AssociateInstanceGroups(jobTemplateID, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{})

	if err != nil {
		return utils.DiagCreate("JobTemplate AssociateInstanceGroups", err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return nil
}

func resourceJobTemplateInstanceGroupsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceJobTemplateInstanceGroupsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)
	res, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("JobTemplate Credential", jobTemplateID, err)
	}

	if _, err = client.JobTemplateService.DisAssociateInstanceGroups(res.ID, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{}); err != nil {
		return utils.DiagDelete("JobTemplate DisAssociateInstanceGroups", jobTemplateID, err)
	}

	d.SetId("")
	return nil
}
