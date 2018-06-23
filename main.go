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
	fmt.Println("----- printing orgs -----")
	fmt.Println(orgs)

	//associate app creates with orgs
	for _, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.create&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app creates with orgs: %s", err)
		}
		org.associatedAppCreates = append(org.associatedAppCreates, jsonResponse)
	}

	fmt.Println(" ----- printing orgs with app creates -----")
	fmt.Println(orgs)

	//associate app starts with orgs
	for _, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.start&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app starts with orgs: %s", err)
		}
		org.associatedAppStarts = append(org.associatedAppStarts, jsonResponse)
	}

	fmt.Println(" ----- printing orgs with app starts ----- ")
	fmt.Println(orgs)

	//associate app updates with orgs
	for _, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.update&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating app updates with orgs: %s", err)
		}
		org.associatedAppUpdates = append(org.associatedAppUpdates, jsonResponse)
	}

	fmt.Println(" ----- printing orgs with app updates ----- ")
	fmt.Println(orgs)

	//associate space creates with orgs
	for _, org := range orgs {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.space.create&q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating space creates with orgs: %s", err)
		}
		org.associatedSpaceCreates = append(org.associatedSpaceCreates, jsonResponse)
	}

	fmt.Println(" ----- printing orgs with space creates ----- ")
	fmt.Println(orgs)

	spaces, err := getSpaces(myClient)
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}
	fmt.Println("----- printing spaces -----")
	fmt.Println(spaces)

	//associate app starts with spaces
	for _, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.start&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app starts with spaces: %s", err)
		}
		space.associatedAppStarts = append(space.associatedAppStarts, jsonResponse)
	}

	//associate app creates with spaces
	for _, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.create&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app creates with spaces: %s", err)
		}
		space.associatedAppCreates = append(space.associatedAppCreates, jsonResponse)
	}
	fmt.Println(" ----- printing spaces with app creates ----- ")
	fmt.Println(spaces)

	//associate app updates with spaces
	for _, space := range spaces {
		jsonResponse, err := cfEventsRequest(myClient, "/v2/events?q=type:audit.app.update&q=space_guid:"+space.guid)
		if err != nil {
			bailWith("error associating app updates with spaces: %s", err)
		}
		space.associatedAppUpdates = append(space.associatedAppUpdates, jsonResponse)
	}
	fmt.Println(" ----- printing spaces with app updates ----- ")
	fmt.Println(spaces)
	// for {
	// 	serve()
	// }
}

type Space struct {
	name                      string
	guid                      string
	organizationGUID          string
	associatedAppCreates      []CFEventsAPIResponse
	associatedAppStarts       []CFEventsAPIResponse
	associatedAppUpdates      []CFEventsAPIResponse
	associatedSpaceCreates    []CFEventsAPIResponse
	associatedServiceBindings []*struct{}
}

func getSpaces(myClient Client) ([]Space, error) {
	spaces := []Space{}
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
		spaces = append(spaces, Space{})
		spaces[index].name = resource.Entity.Name
		spaces[index].organizationGUID = resource.Entity.OrganizationGUID
		spaces[index].guid = resource.Metadata.GUID
	}
	return spaces, nil
}

type Org struct {
	name                      string
	guid                      string
	associatedAppCreates      []CFEventsAPIResponse
	associatedAppStarts       []CFEventsAPIResponse
	associatedAppUpdates      []CFEventsAPIResponse
	associatedSpaceCreates    []CFEventsAPIResponse
	associatedServiceBindings []*struct{}
}

func getOrgs(myClient Client) ([]Org, error) {
	orgs := []Org{}
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
		orgs = append(orgs, Org{})
		orgs[index].name = resource.Entity.Name
		orgs[index].guid = resource.Metadata.GUID
	}
	return orgs, nil
}

type CFEventsAPIResponse struct {
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
					Name              string `json:"name"`
					Instances         int    `json:"instances"`
					Memory            int    `json:"memory"`
					State             string `json:"state"`
					EnvironmentJSON   string `json:"environment_json"`
					DockerImage       string `json:"docker_image"`
					DockerCredentials string `json:"docker_credentials"`
				} `json:"request"`
			} `json:"metadata"`
			SpaceGUID        string `json:"space_guid"`
			OrganizationGUID string `json:"organization_guid"`
		} `json:"entity"`
	} `json:"resources"`
}

func cfEventsRequest(myClient Client, endpoint string) (CFEventsAPIResponse, error) {
	resp, err := myClient.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf events endpoint: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return CFEventsAPIResponse{}, err
	}
	var responseyDoo CFEventsAPIResponse
	err = json.Unmarshal(body, &responseyDoo)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return CFEventsAPIResponse{}, err
	}
	return responseyDoo, nil
}

func setup(myClient *Client) error {
	yamlConfig, err := parseConfig("./config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("yaml config parsed: %v \n", *yamlConfig)

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
