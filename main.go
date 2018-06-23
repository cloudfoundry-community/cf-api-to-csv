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
	ansi "github.com/jhunt/go-ansi"
)

func main() {
	myClient := Client{}
	err := setup(&myClient)
	if err != nil {
		bailWith("err setting up client: %s", err)
	}
	orgs, err := getOrgs(myClient)
	if err != nil {
		bailWith("error getting orgs: %s", err)
	}
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("orgs before processing", orgs)

	//associate app creates with orgs
	for index, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.create&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app creates with orgs: %s", err)
		}
		orgs[index].associatedAppCreates = jsonResponse
	}

	//associate app starts with orgs
	for index, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.start&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app starts with orgs: %s", err)
		}
		orgs[index].associatedAppStarts = jsonResponse
	}

	//associate app updates with orgs
	for index, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.update&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app updates with orgs: %s", err)
		}
		orgs[index].associatedAppUpdates = jsonResponse
	}

	//associate space creates with orgs
	for index, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.space.create&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating space creates with orgs: %s", err)
		}
		orgs[index].associatedSpaceCreates = jsonResponse
	}
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("orgs after data processing:", orgs)
	fmt.Println("--------------------------------------------------------------")

	//list all service bindings with orgs
	// for _, org := range orgs {
	// 	jsonResponse, err := cfServiceBindingsRequest(myClient)
	// }

	//grab all the spaces
	spaces, err := getSpaces(myClient)
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	//associate app starts with spaces
	for index, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.start&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app starts with spaces: %s", err)
		}
		spaces[index].associatedAppStarts = jsonResponse
	}

	//associate app creates with spaces
	for index, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.create&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app creates with spaces: %s", err)
		}
		spaces[index].associatedAppCreates = jsonResponse
	}

	//associate app updates with spaces
	for index, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.update&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app updates with spaces: %s", err)
		}
		spaces[index].associatedAppUpdates = jsonResponse
	}
	// fmt.Println(spaces
	// for {
	// 	serve()
	// }
}

type space struct {
	name                      string
	guid                      string
	organizationGUID          string
	associatedAppCreates      cfEventsAPIResponse
	associatedAppStarts       cfEventsAPIResponse
	associatedAppUpdates      cfEventsAPIResponse
	associatedSpaceCreates    cfEventsAPIResponse
	associatedServiceBindings []*struct{}
}

func getSpaces(myClient Client) ([]space, error) {
	spaces := []space{}
	resp, err := myClient.doGetRequest("/v2/spaces")
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

type org struct {
	name                      string
	guid                      string
	associatedAppCreates      cfEventsAPIResponse
	associatedAppStarts       cfEventsAPIResponse
	associatedAppUpdates      cfEventsAPIResponse
	associatedSpaceCreates    cfEventsAPIResponse
	associatedServiceBindings []*struct{}
}

func getOrgs(myClient Client) ([]org, error) {
	orgs := []org{}
	resp, err := myClient.doGetRequest("/v2/organizations")
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

type cfEventsAPIResponse struct {
	Resources []struct {
		Metadata struct {
			GUID      string    `json:"guid"`
			URL       string    `json:"url"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
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
					Name              string `json:"name,omitempty"`
					Instances         int    `json:"instances,omitempty"`
					Memory            int    `json:"memory,omitempty"`
					State             string `json:"state,omitempty"`
					EnvironmentJSON   string `json:"environment_json,omitempty"`
					DockerImage       string `json:"docker_image,omitempty"`
					DockerCredentials string `json:"docker_credentials,omitempty"`
				} `json:"request"`
			} `json:"metadata"`
			SpaceGUID        string `json:"space_guid"`
			OrganizationGUID string `json:"organization_guid"`
		} `json:"entity"`
	} `json:"resources"`
}

func cfEventsRequest(myClient Client, endpoint string) (cfEventsAPIResponse, error) {
	resp, err := myClient.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return cfEventsAPIResponse{}, err
	}
	var responseyDoo cfEventsAPIResponse
	err = json.Unmarshal(body, &responseyDoo)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return cfEventsAPIResponse{}, err
	}
	return responseyDoo, nil
}

func setup(myClient *Client) error {
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
	*myClient = Client{
		authToken:  token,
		apiURL:     tmpURL,
		httpClient: &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
		cfClient:   goCFClient,
	}
	return nil
}

func bailWith(f string, a ...interface{}) {
	ansi.Fprintf(os.Stderr, fmt.Sprintf("@R{%s}\n", f), a...)
	os.Exit(1)
}
