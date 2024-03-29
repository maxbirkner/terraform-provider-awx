package awx

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagSettingsTitle = "Settings Configration"

var ldapTeamMapAccessMutex sync.Mutex

func resourceSettingsLDAPTeamMap() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `settings_ldap_team_map` manages the LDAP team map settings in AWX.",
		CreateContext: resourceSettingsLDAPTeamMapCreate,
		ReadContext:   resourceSettingsLDAPTeamMapRead,
		DeleteContext: resourceSettingsLDAPTeamMapDelete,
		UpdateContext: resourceSettingsLDAPTeamMapUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this Team",
			},
			"users": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Group DNs to map to this team",
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the team organization",
			},
			"remove": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When True, a user who is not a member of the given groups will be removed from the team",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

type teamMapEntry struct {
	UserDNs      interface{} `json:"users"`
	Organization string      `json:"organization"`
	Remove       bool        `json:"remove"`
}

type teamMap map[string]teamMapEntry

func resourceSettingsLDAPTeamMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	ldapTeamMapAccessMutex.Lock()
	defer ldapTeamMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	res, err := client.SettingService.GetSettingsBySlug("ldap", make(map[string]string))
	if err != nil {
		return utils.DiagCreate(diagSettingsTitle, err)
	}

	/*return buildDiagnosticsMessage(
		"returning as desired",
		"Data: %v", res,
	)*/
	tMaps := make(teamMap)
	if err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tMaps); err != nil {
		return utils.Diagf(
			"Create: failed to parse AUTH_LDAP_TEAM_MAP setting",
			"Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s with input %s", err.Error(), (*res)["AUTH_LDAP_TEAM_MAP"],
		)
	}

	name := d.Get("name").(string)
	if _, ok := tMaps[name]; ok {
		return utils.Diagf(
			"Create: team map already exists",
			"Map for ldap to team map %v already exists", d.Id(),
		)
	}

	newTMap := teamMapEntry{
		UserDNs:      d.Get("users").([]interface{}),
		Organization: d.Get("organization").(string),
		Remove:       d.Get("remove").(bool),
	}

	tMaps[name] = newTMap

	payload := map[string]interface{}{
		"AUTH_LDAP_TEAM_MAP": tMaps,
	}

	if _, err = client.SettingService.UpdateSettings("ldap", payload, make(map[string]string)); err != nil {
		return utils.Diagf(
			"Create: team map not created",
			"failed to save team map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsLDAPTeamMapRead(ctx, d, m)
}

func resourceSettingsLDAPTeamMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	ldapTeamMapAccessMutex.Lock()
	defer ldapTeamMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	res, err := client.SettingService.GetSettingsBySlug("ldap", make(map[string]string))
	if err != nil {
		return utils.Diagf(
			"Update: Unable to fetch settings",
			"Unable to load settings with slug ldap: got %s", err.Error(),
		)
	}

	tMaps := make(teamMap)
	err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tMaps)
	if err != nil {
		return utils.Diagf(
			"Update: failed to parse AUTH_LDAP_TEAM_MAP setting",
			"Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s", err.Error(),
		)
	}

	id := d.Id()
	name := d.Get("name").(string)
	organization := d.Get("organization").(string)
	users := d.Get("users").([]interface{})
	remove := d.Get("remove").(bool)

	if name != id {
		tMaps[name] = tMaps[id]
		delete(tMaps, id)
	}

	utMap := tMaps[name]
	utMap.UserDNs = users
	utMap.Organization = organization
	utMap.Remove = remove
	tMaps[name] = utMap

	payload := map[string]interface{}{
		"AUTH_LDAP_TEAM_MAP": tMaps,
	}

	if _, err = client.SettingService.UpdateSettings("ldap", payload, make(map[string]string)); err != nil {
		return utils.Diagf(
			"Update: team map not created",
			"failed to save team map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsLDAPTeamMapRead(ctx, d, m)
}

func resourceSettingsLDAPTeamMapRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	res, err := client.SettingService.GetSettingsBySlug("ldap", make(map[string]string))
	if err != nil {
		return utils.Diagf(
			"Unable to fetch settings",
			"Unable to load settings with slug ldap: got %s",
			err.Error(),
		)
	}
	tMaps := make(teamMap)
	err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tMaps)
	if err != nil {
		return utils.Diagf(
			"Unable to parse AUTH_LDAP_TEAM_MAP",
			"Unable to parse AUTH_LDAP_TEAM_MAP, got: %s", err.Error(),
		)
	}
	mapdef, ok := tMaps[d.Id()]
	if !ok {
		return utils.DiagFetch("team map", d.Id(), nil)
	}

	/*return buildDiagnosticsMessage(
		"returning as desired",
		"Data: %v %T", mapdef.UserDNs, mapdef.UserDNs,
	)*/

	var users []string
	switch tt := mapdef.UserDNs.(type) {
	case string:
		users = []string{tt}
	case []string:
		users = tt
	case []interface{}:
		for _, v := range tt {
			if dn, ok := v.(string); ok {
				users = append(users, dn)
			}
		}
	}

	if err := d.Set("name", d.Id()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("users", users); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization", mapdef.Organization); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("remove", mapdef.Remove); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSettingsLDAPTeamMapDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ldapTeamMapAccessMutex.Lock()
	defer ldapTeamMapAccessMutex.Unlock()

	client := m.(*awx.AWX)

	res, err := client.SettingService.GetSettingsBySlug("ldap", make(map[string]string))
	if err != nil {
		return utils.Diagf(
			"Delete: Unable to fetch settings",
			"Unable to load settings with slug ldap: got %s", err,
		)
	}

	tmaps := make(teamMap)
	err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tmaps)
	if err != nil {
		return utils.Diagf(
			"Delete: failed to parse AUTH_LDAP_TEAM_MAP setting",
			"Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s", err,
		)
	}

	id := d.Id()
	delete(tmaps, id)

	payload := map[string]interface{}{
		"AUTH_LDAP_TEAM_MAP": tmaps,
	}

	if _, err = client.SettingService.UpdateSettings("ldap", payload, make(map[string]string)); err != nil {
		return utils.DiagDelete("team map", id, err)
	}
	d.SetId("")
	return nil
}
