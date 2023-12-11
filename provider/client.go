package provider

import (
	"context"
	"net/http"

	"github.com/narmi/terraform-provider-pingdom/pingdom"
)

func (c *ProviderConfig) pingdomClient() (*pingdom.ClientWithResponses, error) {
	var client, err = pingdom.NewClient(c.apiURL)
	if err != nil {
		return nil, err
	}
	client.RequestEditors = append(client.RequestEditors, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
		return nil
	})
	return &pingdom.ClientWithResponses{ClientInterface: client}, nil
}
