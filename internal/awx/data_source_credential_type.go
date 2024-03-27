package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentialTypeByID() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialTypeByIDRead,
		Description: "Use this data source to get the details of a credential type by ID.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the credential type",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the credential type",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the credential type",
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kind of the credential type",
			},
			"inputs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The inputs of the credential type",
			},
			"injectors": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The injectors of the credential type",
			},
		},
	}
}

func dataSourceCredentialTypeByIDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id := d.Get("id").(int)
	credType, err := client.CredentialTypeService.GetCredentialTypeByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to fetch credential type with ID: %d. Error: %s", id, err.Error()),
		})
	}

	d.Set("name", credType.Name)
	d.Set("description", credType.Description)
	d.Set("kind", credType.Kind)
	d.Set("inputs", credType.Inputs)
	d.Set("injectors", credType.Injectors)
	d.SetId(strconv.Itoa(id))

	return diags
}
