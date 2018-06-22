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

type Org struct {
	name                      string
	guid                      string
	associatedAppCreates      []createAppResponse
	associatedAppStarts       []*struct{}
	associatedAppUpdates      []*struct{}
	associatedSpaceCreates    []*struct{}
	associatedServiceBindings []*struct{}
}

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

	orgs, err = associateAppCreatesWithOrgs(orgs, myClient)
	if err != nil {
		bailWith("error associating app creates with orgs: %s", err)
	}
	fmt.Println(" ----- printing orgs with app creates -----")
	fmt.Println(orgs)
	for _, org := range orgs {
		fmt.Println(" ----- heres ya gross json Alex ----- ")
		for _, appCreate := range org.associatedAppCreates {
			fmt.Println(appCreate)
		}
	}

}

//loop through all of the resources
// for index, resource := range response.Resources {
// 	if val, ok := orgs[resource.Entity.OrganizationGUID]; !ok {
// 		//orgs doesn't contain the current org so add it
// 		orgs[resource.Entity.OrganizationGUID] = &resource
// 	}
// 	if val, ok = orgs[resource.Entity.SpaceGUID]; !ok {
// 		spaces[resource.Entity.SpaceGUID] = &resource
// 	}
// 	resource.Entity.SpaceGUID
// 	resource.Entity.ActorName
// }

// resp, err = myClient.doGetRequest("/v2/events?q=type:audit.app.start")
// if err != nil {
// 	bailWith("err getting app starts: %s", err)
// }
// body, err = ioutil.ReadAll(resp.Body)
// if err != nil {
// 	bailWith("err reading resp body: %s", err)
// }
// fmt.Println(string(body))

// resp, err = myClient.doGetRequest("/v2/events?q=type:audit.space.create")
// if err != nil {
// 	bailWith("err getting space creations: %s", err)
// }
// body, err = ioutil.ReadAll(resp.Body)
// if err != nil {
// 	bailWith("err reading resp body: %s", err)
// }
// fmt.Println(string(body))

// resp, err = myClient.doGetRequest("/v2/events?q=type:audit.app.update")
// if err != nil {
// 	bailWith("err getting app updates: %s", err)
// }
// body, err = ioutil.ReadAll(resp.Body)
// if err != nil {
// 	bailWith("err reading resp body: %s", err)
// }
// fmt.Println(string(body))

// resp, err = myClient.doGetRequest("/v2/service_bindings")
// if err != nil {
// 	bailWith("err getting service bindings %s", err)
// }
// body, err = ioutil.ReadAll(resp.Body)
// if err != nil {
// 	bailWith("err reading resp body: %s", err)
// }
// fmt.Println(string(body))

// for {
// 	serve()
// }

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

type createAppResponse struct {
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

func associateAppCreatesWithOrgs(orgs []Org, myClient Client) ([]Org, error) {
	for index, org := range orgs {
		resp, err := myClient.doGetRequest("/v2/events?q=type:audit.app.create&q=organization_guid:" + org.name)
		if err != nil {
			bailWith("err getting app creates: %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error reading resp body")
			return nil, err
		}
		var responseyDoo createAppResponse
		err = json.Unmarshal(body, &responseyDoo)
		if err != nil {
			fmt.Println("error unmarshalling resp body into json")
			return nil, err
		}
		orgs[index].associatedAppCreates = append(orgs[index].associatedAppCreates, responseyDoo)
	}
	return orgs, nil
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
