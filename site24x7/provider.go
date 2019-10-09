package site24x7

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"authtoken": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_AUTHTOKEN", nil),
				Description: "site24x7 auth Account (https://www.site24x7.com/help/api/#authentication)",
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
	h := make(http.Header)
	h.Set("Authorization", "Zoho-authtoken "+d.Get("authtoken").(string))
	return &http.Client{
		Transport: &staticHeaderTransport{
			base:   http.DefaultTransport,
			header: h,
		},
	}, nil
}

type staticHeaderTransport struct {
	base   http.RoundTripper
	header http.Header
}

func (t *staticHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.header {
		req.Header[k] = v
	}
	return t.base.RoundTrip(req)
}
