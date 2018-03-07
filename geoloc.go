package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// IPAPIResponse contains the response data from an IP API (https://ip-api.com/json) request
type IPAPIResponse struct {
	Status       string  `json:"status"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"countryCode"`
	Region       string  `json:"region"`
	RegionName   string  `json:"regionName"`
	City         string  `json:"city"`
	ZipCode      string  `json:"zip"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lon"`
	TimeZone     string  `json:"timezone"`
	ISP          string  `json:"isp"`
	Organization string  `json:"org"`
	As           string  `json:"as"`
	QueryIP      string  `json:"query"`
}

type ipAPIClient struct {
	client   *http.Client
	URL      string
	liveReqs int
}

const rateLimitPerMin = 150

// takes an IP API response struct and composes a location string using the data
func (iar *IPAPIResponse) composeLocationString() string {
	return fmt.Sprintf("%s, %s, %s (%f, %f)", iar.City, iar.Region, iar.Country, iar.Latitude, iar.Longitude)
}

func newIPAPIClient(url string) ipAPIClient {
	client := ipAPIClient{}

	client.client = http.DefaultClient
	client.URL = url

	return client
}

// makes a call to a IP-geolocation API, parses the data into a response struct and returns the result
func (iac *ipAPIClient) locateIP(ip string) (IPAPIResponse, error) {
	resp, err := iac.Get(iac.URL + ip)
	if err != nil {
		return IPAPIResponse{}, err
	}
	defer resp.Body.Close()

	bData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPAPIResponse{}, err
	}

	location := IPAPIResponse{}
	err = json.Unmarshal(bData, &location)
	if err != nil {
		return IPAPIResponse{}, err
	}

	return location, nil
}

func (iac *ipAPIClient) Get(target string) (*http.Response, error) {
	// naive attempt at client-side rate limiting; refactor this with a better strategy later
	iac.liveReqs++
	if iac.liveReqs >= rateLimitPerMin {
		time.Sleep(1 * time.Minute)
		iac.liveReqs = 0
	}

	return iac.client.Get(target)
}
