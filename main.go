package main

import (
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
	ansi "github.com/jhunt/go-ansi"
)

type cfAPIResponse struct {
	TotalResults int             `json:"total_results"`
	TotalPages   int             `json:"total_pages"`
	PrevURL      string          `json:"prev_url"`
	NextURL      string          `json:"next_url"`
	Resources    []cfAPIResource `json:"resources"`
}

func main() {
	var client Client
	err := client.setup()
	if err != nil {
		bailWith("err setting up client: %s", err)
	}

	//start up progress bars
	fmt.Println("getting orgs")
	orgs, err := client.getOrgs()
	if err != nil {
		bailWith("error getting orgs: %s", err)
	}

	//start up ui progress bars
	uiprogress.Start()

	//associate app creates with orgs "/v2/events?q=type:audit.app.create&q=organization_guid:"
	err = client.getEndpointData(orgs, FieldAppCreates, "/v2/events?q=type:audit.app.create&q=organization_guid:", "associating app creates with orgs")
	if err != nil {
		bailWith("error associating app creates with orgs: %s", err)
	}

	//associate app starts with orgs
	err = client.getEndpointData(orgs, FieldAppStarts, "/v2/events?q=type:audit.app.start&q=organization_guid:", "associating app starts with orgs")
	if err != nil {
		bailWith("error associating app starts with orgs: %s", err)
	}

	//associate app updates with orgs
	err = client.getEndpointData(orgs, FieldAppUpdates, "/v2/events?q=type:audit.app.update&q=organization_guid:", "associating app updates with orgs")
	if err != nil {
		bailWith("error associating app updates with orgs: %s", err)
	}

	//associate space creates with orgs
	err = client.getEndpointData(orgs, FieldSpaceCreates, "/v2/events?q=type:audit.space.create&q=organization_guid:", "associating space creates with orgs")
	if err != nil {
		bailWith("error associating space creates with orgs: %s", err)
	}

	//associate apps with orgs
	err = client.getEndpointData(orgs, FieldApps, "/v2/apps?q=organization_guid:", "associating apps with orgs")
	if err != nil {
		bailWith("error associating apps with orgs: %s", err)
	}
	//some app stuff for later?
	// for index, org := range orgs {
	// 	for index, app := range orgs[index].apps {
	// 		jsonResponse, err := cfAppsAPIRequest(client, "/v2/service_bindings?q=app_guid:"+orgs[index].apps[index].apps.guid)
	// 	}
	// }

	//get all service bindings based on apps by org
	//todo?

	//grab all the spaces
	fmt.Println("error getting spaces")
	spaces, err := client.getSpaces()
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	//associate app starts with spaces
	err = client.getEndpointData(spaces, FieldAppStarts, "/v2/events?q=type:audit.app.start&q=space_guid:", "associating app starts with spaces")
	if err != nil {
		bailWith("error associating app starts with spaces: %s", err)
	}

	//associate app creates with spaces
	err = client.getEndpointData(spaces, FieldAppCreates, "/v2/events?q=type:audit.app.create&q=space_guid:", "associating app creates with spaces")
	if err != nil {
		bailWith("error associating app creates with spaces: %s", err)
	}

	//associate app updates with spaces
	err = client.getEndpointData(spaces, FieldAppUpdates, "/v2/events?q=type:audit.app.update&q=space_guid:", "associating app updates with spaces")

	//get all apps based on spaces
	err = client.getEndpointData(spaces, FieldApps, "/v2/apps?q=space_guid:", "associating apps with spaces")
	if err != nil {
		bailWith("error associating apps with spaces: %s", err)
	}
	// get all service bindings based on apps by space

	// fmt.Println(spaces
	// for {
	// 	serve()
	// }
	// err = printAsJSON("orgs.json", orgs)
	// if err != nil {
	// 	bailWith("error writing orgs to file %s", err)
	// }
	// err = printAsJSON("spaces.json", spaces)
	// if err != nil {
	// 	bailWith("error writing spaces to file %s", err)
	// }

	err = printAsCSV("orgs.csv", orgs)
	if err != nil {
		bailWith("error writing orgs to csv %s", err)
	}

	err = printAsCSV("spaces.csv", spaces)
	if err != nil {
		bailWith("erorr writing spaces to csv %s", err)
	}
}

func bailWith(f string, a ...interface{}) {
	ansi.Fprintf(os.Stderr, fmt.Sprintf("@R{%s}\n", f), a...)
	os.Exit(1)
}

func sanitizeApps(v *cfAPIResource) {
	m, isMap := v.Entity.(map[string]interface{})
	if !isMap {
		panic("entity isn't a map!")
	}

	delete(m, "environment_json")
}

func sanitizeEvents(v *cfAPIResource) {
	m, isMap := v.Entity.(map[string]interface{})
	if !isMap {
		panic("entity isn't a map!")
	}

	meta, exists := m["metadata"]
	if !exists {
		panic("no metadata in events entity")
	}

	metaMap, isMap := meta.(map[string]interface{})
	if !isMap {
		panic("metadata isn't a map")
	}

	request, exists := metaMap["request"]
	if !exists {
		panic("no request in events metadata")
	}

	reqMap, isMap := request.(map[string]interface{})
	if !isMap {
		panic("request isn't a map")
	}

	delete(reqMap, "environment_json")
}
