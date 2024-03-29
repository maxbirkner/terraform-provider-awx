package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagTeamTitle = "Team"

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamsRead,
		Description: "Use this data source to get the details of a team in AWX.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the team.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the team.",
			},
		},
	}
}

func dataSourceTeamsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if teamName, okName := d.GetOk("name"); okName {
		params["name"] = teamName.(string)
	}

	if teamID, okTeamID := d.GetOk("id"); okTeamID {
		params["id"] = strconv.Itoa(teamID.(int))
	}

	teams, _, err := client.TeamService.ListTeams(params)
	if err != nil {
		return utils.DiagFetch(diagTeamTitle, params, err)
	}

	if len(teams) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one team, %d",
			len(teams),
		)
	}

	if len(teams) == 0 {
		return utils.Diagf(
			"Get: Team does not exist",
			"The Query Returns no Team matching filter %v",
			params,
		)
	}

	entitlements, _, err := client.TeamService.ListTeamRoleEntitlements(teams[0].ID, make(map[string]string))
	if err != nil {
		return utils.DiagFetch(diagTeamTitle, teams[0].ID, err)
	}

	d = setTeamResourceData(d, teams[0], entitlements)
	return diags
}
