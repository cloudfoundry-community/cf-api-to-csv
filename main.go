package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aanelli/cf-metrics/webserver"
	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
	myClient := Client{}
	err := setup(&myClient)
	if err != nil {
		fmt.Println("error during client setup: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("looking for app creations")
	resp, err := myClient.doGetRequest("/v2/events?q=type:audit.app.create")
	if err != nil {
		fmt.Println("error looking for app creations: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("response: ", resp)

	fmt.Println("looking for app start events")
	resp, err = myClient.doGetRequest("/v2/events?q=type:audit.app.start")
	if err != nil {
		fmt.Println("error getting events data", err.Error())
		os.Exit(1)
	}
	fmt.Println("response: ", resp)

	fmt.Println("looking for space creations")
	resp, err = myClient.doGetRequest("/v2/events?q=type:audit.space.create")
	if err != nil {
		fmt.Println("error getting events data", err.Error())
		os.Exit(1)
	}
	fmt.Println("response: ", resp)

	fmt.Println("looking for app update events")
	resp, err = myClient.doGetRequest("/v2/events?q=type:audit.app.update")
	if err != nil {
		fmt.Println("error getting events data", err.Error())
		os.Exit(1)
	}
	fmt.Println("response: ", resp)

	fmt.Println("looking for app service bindings")
	resp, err = myClient.doGetRequest("/v2/service_bindings")
	if err != nil {
		fmt.Println("error getting events data", err.Error())
		os.Exit(1)
	}
	fmt.Println("response: ", resp)
	for {
		webserver.Serve()
	}
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
		fmt.Println("error creating new cfclient")
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
