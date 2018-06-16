package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient"
)

//Client is a struct containing all of the basic parts to make API requests to the Cloud Foundry API
type Client struct {
	authToken  string
	apiURL     *url.URL
	httpClient *http.Client
	cfClient   *cfclient.Client
}

func (c *Client) doGetRequest(path string) (*http.Response, error) {
	fmt.Println("performing GET Request on path: " + c.apiURL.String() + path)
	req, err := http.NewRequest("GET", c.apiURL.String()+path, nil)
	if err != nil {
		fmt.Println("error forming http GET request")
		return &http.Response{}, err
	}
	req.Header.Add("Authorization", c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error attempting http GET request")
		return &http.Response{}, err
	}

	dump, _ := httputil.DumpResponse(resp, true)
	fmt.Println("----------    dump of GET response body    ----------")
	fmt.Printf("%s\n", dump)
	return resp, nil

}
