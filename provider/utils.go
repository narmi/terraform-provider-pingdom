package provider

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func diffSuppresorJSON(k, old, new string, d *schema.ResourceData) bool {
	var oldObject, newObject interface{}

	if err := json.Unmarshal([]byte(old), &oldObject); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newObject); err != nil {
		return false
	}

	return reflect.DeepEqual(oldObject, newObject)
}
