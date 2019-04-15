package searchstax

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrf110/terraform-provider-searchstax/api"
	"github.com/nrf110/terraform-provider-searchstax/api/deployments"
	"github.com/nrf110/terraform-provider-searchstax/util"
)

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Update: resourceDeploymentUpdate,
		Delete: resourceDeploymentDelete,
		Read: resourceDeploymentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The name of the new deployment",
			},
			"application": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The application for which this deployment is being created",
				Default: "Solr",
			},
			"tier": {
				Type: schema.TypeString,
				Computed: true,
			},
			"http_endpoint": {
				Type: schema.TypeString,
				Computed: true,
			},
			"application_version": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The version of application being deployed",
			},
			"termination_lock": {
				Type: schema.TypeBool,
				Required: true,
				ForceNew: true,
				Description: "Whether or not the deployment should be shielded from API deletion",
			},
			"plan_type": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Available plan type could be Shared Deployment or Dedicated Deployment. Currently SearchStax only supports Dedicated Deployment (the default)",
				Default: "DedicatedDeployment",
			},
			"plan": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The plan of the selected plan type. Detailed information regarding plans available can be fetched using List Plans method.",
			},
			"region_id": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The region ID of the region where deployment has to be created. Detailed information regarding regions available can be fetched using List Plans method.",
			},
			"cloud_provider_id": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The cloud provider to be used to create this deployment. Detailed information regarding cloud providers available can be fetched using List Plans method.",
			},
			"num_additional_app_nodes": {
				Type: schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Description: "The number of additional nodes to be added to the deployment in case you want to add nodes to a dedicated cluster's default 3 nodes.",
			},
			"deployment_type": {
				Type: schema.TypeString,
				Computed: true,
			},
			"num_nodes_default": {
				Type: schema.TypeInt,
				Computed: true,
			},
			"num_additional_zookeeper_nodes": {
				Type: schema.TypeInt,
				Computed: true,
			},
			"servers": {
				Type: schema.TypeList,
				Computed: true,
			},
			"zookeeper_ensemble": {
				Type: schema.TypeList,
				Computed: true,
			},
			"specifications": {
				Type: schema.TypeMap,
				Computed: true,
			},
			"tags": {
				Type: schema.TypeList,
				Optional: true,
				Default: new([]string),
			},
		},
	}
}

func resourceDeploymentCreate(data *schema.ResourceData, meta interface{}) error  {
	client := meta.(*api.API)

	details, err := client.Deployments.Create(deployments.CreateDeploymentRequest{
		Name: data.Get("name").(string),
		Application: data.Get("application").(string),
		ApplicationVersion: data.Get("application_version").(string),
		TerminationLock: data.Get("termination_lock").(bool),
		PlanType: data.Get("plan_type").(string),
		Plan: data.Get("plan").(string),
		RegionID: data.Get("region_id").(string),
		CloudProviderID: data.Get("cloud_provider_id").(string),
		NumAdditionalAppNodes: data.Get("num_additional_app_nodes").(int),
	})

	if err != nil {
		return err
	}

	data.SetId(details.UID)

	tags := data.Get("tags").([]string)
	if len(tags) > 0 {
		_, err = client.Deployments.CreateTags(data.Id(), deployments.TagRequest{
			Tags: tags,
		})

		if err != nil {
			return err
		}
	}

	return resourceDeploymentRead(data, meta)
}

func resourceDeploymentUpdate(data *schema.ResourceData, meta interface{}) error  {
	var err error
	client := meta.(*api.API)

	o, n := data.GetChange("tags")
	prev := o.([]string)
	next := n.([]string)

	removed := util.Diff(prev, next)
	added := util.Diff(next, prev)

	if len(removed) > 0 {
		_, err = client.Deployments.DeleteTags(data.Id(), deployments.TagRequest{
			Tags: removed,
		})

		if err != nil {
			return err
		}
	}

	if len(added) > 0 {
		_, err = client.Deployments.CreateTags(data.Id(), deployments.TagRequest{
			Tags: added,
		})

		if err != nil {
			return err
		}
	}

	return resourceDeploymentRead(data, meta)
}

func resourceDeploymentDelete(data *schema.ResourceData, meta interface{}) error  {
	var err error
	client := meta.(*api.API)

	tags := data.Get("tags").([]string)
	if len(tags) > 0 {
		_, err = client.Deployments.DeleteTags(data.Id(), deployments.TagRequest{
			Tags: tags,
		})

		if err != nil {
			return err
		}
	}

	err = client.Deployments.Delete(data.Id())

	if err != nil {
		return err
	}

	return resourceDeploymentRead(data, meta)
}

func resourceDeploymentRead(data *schema.ResourceData, meta interface{}) error  {
	client := meta.(*api.API)

	details, err := client.Deployments.Get(data.Id())

	// TODO : Distinguish between API errors and missing deployment
	if err != nil {
		data.SetId("")
		return nil
	}

	data.Set("tier", details.Tier)
	data.Set("http_endpoint", details.HttpEndpoint)
	data.Set("deployment_type", details.DeploymentType)
	data.Set("num_nodes_default", details.NumNodesDefault)
	data.Set("num_additional_zookeeper_nodes", details.NumAdditionalZookeeperNodes)
	data.Set("servers", details.Servers)
	data.Set("zookeeper_ensemble", details.ZookeeperEnsemble)
	data.Set("specifications", map[string]string{
		"disk_space": details.Specifications.DiskSpace,
		"jvm_heap_memory": details.Specifications.JvmHeapMemory,
		"physical_memory": details.Specifications.PhysicalMemory,
	})
	data.Set("tags", details.Tag)

	return nil
}