package searchstax

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrf110/terraform-provider-searchstax/api"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEARCHSTAX_ACCOUNT_NAME", nil),
				Description: "Account name for this SearchStax account",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SEARCHSTAX_USERNAME", nil),
				Description: "Username used to authenticate to the SearchStax API",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SEARCHSTAX_PASSWORD", nil),
				Description: "Password used to authenticate to the SearchStax API",
			},
			"base_url": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "https://app.searchstax.com/api/rest/v2",
				Optional:    true,
				Description: "Base URL to SearchStax API",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"apikey": resourceApiKey(),
			"deployment": resourceDeployment(),
			"apikey_association": resourceApiKeyAssociation(),
			"user": resourceUser(),
			"alert": resourceAlert(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	return api.New(
		data.Get("account_name").(string),
		data.Get("username").(string),
		data.Get("password").(string),
		data.Get("base_url").(string))
}
