package awx

import (
	"context"
	"strconv"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagElementJobTemplate = "Job Template"

func dataSourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobTemplateRead,
		Description: "Use this data source to get the details of a Job Template",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Job Template",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Job Template",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceJobTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	jobTemplate, _, err := client.JobTemplateService.ListJobTemplates(params)

	if err != nil {
		return utils.DiagFetch(diagElementJobTemplate, params, err)
	}

	for _, template := range jobTemplate {
		log.Printf("loop %v", template.Name)
		if template.Name == params["name"] {
			d = setJobTemplateResourceData(d, template)
			return diags
		}
	}

	if _, okGroupID := d.GetOk("id"); okGroupID {
		log.Printf("byid %v", len(jobTemplate))
		if len(jobTemplate) > 1 {
			return utils.Diagf(
				"Get: find more than one Element",
				"The Query Returns more than one Group, %d",
				len(jobTemplate),
			)
		}
		if len(jobTemplate) == 0 {
			return utils.Diagf(
				"Get: Job Template does not exist",
				"The Query Returns no Job Template matching filter %v",
				params,
			)
		}
		d = setJobTemplateResourceData(d, jobTemplate[0])
		return diags
	}
	return utils.Diagf(
		"Get: find more than one Element",
		"The Query Returns more than one Group, %d",
		len(jobTemplate),
	)
}
