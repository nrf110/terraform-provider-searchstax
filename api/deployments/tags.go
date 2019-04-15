package deployments

import "fmt"

type TagRequest struct {
	Tags []string `json:"tags"`
}

type DeploymentTags struct {
	Deployment string `json:"deployment"`
	Tags []string `json:"tags"`
}

func (deployments *Deployments) CreateTags(deploymentId string, body TagRequest) (*[]DeploymentTags, error) {
	result := new([]DeploymentTags)

	_, err := deployments.sling.New().
		Post(fmt.Sprintf("/%s/deployments/", deploymentId)).
		BodyJSON(body).
		ReceiveSuccess(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (deployments *Deployments) DeleteTags(deploymentId string, body TagRequest) (*[]DeploymentTags, error) {
	result := new([]DeploymentTags)

	_, err := deployments.sling.New().
		Post(fmt.Sprintf("/%s/tags/delete/", deploymentId)).
		BodyJSON(body).
		ReceiveSuccess(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (deployments *Deployments) GetTags(deploymentId string) (*[]DeploymentTags, error) {
	result := new([]DeploymentTags)

	_, err := deployments.sling.New().
		Get(fmt.Sprintf("/%s/tags/", deploymentId)).
		ReceiveSuccess(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
