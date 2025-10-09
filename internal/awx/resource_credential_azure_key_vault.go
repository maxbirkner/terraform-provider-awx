package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const azureKeyVaultCredentialTypeName = "Microsoft Azure Key Vault"

func resourceCredentialAzureKeyVault() *schema.Resource {
	return &schema.Resource{
		Description:   "The `awx_credential_azure_key_vault` resource allows you to manage Azure Key Vault credentials in Ansible AWX.",
		CreateContext: resourceCredentialAzureKeyVaultCreate,
		ReadContext:   resourceCredentialAzureKeyVaultRead,
		UpdateContext: resourceCredentialAzureKeyVaultUpdate,
		DeleteContext: resourceCredentialDelete,
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
			"cloud_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "AzureCloud",
				Description: "The Azure cloud environment. Options: AzureCloud, AzureUSGovernment, AzureChinaCloud, AzureGermanCloud.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialAzureKeyVaultCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	azureKVCredType, err := client.CredentialTypeService.GetCredentialTypeByName(azureKeyVaultCredentialTypeName, map[string]string{})
	if err != nil {
		return utils.DiagCreate("Azure Key Vault Credential Type", err)
	}

	inputs := map[string]interface{}{
		"url":    d.Get("url").(string),
		"client": d.Get("client").(string),
		"secret": d.Get("secret").(string),
		"tenant": d.Get("tenant").(string),
	}

	if cloudName, ok := d.GetOk("cloud_name"); ok {
		inputs["cloud_name"] = cloudName.(string)
	}

	payload := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": azureKVCredType.ID,
		"inputs":          inputs,
	}

	cred, err := client.CredentialsService.CreateCredentials(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate("Azure Key Vault Credential", err)
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialAzureKeyVaultRead(ctx, d, m)

	return diag.Diagnostics{}
}

func resourceCredentialAzureKeyVaultRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagFetch("Azure Key Vault Credential", d.Id(), err)
	}
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch("Azure Key Vault Credential", d.Id(), err)
	}

	if err := d.Set("name", cred.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", cred.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", cred.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", cred.Inputs["url"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client", cred.Inputs["client"]); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "secret", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant", cred.Inputs["tenant"]); err != nil {
		return diag.FromErr(err)
	}
	if cloudName, ok := cred.Inputs["cloud_name"]; ok {
		if err := d.Set("cloud_name", cloudName); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func resourceCredentialAzureKeyVaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	keys := []string{
		"name",
		"description",
		"url",
		"client",
		"secret",
		"tenant",
		"cloud_name",
		"organization_id",
	}

	if d.HasChanges(keys...) {
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return utils.DiagUpdate("Azure Key Vault Credential", d.Id(), err)
		}

		client := m.(*awx.AWX)
		azureKVCredType, err := client.CredentialTypeService.GetCredentialTypeByName(azureKeyVaultCredentialTypeName, map[string]string{})
		if err != nil {
			return utils.DiagUpdate("Azure Key Vault Credential Type", d.Id(), err)
		}

		inputs := map[string]interface{}{
			"url":    d.Get("url").(string),
			"client": d.Get("client").(string),
			"secret": d.Get("secret").(string),
			"tenant": d.Get("tenant").(string),
		}

		if cloudName, ok := d.GetOk("cloud_name"); ok {
			inputs["cloud_name"] = cloudName.(string)
		}

		payload := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": azureKVCredType.ID,
			"inputs":          inputs,
		}

		if _, err = client.CredentialsService.UpdateCredentialsByID(id, payload, map[string]string{}); err != nil {
			return utils.DiagUpdate("Azure Key Vault Credential", d.Id(), err)
		}
	}

	return resourceCredentialAzureKeyVaultRead(ctx, d, m)
}
