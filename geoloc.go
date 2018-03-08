package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/juju/ratelimit"
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
	client *http.Client
	URL    string

	// limiter is a rate limiter designed to keep requests from exceeding the IP API rate limit
	limiter *ratelimit.Bucket
}

const (
	ipapiRatePerSecond float64 = 150 / 60
	ipapiMaxRequests           = 125
)

// takes an IP API response struct and composes a location string using the data
func (iar *IPAPIResponse) composeLocationString() string {
	return fmt.Sprintf("%s, %s, %s (%f, %f)", iar.City, iar.Region, iar.Country, iar.Latitude, iar.Longitude)
}

func newIPAPIClient(url string) ipAPIClient {
	client := ipAPIClient{}

	client.client = http.DefaultClient
	client.URL = url
	client.limiter = ratelimit.NewBucketWithRate(ipapiRatePerSecond, ipapiMaxRequests)

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
	dur := iac.limiter.Take(1)
	time.Sleep(dur)
	return iac.client.Get(target)
}
