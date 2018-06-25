package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
)

//Client is a struct containing all of the basic parts to make API requests to the Cloud Foundry API
type Client struct {
	authToken  string
	apiURL     *url.URL
	httpClient *http.Client
	cfClient   *cfclient.Client
}

type space struct {
	name                      string
	guid                      string
	organizationGUID          string
	apps                      []appsResource
	associatedAppCreates      []eventsResource
	associatedAppStarts       []eventsResource
	associatedAppUpdates      []eventsResource
	associatedSpaceCreates    []eventsResource
	associatedServiceBindings []serviceBindingsResource
}

type org struct {
	name                      string
	guid                      string
	apps                      []appsResource
	associatedAppCreates      []eventsResource
	associatedAppStarts       []eventsResource
	associatedAppUpdates      []eventsResource
	associatedSpaceCreates    []eventsResource
	associatedServiceBindings []serviceBindingsResource
}

func (client *Client) setup() error {
	yamlConfig, err := parseConfig("./config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Printf("yaml config parsed: %v \n", *yamlConfig)

	goCFConfig := &cfclient.Config{
		ApiAddress:        yamlConfig.APIAddress,
		Username:          yamlConfig.Username,
		Password:          yamlConfig.Password,
		SkipSslValidation: true,
	}

	goCFClient, err := cfclient.NewClient(goCFConfig)
	if err != nil {
		fmt.Println("error creating cfclient")
		return err
	}
	token, err := goCFClient.GetToken()
	if err != nil {
		fmt.Println("error getting token fron cfclient")
		return err
	}
	tmpURL, err := url.Parse(yamlConfig.APIAddress)
	if err != nil {
		fmt.Println("error parsing yaml config api address into URL")
		return err
	}
	client.authToken = token
	client.apiURL = tmpURL
	client.httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	client.cfClient = goCFClient

	return nil
}

func (client *Client) doGetRequest(path string) (*http.Response, error) {
	//fmt.Println("performing GET Request on path: " + c.apiURL.String() + path)
	req, err := http.NewRequest("GET", client.apiURL.String()+path, nil)
	if err != nil {
		fmt.Println("error forming http GET request")
		return &http.Response{}, err
	}
	req.Header.Add("Authorization", client.authToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println("error attempting http GET request")
		return &http.Response{}, err
	}

	// dump, _ := httputil.DumpResponse(resp, true)
	// fmt.Println("----------    dump of GET response body    ----------")
	// fmt.Printf("%s\n", dump)
	return resp, nil
}

func (client *Client) getOrgs() ([]org, error) {
	var orgs []org
	resp, err := client.doGetRequest("/v2/organizations")
	var in struct {
		Resources []struct {
			Metadata struct {
				GUID string `json:"guid"`
			} `json:"metadata"`
			Entity struct {
				Name string `json:"name"`
			} `json:"entity"`
		} `json:"resources"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		return nil, err
	}

	for index, resource := range in.Resources {
		orgs = append(orgs, org{})
		orgs[index].name = resource.Entity.Name
		orgs[index].guid = resource.Metadata.GUID
	}
	return orgs, nil
}

func (client *Client) getSpaces() ([]space, error) {
	var spaces []space
	resp, err := client.doGetRequest("/v2/spaces")
	var in struct {
		Resources []struct {
			Metadata struct {
				GUID string `json:"guid"`
			} `json:"metadata"`
			Entity struct {
				Name             string `json:"name"`
				OrganizationGUID string `json:"organization_guid"`
			} `json:"entity"`
		} `json:"resources"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		return nil, err
	}

	for index, resource := range in.Resources {
		spaces = append(spaces, space{})
		spaces[index].name = resource.Entity.Name
		spaces[index].organizationGUID = resource.Entity.OrganizationGUID
		spaces[index].guid = resource.Metadata.GUID
	}
	return spaces, nil
}

func (client *Client) cfServiceBindingsAPIRequest(endpoint string, returnStruct *serviceBindingsAPIResponse) error {
	resp, err := client.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return err
	}
	err = json.Unmarshal(body, returnStruct)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return err
	}
	return nil
}

func (client *Client) cfEventsAPIRequest(endpoint string, returnStruct *eventsAPIResponse) error {
	resp, err := client.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return err
	}
	err = json.Unmarshal(body, returnStruct)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return err
	}
	return nil
}

func (client *Client) cfAppsAPIRequest(endpoint string, returnStruct *appsAPIResponse) error {
	resp, err := client.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return err
	}
	err = json.Unmarshal(body, returnStruct)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return err
	}
	return nil
}

