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
	fmt.Printf("\n\n\n\n\n")
	fmt.Println("orgs before processing", orgs)
	fmt.Printf("\n\n\n\n\n")

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

	//get all apps based on org
	for index, org := range orgs {
		jsonResponse, err := cfAppsRequest(myClient, "/v2/apps?q=organization_guid:"+org.guid)
		if err != nil {
			bailWith("error associating apps with orgs: %s", err)
		}
		orgs[index].apps = jsonResponse.Resources
	}

	for index, org := range orgs {
		for index, app := range orgs[index].apps {
			jsonResponse, err := cfAppsRequest(myClient, "/v2/service_bindings?q=app_guid:"+orgs[index].apps[index].apps.guid)

		}
	}

	//get all service bindings based on apps by org

	fmt.Printf("\n\n\n\n\n")
	fmt.Println("orgs after data processing:", orgs)
	fmt.Printf("\n\n\n\n\n")

	//grab all the spaces
	spaces, err := getSpaces(myClient)
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	fmt.Printf("\n\n\n\n\n")
	fmt.Printf("spaces before data processing\n")
	fmt.Println(spaces)
	fmt.Printf("\n\n\n\n\n")

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

	fmt.Printf("\n\n\n\n\n")
	fmt.Printf("spaces after data processing\n")
	fmt.Println(spaces)
	fmt.Printf("\n\n\n\n\n")

	//get all apps based on spaces

	// get all service bindings based on apps by space

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
	apps                      []Resources
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

type cfApp struct {
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
type cfEventsAPIResponse struct {
	apps []cfApp
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

type Resources struct {
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
type cfAppsAPIResponse struct {
	Resources []struct {
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
	} `json:"resources"`
}

func cfAppsRequest(myClient Client, endpoint string) (cfAppsAPIResponse, error) {
	resp, err := myClient.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf apps endpoint: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return cfAppsAPIResponse{}, err
	}
	var responseyDoo cfAppsAPIResponse
	err = json.Unmarshal(body, &responseyDoo)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return cfAppsAPIResponse{}, err
	}
	return responseyDoo, nil
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

func cfServiceBindingsRequest(myClient Client, endpoint string) (cfServiceBindingsAPIResponse, error) {
	resp, err := myClient.doGetRequest(endpoint)
	if err != nil {
		bailWith("err hitting cf apps endpoint: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body")
		return cfServiceBindingsAPIResponse{}, err
	}
	var responseyDoo cfServiceBindingsAPIResponse
	err = json.Unmarshal(body, &responseyDoo)
	if err != nil {
		fmt.Println("error unmarshalling resp body into json")
		return cfServiceBindingsAPIResponse{}, err
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
