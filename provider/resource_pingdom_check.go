package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/narmi/terraform-provider-pingdom/pingdom"
)

var pingdomCheckSchemaVersion int = 0

// pingdom.CreateCheckType
var validCheckTypes = []string{
	"dns",
	"http",
	"httpcustom",
	"imap",
	"ping",
	"pop3",
	"smtp",
	"tcp",
	"udp",
}

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: pingdomCheckSchemaVersion,

		Create: resourcePingdomCheckCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// from "type CreateCheck" in pingdom.gen.go
		//
		Schema: map[string]*schema.Schema{
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateCheckType,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"parameters": {
				Type:             schema.TypeString,
				Required:         false,
				ForceNew:         false,
				DiffSuppressFunc: diffSuppresorJSON,
				ValidateFunc:     validation.StringIsJSON,
			},

			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: false,
			},
		},
	}
}

func validateCheckType(val interface{}, path cty.Path) diag.Diagnostics {
	var checkType string = val.(string)

	for _, validCheckType := range validCheckTypes {
		if checkType == validCheckType {
			return nil
		}
	}
	return diag.Errorf("Invalid check type: %s. Allowed values are: %v", checkType, validCheckTypes)
}

func resourcePingdomCheckCreate(d *schema.ResourceData, meta interface{}) error {
	var client *pingdom.ClientWithResponses
	var err error
	var ctx = context.Background()
	var checkType, checkName, checkHost, checkParameters string
	var response *pingdom.PostChecksResponse
	var responseID int

	client, err = meta.(*ProviderConfig).pingdomClient()
	if err != nil {
		return fmt.Errorf("Error instantiating Pingdom client: %s", err)
	}

	checkType = d.Get("type").(string)
	checkName = d.Get("name").(string)
	checkHost = d.Get("host").(string)
	checkParameters = d.Get("parameters").(string)

	var checkCreator pingdom.CreateCheck
	checkCreator.UnmarshalJSON([]byte(checkParameters))

	checkCreator.Name = checkName
	checkCreator.Host = checkHost
	checkCreator.Type = pingdom.CreateCheckType(checkType)

	response, err = client.PostChecksWithResponse(ctx, checkCreator)
	if err != nil {
		return fmt.Errorf("Error creating Pingdom check: %s", err)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("Server returned failure while creating Pingdom check: %s", response.Status())
	}

	responseID = *response.JSON200.Check.Id
	if err = d.Set("id", responseID); err != nil {
		return fmt.Errorf("Error setting chcek id: %s", err)
	}

	return nil
}
