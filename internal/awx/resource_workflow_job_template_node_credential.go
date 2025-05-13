package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceWorkflowJobTemplateNodeCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Associates a credential to a workflow job template node",
		CreateContext: resourceWorkflowJobTemplateNodeCredentialCreate,
		DeleteContext: resourceWorkflowJobTemplateNodeCredentialDelete,
		ReadContext:   resourceWorkflowJobTemplateNodeCredentialRead,
		Schema: map[string]*schema.Schema{
			"workflow_job_template_node_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the job template to associate the credential with",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the credential to associate with the job template",
			},
		},
	}
}

func resourceWorkflowJobTemplateNodeCredentialCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	workflowJobTemplateNodeID := d.Get("workflow_job_template_node_id").(int)
	res, err := client.WorkflowJobTemplateNodeService.GetWorkflowJobTemplateNodeByID(workflowJobTemplateNodeID, make(map[string]string))

	if err != nil {
		return utils.DiagNotFound("Workflow Job Template Node", workflowJobTemplateNodeID, err)
	}

	if err = client.WorkflowJobTemplateNodeService.AssociateCredential(res.ID, d.Get("credential_id").(int)); err != nil {
		return utils.DiagCreate("JobTemplate AssociateCredentials", err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return nil
}

func resourceWorkflowJobTemplateNodeCredentialRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowJobTemplateNodeCredentialDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	workflowJobTemplateNodeID := d.Get("workflow_job_template_node_id").(int)
	res, err := client.WorkflowJobTemplateNodeService.GetWorkflowJobTemplateNodeByID(workflowJobTemplateNodeID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("Workflow Job Template Node", workflowJobTemplateNodeID, err)
	}

	if err = client.WorkflowJobTemplateNodeService.DisassociateCredential(res.ID, d.Get("credential_id").(int)); err != nil {
		return utils.DiagDelete("JobTemplate DisassociateCredentials", workflowJobTemplateNodeID, err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return nil
}
