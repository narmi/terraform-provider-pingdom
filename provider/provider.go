package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	apiToken string
	apiURL   string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PINGDOM_API_TOKEN", nil),
				Description: "Pingdom API Token",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PINGDOM_API_URL", "https://api.pingdom.com/api/3.1/"),
				Description: "Pingdom API URL",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"pingdom_check": resourcePingdomCheck(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pingdom_team": dataSourcePingdomTeam(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var api_token string = d.Get("api_token").(string)
	var api_url string = d.Get("api_url").(string)

	return &ProviderConfig{
		apiToken: api_token,
		apiURL:   api_url,
	}, nil
}
