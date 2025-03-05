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

const diagInstanceGroupTitle = "Instance Group"

func resourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_instance_group` manages instance groups within an AWX instance.",
		CreateContext: resourceInstanceGroupCreate,
		ReadContext:   resourceInstanceGroupRead,
		UpdateContext: resourceInstanceGroupUpdate,
		DeleteContext: resourceInstanceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the instance group.",
			},
			"is_container_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the instance group is a container group.",
			},
			"policy_instance_minimum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The minimum number of instances to run in the instance group.",
			},
			"policy_instance_percentage": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The percentage of instances to run in the instance group.",
			},
			"pod_spec_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   utils.Normalize,
				Description: "The pod spec override for the instance group.",
			},
			"credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   utils.Normalize,
				Description: "ID of the credential of type 'OpenShift or Kubernetes API Bearer Token' to use as remote cluster.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.InstanceGroupsService

	result, err := awxService.CreateInstanceGroup(map[string]interface{}{
		"name":                       d.Get("name").(string),
		"policy_instance_minimum":    d.Get("policy_instance_minimum").(int),
		"is_container_group":         d.Get("is_container_group").(bool),
		"policy_instance_percentage": d.Get("policy_instance_percentage").(int),
		"pod_spec_override":          d.Get("pod_spec_override").(string),
		"credential":                 d.Get("credential_id").(string),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagInstanceGroupTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInstanceGroupRead(ctx, d, m)

}

func resourceInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.InstanceGroupsService.UpdateInstanceGroup(id, map[string]interface{}{
		"name":                       d.Get("name").(string),
		"policy_instance_minimum":    d.Get("policy_instance_minimum").(int),
		"is_container_group":         d.Get("is_container_group").(bool),
		"policy_instance_percentage": d.Get("policy_instance_percentage").(int),
		"pod_spec_override":          d.Get("pod_spec_override").(string),
		"credential":                 d.Get("credential_id").(string),
	}, nil); err != nil {
		return utils.DiagUpdate(diagInstanceGroupTitle, id, err)
	}

	return resourceInstanceGroupRead(ctx, d, m)

}

func resourceInstanceGroupDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.InstanceGroupsService.DeleteInstanceGroup(id); err != nil {
		return utils.DiagDelete(diagInstanceGroupTitle, id, err)
	}
	d.SetId("")
	return nil
}

func resourceInstanceGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := client.InstanceGroupsService.GetInstanceGroupByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagInstanceGroupTitle, id, err)
	}
	d = setInstanceGroupResourceData(d, res)
	return diag.Diagnostics{}
}

func setInstanceGroupResourceData(d *schema.ResourceData, r *awx.InstanceGroup) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("is_container_group", r.IsContainerGroup); err != nil {
		fmt.Println("Error setting is_container_group", err)
	}
	if err := d.Set("pod_spec_override", utils.Normalize(r.PodSpecOverride)); err != nil {
		fmt.Println("Error setting pod_spec_override", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