type appsResource struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name                     string      `json:"name"`
		Production               bool        `json:"production"`
		SpaceGUID                string      `json:"space_guid"`
		StackGUID                string      `json:"stack_guid"`
		Buildpack                interface{} `json:"buildpack"`
		DetectedBuildpack        interface{} `json:"detected_buildpack"`
		DetectedBuildpackGUID    interface{} `json:"detected_buildpack_guid"`
		EnvironmentJSON          interface{} `json:"environment_json"`
		Memory                   int         `json:"memory"`
		Instances                int         `json:"instances"`
		DiskQuota                int         `json:"disk_quota"`
		State                    string      `json:"state"`
		Version                  string      `json:"version"`
		Command                  interface{} `json:"command"`
		Console                  bool        `json:"console"`
		Debug                    interface{} `json:"debug"`
		StagingTaskID            interface{} `json:"staging_task_id"`
		PackageState             string      `json:"package_state"`
		HealthCheckHTTPEndpoint  string      `json:"health_check_http_endpoint"`
		HealthCheckType          string      `json:"health_check_type"`
		HealthCheckTimeout       interface{} `json:"health_check_timeout"`
		StagingFailedReason      interface{} `json:"staging_failed_reason"`
		StagingFailedDescription interface{} `json:"staging_failed_description"`
		Diego                    bool        `json:"diego"`
		DockerImage              interface{} `json:"docker_image"`
		DockerCredentials        struct {
			Username interface{} `json:"username"`
			Password interface{} `json:"password"`
		} `json:"docker_credentials"`
		PackageUpdatedAt     time.Time   `json:"package_updated_at"`
		DetectedStartCommand string      `json:"detected_start_command"`
		EnableSSH            bool        `json:"enable_ssh"`
		Ports                interface{} `json:"ports"`
		SpaceURL             string      `json:"space_url"`
		StackURL             string      `json:"stack_url"`
		RoutesURL            string      `json:"routes_url"`
		EventsURL            string      `json:"events_url"`
		ServiceBindingsURL   string      `json:"service_bindings_url"`
		RouteMappingsURL     string      `json:"route_mappings_url"`
	} `json:"entity"`
}

func (client *Client) cfAppsResourcesFromResponse(response appsAPIResponse) ([]appsResource, error) {
	totalPages := response.TotalPages
	var resourceList []appsResource
	for i := 0; i < totalPages; i++ {
		for _, resource := range response.Resources {
			resourceList = append(resourceList, resource)
		}
		//set the page into the next page
		err := client.cfAppsAPIRequest(string(response.NextURL), &response)
		if err != nil {
			return nil, err
		}
	}
	return resourceList, nil
}

type eventsResource struct {
	Metadata struct {
		GUID      string      `json:"guid"`
		URL       string      `json:"url"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Type      string    `json:"type"`
		Actor     string    `json:"actor"`
		ActorType string    `json:"actor_type"`
		ActorName string    `json:"actor_name"`
		Actee     string    `json:"actee"`
		ActeeType string    `json:"actee_type"`
		ActeeName string    `json:"actee_name"`
		Timestamp time.Time `json:"timestamp"`
		Metadata  struct {
			Request struct {
				Name                  string `json:"name"`
				Instances             int    `json:"instances"`
				Memory                int    `json:"memory"`
				State                 string `json:"state"`
				EnvironmentJSON       string `json:"environment_json"`
				DockerCredentialsJSON string `json:"docker_credentials_json"`
			} `json:"request"`
		} `json:"metadata"`
		SpaceGUID        string `json:"space_guid"`
		OrganizationGUID string `json:"organization_guid"`
	} `json:"entity"`
}

func (client *Client) cfEventsResourcesFromResponse(response eventsAPIResponse) ([]eventsResource, error) {
	totalPages := response.TotalPages
	resourceList := []eventsResource{}
	for i := 0; i < totalPages; i++ {
		for _, resource := range response.Resources {
			resourceList = append(resourceList, resource)
		}
		//set the page into the next page
		err := client.cfEventsAPIRequest(string(response.NextURL), &response)
		if err != nil {
			return nil, err
		}
	}
	return resourceList, nil
}

type serviceBindingsResource struct {
	Metadata struct {
		GUID      string      `json:"guid"`
		URL       string      `json:"url"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		AppGUID             string `json:"app_guid"`
		ServiceInstanceGUID string `json:"service_instance_guid"`
		Credentials         struct {
			CredsKey100 string `json:"creds-key-100"`
		} `json:"credentials"`
		BindingOptions struct {
		} `json:"binding_options"`
		GatewayData        interface{} `json:"gateway_data"`
		GatewayName        string      `json:"gateway_name"`
		SyslogDrainURL     interface{} `json:"syslog_drain_url"`
		AppURL             string      `json:"app_url"`
		ServiceInstanceURL string      `json:"service_instance_url"`
	} `json:"entity"`
}

func (client *Client) cfServiceBindingsResourcesFromResponse(response serviceBindingsAPIResponse) ([]serviceBindingsResource, error) {
	totalPages := response.TotalPages
	var resourceList []serviceBindingsResource
	for i := 0; i < totalPages; i++ {
		for _, resource := range response.Resources {
			resourceList = append(resourceList, resource)
		}
		//set the page into the next page
		err := client.cfServiceBindingsAPIRequest(string(response.NextURL), &response)
		if err != nil {
			return nil, err
		}
	}
	return resourceList, nil
}
