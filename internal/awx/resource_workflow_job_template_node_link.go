package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceWorkflowJobTemplateNodeLink() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows you to associate and disassociate a workflow node to another one.",
		CreateContext: resourceWorkflowJobTemplateNodeLinkCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeLinkRead,
		//UpdateContext: resourceWorkflowJobTemplateNodeLinkUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeLinkDelete,
		Schema: map[string]*schema.Schema{
			"origin_node_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the awx_workflow_job_template_node from which the link is starting.",
				ForceNew:    true,
			},
			"next_node_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the awx_workflow_job_template_node to which the link is arriving.",
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"success", "failure", "always"}, false),
				Description:  "The type of the link between 'origin_node_id' and 'next_node_id'. One of \"success\", \"failure\", \"always\"",
				ForceNew:     true,
			},
		},
	}
}

func resourceWorkflowJobTemplateNodeLinkCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	originNodeID := d.Get("origin_node_id").(int)
	nextNodeID := d.Get("next_node_id").(int)
	linkType := d.Get("type").(string)

	client := m.(*awx.AWX)
	if err := client.WorkflowJobTemplateNodeService.AssociateNode(originNodeID, nextNodeID, linkType); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(originNodeID))
	return nil
}

func resourceWorkflowJobTemplateNodeLinkRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowJobTemplateNodeLinkDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	originNodeID := d.Get("origin_node_id").(int)
	nextNodeID := d.Get("next_node_id").(int)
	linkType := d.Get("type").(string)

	client := m.(*awx.AWX)

	if err := client.WorkflowJobTemplateNodeService.DisassociateNode(originNodeID, nextNodeID, linkType); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
