package main

import (
	"fmt"
	"os"

	ansi "github.com/jhunt/go-ansi"
)

func main() {
	var client Client
	err := client.setup()
	if err != nil {
		bailWith("err setting up client: %s", err)
	}
	orgs, err := client.getOrgs()
	if err != nil {
		bailWith("error getting orgs: %s", err)
	}

	//associate app creates with orgs
	for index, org := range orgs {
		var response cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.create&q=organization_guid:"+org.GUID, &response)
		if err != nil {
			bailWith("error associating app creates with orgs: %s", err)
		}
		resourceList, err := client.cfResourcesFromResponse(response)
		if err != nil {
			bailWith("error getting resources out of api response %s", err)
		}
		orgs[index].AppCreates = resourceList
	}

	//associate app starts with orgs
	for index, org := range orgs {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.start&q=organization_guid:"+org.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating app starts with orgs: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error getting resources out of api resp %s", err)
		}
		orgs[index].AppStarts = responseList
	}

	//associate app updates with orgs
	for index, org := range orgs {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.update&q=organization_guid:"+org.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating app updates with orgs: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app updates with orgs %s", err)
		}

		for _, v := range responseList {
			sanitizeEvents(&v)
		}
		orgs[index].AppStarts = responseList
	}

	//associate space creates with orgs
	for index, org := range orgs {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.space.create&q=organization_guid:"+org.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating space creates with orgs: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating space creates with orgs %s", err)
		}

		orgs[index].AppStarts = responseList
	}

	//get all apps based on org
	for index, org := range orgs {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/apps?q=organization_guid:"+org.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating apps with orgs: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating apps with orgs: %s", err)
		}

		for _, v := range responseList {
			sanitizeApps(&v)
		}
		orgs[index].Apps = responseList
	}

	//some app stuff for later?
	// for index, org := range orgs {
	// 	for index, app := range orgs[index].apps {
	// 		jsonResponse, err := cfAppsAPIRequest(client, "/v2/service_bindings?q=app_guid:"+orgs[index].apps[index].apps.guid)

	// 	}
	// }

	//get all service bindings based on apps by org

	//grab all the spaces
	spaces, err := client.getSpaces()
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	//associate app starts with spaces
	for index, space := range spaces {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.start&q=space_guid:"+space.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating app starts with spaces: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app starts with spaces %s", err)
		}
		spaces[index].AppStarts = responseList
	}

	//associate app creates with spaces
	for index, space := range spaces {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.create&q=space_guid:"+space.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating app creates with spaces: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app creates with spaces %s", err)
		}
		spaces[index].AppStarts = responseList
	}

	//associate app updates with spaces
	for index, space := range spaces {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/events?q=type:audit.app.update&q=space_guid:"+space.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating app updates with spaces: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app updates with spaces %s", err)
		}

		for _, v := range responseList {
			sanitizeEvents(&v)
		}

		spaces[index].AppStarts = responseList
	}

	//get all apps based on spaces

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
