package apiclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WorkloadCreateRequest struct {
	EnvironmentName string                        `json:"-"`
	Name            string                        `json:"name"`
	Deployments     []WorkloadAutoscaleDeployment `json:"deployments"`
	Image           string                        `json:"image"`
	Specs           string                        `json:"specs"`
	Type            string                        `json:"type"`
	//Optional fields
	AddAnyCastIpAddress           bool                          `json:"addAnyCastIpAddress,omitempty"`
	AddImagePullCredentialsOption bool                          `json:"addImagePullCredentialsOption,omitempty"`
	Commands                      []string                      `json:"commands,omitempty"`
	ContainerEmail                string                        `json:"containerEmail,omitempty"`
	ContainerPassword             string                        `json:"containerPassword,omitempty"`
	ContainerServer               string                        `json:"containerServer,omitempty"`
	ContainerUsername             string                        `json:"containerUsername,omitempty"`
	EnvironmentVariables          []WorkloadEnvironmentVariable `json:"environmentVariables,omitempty"`
	PersistentStorage             []WorkloadPersistentStorage   `json:"persistentStorage,omitempty"`
	Ports                         []WorkloadPort                `json:"ports,omitempty"`
	SecretEnvironmentVariables    []WorkloadEnvironmentVariable `json:"secretEnvironmentVariables,omitempty"`
	Slug                          string                        `json:"slug,omitempty"`
}

//GetWorkloads Get workloads in account
func (c *Client) GetWorkloads(environmentName string) ([]Workload, error) {
	request, err := http.NewRequest("GET",
		CoxEdgeAPIBase+"/services/"+CoxEdgeServiceCode+"/"+environmentName+"/workloads",
		nil,
	)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	var wrappedAPIStruct WrappedWorkloads
	err = json.Unmarshal(respBytes, &wrappedAPIStruct)
	if err != nil {
		return nil, err
	}
	return wrappedAPIStruct.Data, nil
}

//GetWorkload Get workload in account by id
func (c *Client) GetWorkload(environmentName string, id string) (*Workload, error) {
	//Create the request
	request, err := http.NewRequest("GET",
		CoxEdgeAPIBase+"/services/"+CoxEdgeServiceCode+"/"+environmentName+"/workloads/"+id,
		nil,
	)
	if err != nil {
		return nil, err
	}

	//Execute request
	respBytes, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	//Unmarshal, unwrap, and return
	var wrappedAPIStruct WrappedWorkload
	err = json.Unmarshal(respBytes, &wrappedAPIStruct)
	if err != nil {
		return nil, err
	}
	return &wrappedAPIStruct.Data, nil
}

//CreateWorkload Create the workload
func (c *Client) CreateWorkload(newWorkload WorkloadCreateRequest) (*TaskStatusResponse, error) {
	//Marshal the request
	jsonBytes, err := json.Marshal(newWorkload)
	if err != nil {
		return nil, err
	}
	//Wrap bytes in reader
	bReader := bytes.NewReader(jsonBytes)
	//Create the request
	request, err := http.NewRequest("POST",
		CoxEdgeAPIBase+"/services/"+CoxEdgeServiceCode+"/"+newWorkload.EnvironmentName+"/workloads",
		bReader,
	)
	request.Header.Set("Content-Type", "application/json")
	//Execute request
	respBytes, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}
	//Return struct
	var wrappedAPIStruct TaskStatusResponse
	err = json.Unmarshal(respBytes, &wrappedAPIStruct)
	if err != nil {
		return nil, err
	}
	return &wrappedAPIStruct, nil
}

//UpdateWorkload Update a workload
func (c *Client) UpdateWorkload(workloadId string, newWorkload WorkloadCreateRequest) (*TaskStatusResponse, error) {
	//Marshal the request
	jsonBytes, err := json.Marshal(newWorkload)
	if err != nil {
		return nil, err
	}
	//Wrap bytes in reader
	bReader := bytes.NewReader(jsonBytes)
	//Create the request
	request, err := http.NewRequest("PUT",
		CoxEdgeAPIBase+"/services/"+CoxEdgeServiceCode+"/"+newWorkload.EnvironmentName+"/workloads/"+workloadId,
		bReader,
	)
	request.Header.Set("Content-Type", "application/json")
	//Execute request
	respBytes, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}
	//Return struct
	var wrappedAPIStruct TaskStatusResponse
	err = json.Unmarshal(respBytes, &wrappedAPIStruct)
	if err != nil {
		return nil, err
	}
	return &wrappedAPIStruct, nil
}

//DeleteWorkload Delete workload in account by id
func (c *Client) DeleteWorkload(environmentName string, id string) error {
	//Create the request
	request, err := http.NewRequest("DELETE",
		CoxEdgeAPIBase+"/services/"+CoxEdgeServiceCode+"/"+environmentName+"/workloads/"+id,
		nil,
	)
	if err != nil {
		return err
	}

	//Execute request
	_, err = c.doRequest(request)
	if err != nil {
		return err
	}
	return nil
}
