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
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"id": {
				Type:     schema.TypeInt,
				Required: false,
				Computed: true,
				ForceNew: false,
				ExactlyOneOf: []string{
					"name",
					"id",
				},
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePingdomTeamRead(d *schema.ResourceData, meta interface{}) error {
	var client *pingdom.ClientWithResponses
	var err error
	var got bool
	var val interface{}
	var teamName string
	var teamID int
	var teamFound bool
	var teamIDResponse *pingdom.GetAlertingTeamsTeamidResponse
	var teamsResponse *pingdom.GetAlertingTeamsResponse
	var member_ids []int = make([]int, 0)
	var ctx = context.Background()

	client, err = meta.(*ProviderConfig).pingdomClient()
	if err != nil {
		return fmt.Errorf("Error instantiating Pingdom client: %s", err)
	}

	val, got = d.GetOk("name")
	if got {
		teamName = val.(string)

		teamsResponse, err = client.GetAlertingTeamsWithResponse(ctx)
		if err != nil {
			return fmt.Errorf("Error retrieving teams: %s", err)
		}

		if teamsResponse.StatusCode() != 200 {
			return fmt.Errorf("Server returned failure while reading Pingdom teams: %s", teamsResponse.Status())
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
	} else { // id has to exist
		val, got = d.GetOk("id")
		if !got {
			return fmt.Errorf("Team name or id must be specified")
		}
		teamID = val.(int)

		teamIDResponse, err = client.GetAlertingTeamsTeamidWithResponse(ctx, teamID)
		if err != nil {
			return fmt.Errorf("Error retrieving team: %s", err)
		}
		if teamIDResponse.StatusCode() != 200 {
			return fmt.Errorf("Server returned failure while reading Pingdom team: %s", teamIDResponse.Status())
		}
		if err = d.Set("name", teamIDResponse.JSON200.Team.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
		// extract member IDs (member.ID)
		for _, member := range *teamIDResponse.JSON200.Team.Members {
			member_ids = append(member_ids, *member.Id)
		}
		if err = d.Set("member_ids", member_ids); err != nil {
			return fmt.Errorf("Error setting member_ids: %s", err)
		}
	}
	return nil
}
