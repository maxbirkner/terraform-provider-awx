package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

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

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if teamName, okName := d.GetOk("name"); okName {
		params["name"] = teamName.(string)
	}

	if teamID, okTeamID := d.GetOk("id"); okTeamID {
		params["id"] = strconv.Itoa(teamID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}
	Teams, _, err := client.TeamService.ListTeams(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Team",
			"Fail to find the team got: %s",
			err.Error(),
		)
	}

	if len(Teams) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one team, %d",
			len(Teams),
		)
	}

	if len(Teams) == 0 {
		return buildDiagnosticsMessage(
			"Get: Team does not exist",
			"The Query Returns no Team matching filter %v",
			params,
		)
	}

	Team := Teams[0]
	Entitlements, _, err := client.TeamService.ListTeamRoleEntitlements(Team.ID, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch team role entitlements",
			"Fail to retrieve team role entitlements got: %s",
			err.Error(),
		)
	}

	d = setTeamResourceData(d, Team, Entitlements)
	return diags
}
