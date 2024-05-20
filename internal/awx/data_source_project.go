package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagProjectTitle = "Project"

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Description: "This data source allows you to get a project from AWX.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the project.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the project.",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceProjectsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	projects, _, err := client.ProjectService.ListProjects(params)
	if err != nil {
		return utils.DiagFetch(diagProjectTitle, params, err)
	}
	if len(projects) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(projects),
		)
	}
	if len(projects) == 0 {
		return utils.Diagf(
			"Get: Project does not exist",
			"The Query Returns no Project matching filter %v",
			params,
		)
	}

	project := projects[0]
	d = setProjectResourceData(d, project)
	return diags
}
