package awx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceWorkflowJobTemplateNodeSuccess() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows you to create, update, and delete a Workflow Job Template Node Success.",
		CreateContext: resourceWorkflowJobTemplateNodeSuccessCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,
		Schema:        workflowJobNodeSchema,
	}
}

func resourceWorkflowJobTemplateNodeSuccessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeSuccessService
	return createNodeForWorkflowJob(awxService, ctx, d, m)
}
