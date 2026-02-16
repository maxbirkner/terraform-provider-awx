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

const diagCredentialTitle = "Credential"

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "The `awx_credential` resource allows you to create and manage credentials in Ansible Tower.",
		CreateContext: resourceCredentialCreate,
		ReadContext:   resourceCredentialRead,
		UpdateContext: resourceCredentialUpdate,
		DeleteContext: resourceCredentialDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the credential",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The organization ID that the credential belongs to",
			},
			"credential_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specify the type of credential you want to create. Refer to the Ansible Tower documentation for details on each type",
			},
			"inputs": {
				Type:        schema.TypeMap,
				Required:    true,
				Sensitive:   true,
				Description: "The inputs to be created with the credential.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inputsMap := d.Get("inputs")

	payload := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": d.Get("credential_type_id").(int),
		"inputs":          inputsMap,
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialsService.CreateCredentials(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagCredentialTitle, err)
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialRead(ctx, d, m)
	return diag.Diagnostics{}
}

func resourceCredentialRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagFetch(diagCredentialTitle, d.Id(), err)
	}
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch(diagCredentialTitle, d.Id(), err)
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

	// Fetch the credential type to identify which input fields are secret.
	// AWX returns "$encrypted$" for secret fields, which would cause
	// perpetual drift. For those fields, preserve the current state value.
	secretFields := getSecretFields(client, cred.CredentialTypeID)
	inputs := cred.Inputs
	if len(secretFields) > 0 {
		currentInputs, ok := d.GetOk("inputs")
		var stateInputs map[string]interface{}
		if ok {
			stateInputs, _ = currentInputs.(map[string]interface{})
		}
		inputs = sanitizeEncryptedInputs(inputs, stateInputs, secretFields)
	}

	if err := d.Set("inputs", inputs); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// sanitizeEncryptedInputs replaces "$encrypted$" values in api-returned
// inputs with the corresponding values from the current Terraform state,
// but only for fields that are known to be secret. Non-secret fields and
// secret fields that do not contain "$encrypted$" are left unchanged.
func sanitizeEncryptedInputs(apiInputs, stateInputs map[string]interface{}, secretFields map[string]struct{}) map[string]interface{} {
	result := make(map[string]interface{}, len(apiInputs))
	for k, v := range apiInputs {
		result[k] = v
	}

	for fieldName := range secretFields {
		val, exists := result[fieldName]
		if !exists {
			continue
		}
		strVal, isStr := val.(string)
		if isStr && strVal == "$encrypted$" && stateInputs != nil {
			if stateVal, hasState := stateInputs[fieldName]; hasState {
				result[fieldName] = stateVal
			}
		}
	}

	return result
}

// getSecretFields fetches the credential type by ID and returns a set of
// field IDs that have "secret": true in the credential type's input schema.
// The AWX credential type inputs schema looks like:
//
//	{"fields": [{"id": "username", ...}, {"id": "password", "secret": true, ...}]}
//
// On any error (network, parsing), it returns an empty map so the read
// operation degrades gracefully rather than failing.
func getSecretFields(client *awx.AWX, credentialTypeID int) map[string]struct{} {
	credType, err := client.CredentialTypeService.GetCredentialTypeByID(credentialTypeID, map[string]string{})
	if err != nil {
		fmt.Printf("[WARN] Unable to fetch credential type %d to determine secret fields: %s\n", credentialTypeID, err)
		return make(map[string]struct{})
	}

	return parseSecretFieldsFromInputs(credType.Inputs)
}

// parseSecretFieldsFromInputs parses the credential type's Inputs schema
// (an interface{} that should be map[string]interface{} with a "fields" key)
// and returns the set of field IDs that have "secret": true.
func parseSecretFieldsFromInputs(inputs interface{}) map[string]struct{} {
	secretFields := make(map[string]struct{})

	inputsMap, ok := inputs.(map[string]interface{})
	if !ok {
		return secretFields
	}

	fieldsRaw, ok := inputsMap["fields"]
	if !ok {
		return secretFields
	}

	fields, ok := fieldsRaw.([]interface{})
	if !ok {
		return secretFields
	}

	for _, fieldRaw := range fields {
		field, ok := fieldRaw.(map[string]interface{})
		if !ok {
			continue
		}
		secret, hasSecret := field["secret"]
		if !hasSecret {
			continue
		}
		if secretBool, ok := secret.(bool); ok && secretBool {
			if fieldID, ok := field["id"].(string); ok {
				secretFields[fieldID] = struct{}{}
			}
		}
	}

	return secretFields
}

func resourceCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	keys := []string{
		"name",
		"description",
		"organization_id",
		"inputs",
	}

	if d.HasChanges(keys...) {
		var err error

		inputsMap := d.Get("inputs")

		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return utils.DiagUpdate(diagCredentialTitle, d.Id(), err)
		}
		update := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": d.Get("credential_type_id"),
			"inputs":          inputsMap,
		}

		client := m.(*awx.AWX)
		if _, err = client.CredentialsService.UpdateCredentialsByID(id, update, map[string]string{}); err != nil {
			return utils.DiagUpdate(diagCredentialTitle, d.Id(), err)
		}
	}

	return resourceCredentialRead(ctx, d, m)
}

func resourceCredentialDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagDelete(diagCredentialTitle, d.Id(), err)
	}
	client := m.(*awx.AWX)
	if err := client.CredentialsService.DeleteCredentialsByID(id, map[string]string{}); err != nil {
		return utils.DiagDelete(diagCredentialTitle, d.Id(), err)
	}

	return diag.Diagnostics{}
}
