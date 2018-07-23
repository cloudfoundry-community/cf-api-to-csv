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

	//associate app creates with orgs
	numOfOrgs := len(orgs)
	uiprogress.Start()
	bar1 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Associating app creates with orgs")
	})
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

		bar1.Incr()
	}

	//associate app starts with orgs
	bar2 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Associating app starts with orgs")
	})
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
		bar2.Incr()
	}

	//associate app updates with orgs
	bar3 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Associating app updates with orgs")
	})
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
		orgs[index].AppUpdates = responseList
		bar3.Incr()
	}

	//associate space creates with orgs
	bar4 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Associating space creates with orgs")
	})
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

		orgs[index].SpaceCreates = responseList
		bar4.Incr()
	}

	//get all apps based on org
	fmt.Println("associating apps with orgs")
	//bar5 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed()
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
		//	bar5.Incr()
	}

	//some app stuff for later?
	// for index, org := range orgs {
	// 	for index, app := range orgs[index].apps {
	// 		jsonResponse, err := cfAppsAPIRequest(client, "/v2/service_bindings?q=app_guid:"+orgs[index].apps[index].apps.guid)

	// 	}
	// }

	//get all service bindings based on apps by org

	//grab all the spaces
	fmt.Println("error getting spaces")
	spaces, err := client.getSpaces()
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	//associate app starts with spaces
	fmt.Println("associating app starts with spaces")
	//bar6 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed()
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
		//	bar6.Incr()
	}

	//associate app creates with spaces
	fmt.Println("associating app creates with spaces")
	//bar7 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed()
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
		spaces[index].AppCreates = responseList
		//	bar7.Incr()
	}

	//associate app updates with spaces
	fmt.Println("associating app updates with spaces")
	//bar8 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed()
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

		spaces[index].AppUpdates = responseList
		//	bar8.Incr()
	}

	//get all apps based on spaces
	fmt.Println("associating apps with spaces")
	//bar9 := uiprogress.AddBar(numOfOrgs).AppendCompleted().PrependElapsed()
	for index, space := range spaces {
		var returnStruct cfAPIResponse
		err := client.cfAPIRequest("/v2/apps?q=space_guid:"+space.GUID, &returnStruct)
		if err != nil {
			bailWith("error associating apps with spaces: %s", err)
		}
		responseList, err := client.cfResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating apps with spaces: %s", err)
		}

		for _, v := range responseList {
			sanitizeApps(&v)
		}
		spaces[index].Apps = responseList
		//	bar9.Incr()
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
