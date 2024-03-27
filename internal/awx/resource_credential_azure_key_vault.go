package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceCredentialAzureKeyVault() *schema.Resource {
	return &schema.Resource{
		Description:   "The `awx_credential_azure_key_vault` resource allows you to manage Azure Key Vault credentials in Ansible AWX.",
		CreateContext: resourceCredentialAzureKeyVaultCreate,
		ReadContext:   resourceCredentialAzureKeyVaultRead,
		UpdateContext: resourceCredentialAzureKeyVaultUpdate,
		DeleteContext: CredentialsServiceDeleteByID,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the credential.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The organization ID that the credential belongs to.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the Azure Key Vault.",
			},
			"client": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The client ID of the Azure Key Vault.",
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The secret of the Azure Key Vault.",
			},
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tenant ID of the Azure Key Vault.",
			},
		},
	}
}

func resourceCredentialAzureKeyVaultCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 19, // Azure Key Vault
		"inputs": map[string]interface{}{
			"url":    d.Get("url").(string),
			"client": d.Get("client").(string),
			"secret": d.Get("secret").(string),
			"tenant": d.Get("tenant").(string),
		},
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialsService.CreateCredentials(newCredential, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to create new credentials: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialAzureKeyVaultRead(ctx, d, m)

	return diags
}

func resourceCredentialAzureKeyVaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", cred.Name)
	d.Set("description", cred.Description)
	d.Set("organization_id", cred.OrganizationID)
	d.Set("url", cred.Inputs["url"])
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("tenant", cred.Inputs["tenant"])

	return diags
}

func resourceCredentialAzureKeyVaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"url",
		"client",
		//"secret",
		"tenant",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 19, // Azure Key Vault
			"inputs": map[string]interface{}{
				"url":    d.Get("url").(string),
				"client": d.Get("client").(string),
				"secret": d.Get("secret").(string),
				"tenant": d.Get("tenant").(string),
			},
		}

		client := m.(*awx.AWX)
		_, err = client.CredentialsService.UpdateCredentialsByID(id, updatedCredential, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credentials",
				Detail:   fmt.Sprintf("Unable to update existing credentials with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialAzureKeyVaultRead(ctx, d, m)
}
