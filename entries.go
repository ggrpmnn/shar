package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// tracks the login attempts (per-IP, single day)
type authEntry struct {
	IP       string   `json:"ip"`
	Location string   `json:"location"`
	Count    int      `json:"count"`
	Users    []string `json:"usernames"`
}

// associates authEntryList objects with a particular date
type datedAuthEntries struct {
	Date    string      `json:"date"`
	Entries []authEntry `json:"entries"`
}

// slice for containing all dated entries
type allEntries []datedAuthEntries

// IPAPIResponse contains the response data to the IP API (https://ip-api.com/json) request
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

// makes a call to a IP-geolocation API, parses the data into a response struct and returns the result
func locateIP(ip string) (IPAPIResponse, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
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

// takes an IP API response struct and composes a location string using the data
func (iar IPAPIResponse) composeLocationString() string {
	return fmt.Sprintf("%s, %s, %s (%f, %f)", iar.City, iar.RegionName, iar.Country, iar.Latitude, iar.Longitude)
}

// adds the username to the list for the given IP AuthEntry struct
func (ae *authEntry) addUser(user string) {
	ae.Count++
	for _, un := range ae.Users {
		if un == user {
			return
		}
	}
	ae.Users = append(ae.Users, user)
}

// returns true if the IP string exists in the given map
func (dae *datedAuthEntries) exists(ip string) (int, bool) {
	for idx, ae := range dae.Entries {
		if ae.IP == ip {
			return idx, true
		}
	}
	return 0, false
}

func (ae allEntries) print() {
	for _, dae := range ae {
		color.Set(color.FgGreen, color.Bold)
		fmt.Println("Date: " + dae.Date)
		color.Unset()
		for _, ae := range dae.Entries {
			color.Set(color.FgBlue, color.Bold)
			fmt.Printf("IP: %s\n", ae.IP)
			color.Unset()
			color.Set(color.FgYellow)
			fmt.Print("Location: ")
			color.Unset()
			fmt.Println(ae.Location)
			color.Set(color.FgYellow)
			fmt.Print("Attempts: ")
			color.Unset()
			fmt.Println(ae.Count)
			color.Set(color.FgYellow)
			fmt.Print("Usernames: ")
			color.Unset()
			fmt.Println(strings.Join(ae.Users, ", "))
		}
		fmt.Println()
	}
}

func (ae allEntries) jsonPrint() {
	bytes, _ := json.MarshalIndent(ae, "", "    ")
	fmt.Println(string(bytes))
}
