package awx

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goawx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func labelAssociationStateID(parentID int, organizationID int, name string) string {
	return fmt.Sprintf("%d:%d:%s", parentID, organizationID, url.QueryEscape(name))
}

func findAssociatedLabel(labels []*goawx.Label, name string, organizationID int) *goawx.Label {
	for _, label := range labels {
		if label == nil {
			continue
		}

		if label.Name == name && label.Organization == organizationID {
			return label
		}
	}

	return nil
}

func syncLabelAssociationState(d *schema.ResourceData, parentIDField string, label *goawx.Label) diag.Diagnostics {
	if err := d.Set("name", label.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("organization_id", label.Organization); err != nil {
		return diag.FromErr(fmt.Errorf("error setting organization_id: %w", err))
	}

	d.SetId(labelAssociationStateID(d.Get(parentIDField).(int), label.Organization, label.Name))
	return nil
}
