package searchstax

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrf110/terraform-provider-searchstax/api"
	"github.com/nrf110/terraform-provider-searchstax/api/apikeys"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyCreate,
		Read: resourceApiKeyRead,
		Delete: resourceApiKeyDelete,

		Schema: map[string]*schema.Schema{
			"scope": {
				Type: schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: "The scope in which this apikey will be effective",
				Default: []string{ "deployment.dedicateddeployment" },
			},
		},
	}
}

func resourceApiKeyCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.API)

	apiKey, err := client.ApiKeys.Create(apikeys.CreateApiKeyRequest{
		Scope: data.Get("scope").([]string),
	})

	if err != nil {
		return err
	}

	data.SetId(apiKey.ApiKey)

	return nil
}

func resourceApiKeyDelete(data *schema.ResourceData, meta interface{}) error {
	var err error
	var deployments []string

	client := meta.(*api.API)

	deployments, err = client.ApiKeys.GetForDeployment(data.Id())
	if err != nil {
		return err
	}

	for _, deployment := range deployments {
		_, err = client.ApiKeys.Disassociate(apikeys.ApiKeyAssociationRequest{
			ApiKey: data.Id(),
			Deployment: deployment,
		})

		if err != nil {
			return err
		}
	}

	_, err = client.ApiKeys.Revoke(apikeys.RevokeApiKeyRequest{
		ApiKey: data.Get("apikey").(string),
	})

	if err != nil {
		data.SetId("")
	}

	return err
}

func resourceApiKeyRead(data *schema.ResourceData, meta interface{}) error {
	return nil
}