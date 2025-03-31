package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

// JobTemplate and WorkflowJobTemplate share the same SurveySpec, use one resource for both
const workflowPrefix = "workflow_"
const workflowSpecTitlePrefix = "Workflow "
const diagGenericSurveySpecTitle = "%sJob Template Survey Spec"

var diagSurveySpecTitle = ""

//nolint:funlen
func resourceSurveySpec(isWorkflow bool) *schema.Resource {
	prefix := ""
	if isWorkflow {
		prefix = workflowPrefix
		diagSurveySpecTitle = fmt.Sprintf(diagGenericSurveySpecTitle, workflowSpecTitlePrefix)
	} else {
		diagSurveySpecTitle = fmt.Sprintf(diagGenericSurveySpecTitle, "")
	}

	return &schema.Resource{
		Description:   fmt.Sprintf("Resource `awx_%sjob_template_survey_spec` manages job templates surveys within AWX.", prefix),
		CreateContext: resourceSurveySpecCreate(isWorkflow),
		ReadContext:   resourceSurveySpecRead(isWorkflow),
		UpdateContext: resourceSurveySpecCreate(isWorkflow),
		DeleteContext: resourceSurveySpecDelete(isWorkflow),

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the Job Template to which create the survey spec.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Job Template survey spec.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the job template survey spec.",
			},
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the job template survey. One block per question in the survey.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"text", "password", "integer", "float", "multiplechoice", "multiselect"}, false),
							Description:  "The type of the question. One of \"text\", \"password\", \"integer\", \"float\", \"multiplechoice\" or \"multiselect\"",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Boolean to decide if this variable is required for the survey.",
							Default:     true,
						},
						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default answer value for the survey question.",
							Default:     "",
						},
						"variable": {
							Type:        schema.TypeString,
							Description: "Name of the ansible variable that will be provisioned with the response.",
							Required:    true,
						},
						"question_name": {
							Type:        schema.TypeString,
							Description: "Name of the question that will be asked to the user.",
							Required:    true,
						},
						"question_description": {
							Type:        schema.TypeString,
							Description: "Description of the question that will be asked to the user.",
							Optional:    true,
						},
						"min": {
							Type:        schema.TypeInt,
							Description: "Minimum length of the answer for type 'text' or 'password' and minimum value of the response for type 'integer' or 'float'. Defaults to 0",
							Optional:    true,
							Default:     0,
						},
						"max": {
							Type:        schema.TypeInt,
							Description: "Maximum length of the answer for type 'text' or 'password' and maximum value of the response for type 'integer' or 'float'. Defaults to 1024",
							Optional:    true,
							Default:     1024,
						},
						"choices": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of choices when type is 'multiplechoice' or 'multiselect'.",
							Optional:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSurveySpecRead(isWorkflow bool) func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		jobTemplateID := d.Get("job_template_id").(int)

		surveySpec, err := client.SurveySpecService.GetSurveySpec(isWorkflow, jobTemplateID, map[string]string{})
		if err != nil {
			return utils.DiagNotFound(diagSurveySpecTitle, jobTemplateID, err)
		}

		if err := d.Set("name", surveySpec.Name); err != nil {
			return nil
		}
		if err := d.Set("description", surveySpec.Description); err != nil {
			return nil
		}
		if err := d.Set("spec", surveySpec.Spec); err != nil {
			return nil
		}
		d.SetId(strconv.Itoa(jobTemplateID))
		return nil
	}
}

func resourceSurveySpecCreate(isWorkflow bool) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		jobTemplateID := d.Get("job_template_id").(int)

		_, err := client.SurveySpecService.CreateSurveySpec(isWorkflow, jobTemplateID, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"spec":        d.Get("spec").([]interface{}),
		})
		if err != nil {
			return utils.DiagCreate(diagSurveySpecTitle, err)
		}

		resourceSurveySpecRead(isWorkflow)
		d.SetId(strconv.Itoa(jobTemplateID))
		return nil
	}
}

func resourceSurveySpecDelete(isWorkflow bool) func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*awx.AWX)
		jobTemplateID := d.Get("job_template_id").(int)

		if err := client.SurveySpecService.DeleteSurveySpec(isWorkflow, jobTemplateID); err != nil {
			return utils.DiagDelete(diagSurveySpecTitle, jobTemplateID, err)
		}
		d.SetId("")
		return nil
	}
}
