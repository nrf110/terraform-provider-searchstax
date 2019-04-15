package deployments

import (
	"github.com/dghubble/sling"
	"time"
)

type Deployments struct {
	sling *sling.Sling
}

type CreateDeploymentRequest struct {
	Name                  string `json:"name"`
	Application           string `json:"application"`
	ApplicationVersion    string `json:"application_version"`
	TerminationLock       bool   `json:"termination_lock"`
	PlanType              string `json:"plan_type"`
	Plan                  string `json:"plan"`
	RegionID              string `json:"region_id"`
	CloudProviderID       string `json:"cloud_provider_id"`
	NumAdditionalAppNodes int    `json:"num_additional_app_nodes"`
}

type DeploymentDetails struct {
	Name                        string   `json:"name"`
	UID                         string   `json:"uid"`
	Application                 string   `json:"application"`
	ApplicationVersion          string   `json:"application_version"`
	Tier                        string   `json:"tier"`
	HttpEndpoint                string   `json:"http_endpoint"`
	Status                      string   `json:"status"`
	ProvisionState              string   `json:"provision_state"`
	TerminationLock             bool     `json:"termination_lock"`
	PlanType                    string   `json:"plan_type"`
	Plan                        string   `json:"plan"`
	RegionID                    string   `json:"region_id"`
	CloudProvider               string   `json:"cloud_provider"`
	CloudProviderID             string   `json:"cloud_provider_id"`
	NumAdditionalAppNodes       int      `json:"num_additional_app_nodes"`
	DeploymentType              string   `json:"deployment_type"`
	NumNodesDefault             int      `json:"num_nodes_default"`
	NumAdditionalZookeeperNodes int      `json:"num_additional_zookeeper_nodes"`
	Servers                     []string `json:"servers"`
	ZookeeperEnsemble           string   `json:"zookeeper_ensemble"`
	Tag                         []string `json:"tag"`
	Specifications              struct {
		DiskSpace      string `json:"disk_space"`
		JvmHeapMemory  string `json:"jvm_heap_memory"`
		PhysicalMemory string `json:"physical_memory"`
	} `json:"specifications"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func New(client *sling.Sling) *Deployments {
	return &Deployments{
		sling: client.New().
			Path("/deployment"),
	}
}

func (deployments *Deployments) waitUntilReady(uid string) (*DeploymentDetails, error) {
	started := time.Now()

	for {
		details, err := deployments.Get(uid)
		if err != nil {
			return nil, err
		}
		if details.ProvisionState == "Done" && details.Status == "Running" {
			return details, nil
		}
		// TODO: How do we determine if provisioning failed?
		if int(time.Since(started).Minutes()) > 65 {
			// TODO: timeout
		}
		time.Sleep(30 * time.Second)
	}
}

func (deployments *Deployments) Get(uid string) (*DeploymentDetails, error) {
	details := new(DeploymentDetails)

	_, err := deployments.sling.New().
		Get(uid).
		ReceiveSuccess(details)

	if err != nil {
		return nil, err
	}

	return details, nil
}

func (deployments *Deployments) Create(body CreateDeploymentRequest) (*DeploymentDetails, error) {
	details := new(DeploymentDetails)

	_, err := deployments.sling.New().
		Post("/").
		BodyJSON(body).
		ReceiveSuccess(details)

	if err != nil {
		return nil, err
	}

	return deployments.waitUntilReady(details.UID)
}

func (deployments *Deployments) Delete(uid string) error {
	errorResponse := new(ErrorResponse)

	_, err := deployments.sling.New().
		Delete(uid).
		Receive(nil, errorResponse)

	return err
}
