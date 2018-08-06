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
	"time"

	"github.com/gosuri/uiprogress"
)

//Client is a struct containing all of the basic parts to make API requests to the Cloud Foundry API
type Client struct {
	authToken    string
	refreshToken string
	uaaClient    string
	uaaSecret    string
	apiURL       *url.URL
	uaaURL       *url.URL
	httpClient   *http.Client
}

type cfAPIResource struct {
	Metadata cfAPIMetadata `json:"metadata"`
	Entity   interface{}   `json:"entity"`
}

type cfAPIMetadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type cfData struct {
	Name             string
	GUID             string
	OrganizationGUID string
	Apps             []cfAPIResource
	AppCreates       []cfAPIResource
	AppStarts        []cfAPIResource
	AppUpdates       []cfAPIResource
	SpaceCreates     []cfAPIResource
	ServiceBindings  []cfAPIResource
}
type DataField int

const (
	FieldApps DataField = iota
	FieldAppCreates
	FieldAppStarts
	FieldAppUpdates
	FieldSpaceCreates
	FieldServiceBindings
)

func (client *Client) setup() error {
	//old way with yaml parsing

	myConf, err := grabCFCLIENV()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Printf("yaml config parsed: %v \n", *yamlConfig)

	tmpURL, err := url.Parse(myConf.Target)
	if err != nil {
		fmt.Println("error parsing config api address into URL")
		return err
	}
	tmp2URL, err := url.Parse(myConf.UAAEndpoint)
	if err != nil {
		fmt.Println("error parsing uaa api address into URL")
		return err
	}

	client.authToken = myConf.AccessToken
	client.refreshToken = myConf.RefreshToken
	client.uaaClient = myConf.UAAClientID
	client.uaaSecret = myConf.UAAClientSecret
	client.apiURL = tmpURL
	client.uaaURL = tmp2URL
	client.httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	return nil
}

func (client *Client) refreshAccessToken() error {
	req, err := http.NewRequest("GET", client.uaaURL.String()+"/oauth/token", nil)
	if err != nil {
		fmt.Println("error forming http GET request")
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	myURLEncoding := url.Values{}
	myURLEncoding.Add("grant_type", "refresh_token")
	myURLEncoding.Add("refresh_token", client.refreshToken)
	myURLEncoding.Add("client_id", client.uaaClient)
	myURLEncoding.Add("client_secret", client.uaaSecret)
	req.URL.RawQuery = myURLEncoding.Encode()
	resp, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println("error attempting http GET request")
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("error: non 200 response code from uaa when attempting to refresh token")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read refresh response body: %s", err))
	}

	type refreshResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	contents := refreshResponse{}
	err = json.Unmarshal(b, &contents)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal refresh response JSON: %s", err))
	}
	client.authToken = fmt.Sprintf("bearer %s", contents.AccessToken)
	client.refreshToken = contents.RefreshToken

	return nil
}

func (client *Client) getOrgs() ([]cfData, error) {
	var orgs []cfData
	var resp cfAPIResponse
	err := client.cfAPIRequest("/v2/organizations", &resp)
	if err != nil {
		return nil, err
	}
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
		orgs = append(orgs, cfData{})
		orgs[index].Name = resource.Entity.Name
		orgs[index].GUID = resource.Metadata.GUID
	}
	return orgs, nil
}

func (client *Client) getSpaces() ([]cfData, error) {
	var spaces []cfData
	var resp cfAPIResponse
	err := client.cfAPIRequest("/v2/spaces", &resp)
	if err != nil {
		return nil, err
	}
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
		spaces = append(spaces, cfData{})
		spaces[index].Name = resource.Entity.Name
		spaces[index].OrganizationGUID = resource.Entity.OrganizationGUID
		spaces[index].GUID = resource.Metadata.GUID
	}
	return spaces, nil
}

func (client *Client) cfAPIRequest(endpoint string, returnStruct *cfAPIResponse, secondAttempt ...bool) error {

	//fmt.Println("performing GET Request on path: " + client.apiURL.String() + path)
	req, err := http.NewRequest("GET", client.apiURL.String()+endpoint, nil)
	if err != nil {
		fmt.Println("error forming http GET request")
		return err
	}
	req.Header.Add("Authorization", client.authToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println("error attempting http GET request")
		return err
	}

	if (resp.StatusCode == 401 || resp.StatusCode == 403) && len(secondAttempt) == 0 {
		err = client.refreshAccessToken()
		if err != nil {
			return fmt.Errorf("Error refreshing token: %s", err)
		}
		return client.cfAPIRequest(endpoint, returnStruct, true)
	}

	if resp.StatusCode >= 400 || resp.StatusCode <= 500 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New("bad response code in response, dumping body: " + string(bodyBytes))
	}

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

func (client *Client) getEndpointData(dataList []cfData, listToUpdate DataField, endpoint string, whatYoureDoing string) error {
	if len(whatYoureDoing) < 36 {
		//pad length to 36 chars to make it less ugly in the terminal
		for len(whatYoureDoing) < 36 {
			whatYoureDoing = whatYoureDoing + " "
		}
	}
	//add in terminal ui progress bars with comments
	bar := uiprogress.AddBar(len(dataList)).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf(whatYoureDoing)
	})

	//iterate over the list of orgs/spaces and ping the endpoint of choice
	for index, datapoint := range dataList {
		var response cfAPIResponse
		err := client.cfAPIRequest(endpoint+datapoint.GUID, &response)
		if err != nil {
			fmt.Println("error making cf api request", whatYoureDoing, ":", err)
			return err
		}

		//grab the data from said endpoint
		cfResources, err := client.cfResourcesFromResponse(response)
		if err != nil {
			fmt.Println("error getting resources out of api response:", err, "while attempting:", whatYoureDoing)
			return err
		}

		//add in the data in the chosen struct field
		switch listToUpdate {
		case FieldApps:
			for _, v := range cfResources {
				sanitizeApps(&v)
			}
			dataList[index].Apps = cfResources
		case FieldAppCreates:
			dataList[index].AppCreates = cfResources
		case FieldAppStarts:
			dataList[index].AppStarts = cfResources
		case FieldAppUpdates:
			for _, v := range cfResources {
				sanitizeEvents(&v)
			}
			dataList[index].AppUpdates = cfResources
		case FieldServiceBindings:
			dataList[index].ServiceBindings = cfResources
		case FieldSpaceCreates:
			dataList[index].SpaceCreates = cfResources
		}

		//update the terminal ui
		bar.Incr()
	}

	return nil
}

func (client *Client) cfResourcesFromResponse(response cfAPIResponse) ([]cfAPIResource, error) {
	totalPages := response.TotalPages
	var resourceList []cfAPIResource
	//iterate over the pages of the response until you get the full list of data
	for i := 0; i < totalPages; i++ {
		for _, resource := range response.Resources {
			resourceList = append(resourceList, resource)
		}
		//keep pinging the api until you get all of the data
		if i-1 < totalPages {
			//set the page into the next page
			err := client.cfAPIRequest(string(response.NextURL), &response)
			if err != nil {
				return nil, err
			}
		}
	}
	return resourceList, nil
}
