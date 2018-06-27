package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

//Client is a struct containing all of the basic parts to make API requests to the Cloud Foundry API
type Client struct {
	authToken  string
	apiURL     *url.URL
	httpClient *http.Client
}

type space struct {
	Name                      string
	GUID                      string
	OrganizationGUID          string
	Apps                      []cfAPIResource
	AssociatedAppCreates      []cfAPIResource
	AssociatedAppStarts       []cfAPIResource
	AssociatedAppUpdates      []cfAPIResource
	AssociatedSpaceCreates    []cfAPIResource
	AssociatedServiceBindings []cfAPIResource
}

type org struct {
	Name                      string
	GUID                      string
	Apps                      []cfAPIResource
	AssociatedAppCreates      []cfAPIResource
	AssociatedAppStarts       []cfAPIResource
	AssociatedAppUpdates      []cfAPIResource
	AssociatedSpaceCreates    []cfAPIResource
	AssociatedServiceBindings []cfAPIResource
}

func (client *Client) setup() error {
	//old way with yaml parsing

	myConf, err := GrabCFCLIENV()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Printf("yaml config parsed: %v \n", *yamlConfig)

	token := myConf.AccessToken
	if err != nil {
		fmt.Println("error getting token")
		return err
	}
	tmpURL, err := url.Parse(myConf.Target)
	if err != nil {
		fmt.Println("error parsing config api address into URL")
		return err
	}
	client.authToken = token
	client.apiURL = tmpURL
	client.httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	return nil
}

func (client *Client) doGetRequest(path string) (*http.Response, error) {
	//fmt.Println("performing GET Request on path: " + client.apiURL.String() + path)
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

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return resp, nil
	}
	//if we hit this code we have a bad response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return nil, errors.New("bad response code in response, dumping body: " + string(bodyBytes))
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
	//fmt.Println("body received from get request", string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		return nil, err
	}
	//fmt.Println("using json from", in, "to build orgs")

	for index, resource := range in.Resources {
		orgs = append(orgs, org{})
		orgs[index].Name = resource.Entity.Name
		orgs[index].GUID = resource.Metadata.GUID
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
		spaces[index].Name = resource.Entity.Name
		spaces[index].OrganizationGUID = resource.Entity.OrganizationGUID
		spaces[index].GUID = resource.Metadata.GUID
	}
	return spaces, nil
}

func (client *Client) cfAPIRequest(endpoint string, returnStruct *cfAPIResponse) error {
	resp, err := client.doGetRequest(endpoint)
	//fmt.Println("got response from endpoint", endpoint)
	if err != nil {
		bailWith("err hitting cf endpoint: %s", err)
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
	//fmt.Println("returning json", returnStruct)
	return nil
}

func (client *Client) cfResourcesFromResponse(response cfAPIResponse) ([]cfAPIResource, error) {
	totalPages := response.TotalPages
	var resourceList []cfAPIResource
	for i := 0; i < totalPages; i++ {
		for _, resource := range response.Resources {
			resourceList = append(resourceList, resource)
		}
		//set the page into the next page
		err := client.cfAPIRequest(string(response.NextURL), &response)
		if err != nil {
			return nil, err
		}
	}
	return resourceList, nil
}
