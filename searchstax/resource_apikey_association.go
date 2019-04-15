package searchstax

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrf110/terraform-provider-searchstax/api"
	"github.com/nrf110/terraform-provider-searchstax/api/apikeys"
)

func resourceApiKeyAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyAssociationCreate,
		Delete: resourceApiKeyAssociationDelete,
		Read: resourceApiKeyAssociationRead,
		Exists: resourceApiKeyAssociationExists,

		Schema: map[string]*schema.Schema{
			"apikey": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deployment": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApiKeyAssociationCreate(data *schema.ResourceData, meta interface{}) error  {
	client := meta.(*api.API)
	apiKey := data.Get("apikey").(string)
	deployment := data.Get("deployment").(string)

	_, err := client.ApiKeys.Associate(apikeys.ApiKeyAssociationRequest{
		ApiKey: apiKey,
		Deployment: deployment,
	})

	data.SetId(fmt.Sprintf("%s-%s", apiKey, deployment))

	return err
}

func resourceApiKeyAssociationDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.API)
	_, err := client.ApiKeys.Disassociate(apikeys.ApiKeyAssociationRequest{
		ApiKey: data.Get("apikey").(string),
		Deployment: data.Get("deployment").(string),
	})

	if err != nil {
		return err
	}

	return resourceApiKeyAssociationRead(data, meta)
}

func resourceApiKeyAssociationRead(data *schema.ResourceData, meta interface{}) error {
	if exists, err := resourceApiKeyAssociationExists(data, meta); !exists || err != nil {
		data.SetId("")
		return err
	}

	return nil
}

func resourceApiKeyAssociationExists(data *schema.ResourceData, meta interface{}) (bool, error)  {
	client := meta.(*api.API)

	apiKeys, err := client.ApiKeys.GetForDeployment(data.Get("deployment").(string))
	if err != nil {
		return false, err
	}

	for _, apiKey := range apiKeys {
		if apiKey == data.Get("apikey").(string) {
			return true, nil
		}
	}

	return false, nil
}
