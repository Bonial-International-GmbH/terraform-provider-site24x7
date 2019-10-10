package site24x7

import (
	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"oauth2_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_CLIENT_ID", nil),
				Description: "OAuth2 Client ID",
			},
			"oauth2_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_CLIENT_SECRET", nil),
				Description: "OAuth2 Client Secret",
			},
			"oauth2_refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_REFRESH_TOKEN", nil),
				Description: "OAuth2 Refresh Token",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"site24x7_dns_monitor":     resourceSite24x7DnsMonitor(),
			"site24x7_website_monitor": resourceSite24x7WebsiteMonitor(),
			"site24x7_monitor_group":   resourceSite24x7MonitorGroup(),
			"site24x7_action":          resourceSite24x7Action(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := site24x7.Config{
		ClientID:     d.Get("oauth2_client_id").(string),
		ClientSecret: d.Get("oauth2_client_secret").(string),
		RefreshToken: d.Get("oauth2_refresh_token").(string),
	}

	return site24x7.New(config), nil
}
