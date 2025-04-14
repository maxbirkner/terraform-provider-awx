package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_team` manages teams within an organization.",
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		DeleteContext: resourceTeamDelete,
		UpdateContext: resourceTeamUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this Team",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this Team.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric ID of the Team organization",
			},
			"role_entitlement": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of role IDs of the role entitlements",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	orgID := d.Get("organization_id").(int)
	teamName := d.Get("name").(string)
	_, res, err := client.TeamService.ListTeams(map[string]string{
		"name":         teamName,
		"organization": strconv.Itoa(orgID),
	})
	if err != nil {
		return utils.Diagf("Create: Fail to find Team", "Fail to find Team %s Organization ID %v, %s", teamName, orgID, err)
	}
	if len(res.Results) >= 1 {
		return utils.Diagf("Create: Already exist", "Team with name %s  already exists in the Organization ID %v", teamName, orgID)
	}

	result, err := client.TeamService.CreateTeam(map[string]interface{}{
		"name":         teamName,
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(int),
	}, map[string]string{})
	if err != nil {
		return utils.Diagf("Create: Team not created", "Team with name %s  in the Organization ID %v not created, %s", teamName, orgID, err)
	}

	d.SetId(strconv.Itoa(result.ID))

	if rent, entOk := d.GetOk("role_entitlement"); entOk {
		entset := rent.(*schema.Set).List()
		if err := roleTeamEntitlementUpdate(m, result.ID, entset, false); err != nil {
			return utils.Diagf("Create: team role entitlement not created", "Role entitlement for team %s not created: %s", teamName, err)
		}
	}

	return resourceTeamRead(ctx, d, m)
}

func roleTeamEntitlementUpdate(m interface{}, teamID int, roles []interface{}, remove bool) error {
	client := m.(*awx.AWX)
	for _, v := range roles {
		emap := v.(map[string]interface{})
		payload := map[string]interface{}{
			"id": emap["role_id"],
		}
		if remove {
			payload["disassociate"] = true // presence of key triggers removal
		}

		if _, err := client.TeamService.UpdateTeamRoleEntitlement(teamID, payload, make(map[string]string)); err != nil {
			return err
		}
	}
	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.TeamService

	id, diags := utils.StateIDToInt("Update Team", d)
	if diags.HasError() {
		return diags
	}
	if d.HasChange("role_entitlement") {
		oi, ni := d.GetChange("role_entitlement")
		if oi == nil {
			oi = new(schema.Set)
		}
		if ni == nil {
			ni = new(schema.Set)
		}
		oe := oi.(*schema.Set)
		ne := ni.(*schema.Set)

		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if err := roleTeamEntitlementUpdate(m, id, remove, true); err != nil {
			return utils.DiagUpdate("Team Role Entitlement", id, err)
		}
		if err := roleTeamEntitlementUpdate(m, id, add, false); err != nil {
			return utils.DiagUpdate("Team Role Entitlement", id, err)
		}
	}
	if _, err := awxService.UpdateTeam(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(int),
	}, map[string]string{}); err != nil {
		return utils.DiagUpdate("Team Role Entitlement", id, err)
	}
	d.Partial(false)
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read Team", d)
	if diags.HasError() {
		return diags
	}

	team, err := client.TeamService.GetTeamByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("team", id, err)
	}
	entitlements, _, err := client.TeamService.ListTeamRoleEntitlements(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("team roles", id, err)
	}

	d = setTeamResourceData(d, team, entitlements)
	return diags
}

func resourceTeamDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)

	id, diags := utils.StateIDToInt("Delete Team", d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.TeamService.DeleteTeam(id); err != nil {
		return utils.DiagDelete("Team", id, err)
	}
	d.SetId("")
	return diags
}

func setTeamResourceData(d *schema.ResourceData, r *awx.Team, e []*awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("organization_id", r.Organization); err != nil {
		fmt.Println("Error setting organization_id", err)
	}

	var entlist []interface{}
	for _, v := range e {
		elem := make(map[string]interface{})
		elem["role_id"] = v.ID
		entlist = append(entlist, elem)
	}
	f := schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id": {Type: schema.TypeInt},
		}})

	ent := schema.NewSet(f, entlist)

	if err := d.Set("role_entitlement", ent); err != nil {
		fmt.Println("Error setting role_entitlement", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
