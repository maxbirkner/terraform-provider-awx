package awx

import (
	"context"
	"reflect"
	"testing"
)

func Test_resourceHostStateUpgradeV0(t *testing.T) {
	rawState := map[string]interface{}{
		"name":         "test-host",
		"description":  "a host",
		"inventory_id": 1,
		"group_ids":    []interface{}{10, 20, 30},
		"enabled":      true,
		"instance_id":  "",
		"variables":    "",
	}

	got, err := resourceHostStateUpgradeV0(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("resourceHostStateUpgradeV0() error = %v", err)
	}
	if !reflect.DeepEqual(got, rawState) {
		t.Errorf("resourceHostStateUpgradeV0() = %v, want %v", got, rawState)
	}
}
