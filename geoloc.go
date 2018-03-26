package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

// ipAPIResponse contains the response data from an IP API (https://ip-api.com/json) request
type ipAPIResponse struct {
	Status       string  `json:"-"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"-"`
	Region       string  `json:"region"`
	RegionName   string  `json:"-"`
	City         string  `json:"city"`
	ZipCode      string  `json:"zip"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lon"`
	TimeZone     string  `json:"timezone"`
	ISP          string  `json:"isp"`
	Organization string  `json:"org"`
	As           string  `json:"-"`
	QueryIP      string  `json:"-"`
}

type ipAPIClient struct {
	client *http.Client
	URL    string

	// limiter is a rate limiter designed to keep requests from exceeding the IP API rate limit
	limiter *ratelimit.Bucket
}

const (
	ipapiRatePerSecond float64 = 2
	ipapiMaxRequests           = 1
)

// takes an IP API response struct and composes a location string using the data
func (iar *ipAPIResponse) composeLocationString() string {
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
func (iac *ipAPIClient) locateIP(ip string) (ipAPIResponse, error) {
	resp, err := iac.Get(iac.URL + "/json/" + ip)
	if err != nil {
		return ipAPIResponse{}, err
	}
	defer resp.Body.Close()

	bData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ipAPIResponse{}, err
	}

	location := ipAPIResponse{}
	err = json.Unmarshal(bData, &location)
	if err != nil {
		return ipAPIResponse{}, err
	}

	return location, nil
}

func (iac *ipAPIClient) Get(target string) (*http.Response, error) {
	dur := iac.limiter.Take(1)
	time.Sleep(dur)
	return iac.client.Get(target)
}
