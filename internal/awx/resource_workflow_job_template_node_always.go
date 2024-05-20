package awx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceWorkflowJobTemplateNodeAlways() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows you to create, read, update, and delete a Workflow Job Template Node Always.",
		CreateContext: resourceWorkflowJobTemplateNodeAlwaysCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,
		Schema:        workflowJobNodeSchema,
	}
}
func resourceWorkflowJobTemplateNodeAlwaysCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeAlwaysService
	return createNodeForWorkflowJob(ctx, awxService, d, m)
}
