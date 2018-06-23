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
	apps                      []appsAPIResource
	associatedAppCreates      []eventsAPIResource
	associatedAppStarts       []eventsAPIResource
	associatedAppUpdates      []eventsAPIResource
	associatedSpaceCreates    []eventsAPIResponse
	associatedServiceBindings []serviceBindingsAPIResource
}

type org struct {
	name                      string
	guid                      string
	apps                      []appsAPIResource
	associatedAppCreates      []eventsAPIResource
	associatedAppStarts       []eventsAPIResource
	associatedAppUpdates      []eventsAPIResource
	associatedSpaceCreates    []eventsAPIResponse
	associatedServiceBindings []serviceBindingsAPIResource
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
	orgs := []org{}
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
	spaces := []space{}
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

type cfServiceBindingsAPIResponse struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      interface{} `json:"prev_url"`
	NextURL      interface{} `json:"next_url"`
	Resources    []struct {
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
	} `json:"resources"`
}

func (client *Client) cfAPIRequest(endpoint string, returnStruct *struct{}) error {
	resp, err := client.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return err
	}
	err = json.Unmarshal(body, &returnStruct)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return err
	}
	return nil
}
