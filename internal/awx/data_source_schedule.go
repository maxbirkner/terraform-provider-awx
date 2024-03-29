package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagScheduleTitle = "Schedule"

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchedulesRead,
		Description: "Data source for AWX Schedule",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the Schedule",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Schedule",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	schedules, _, err := client.ScheduleService.List(params)
	if err != nil {
		return utils.DiagFetch(diagScheduleTitle, params, err)
	}
	if len(schedules) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(schedules),
		)
	}
	if len(schedules) == 0 {
		return utils.Diagf(
			"Get: Schedule does not exist",
			"The Query Returns no Schedule matching filter %v",
			params,
		)
	}

	schedule := schedules[0]
	d = setScheduleResourceData(d, schedule)
	return diags
}
