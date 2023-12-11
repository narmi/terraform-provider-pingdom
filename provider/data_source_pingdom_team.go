package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/narmi/terraform-provider-pingdom/pingdom"
)

func dataSourcePingdomTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePingdomTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePingdomTeamRead(d *schema.ResourceData, meta interface{}) error {
	var client *pingdom.ClientWithResponses
	var err error
	var teamName string
	var teamFound bool
	var teamsResponse *pingdom.GetAlertingTeamsResponse
	var member_ids []int = make([]int, 0)
	var ctx = context.Background()

	client, err = meta.(*ProviderConfig).pingdomClient()
	if err != nil {
		return fmt.Errorf("Error instantiating Pingdom client: %s", err)
	}

	teamName = d.Get("name").(string)

	teamsResponse, err = client.GetAlertingTeamsWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving teams: %s", err)
	}

	for _, team := range *teamsResponse.JSON200.Teams {
		if *team.Name == teamName {
			teamFound = true
			if err = d.Set("name", team.Name); err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
			// extract member IDs (member.ID)
			for _, member := range *team.Members {
				member_ids = append(member_ids, *member.Id)
			}
			if err = d.Set("member_ids", member_ids); err != nil {
				return fmt.Errorf("Error setting member_ids: %s", err)
			}
		}
	}

	if !teamFound {
		return fmt.Errorf("Team not found: %s", teamName)
	}
	return nil
}
