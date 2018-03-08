package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// tracks the login attempts (per-IP, single day)
type authEntry struct {
	IP       string        `json:"ip"`
	Count    int           `json:"count"`
	Users    []string      `json:"usernames"`
	Location IPAPIResponse `json:"location"`
}

// associates authEntryList objects with a particular date
type datedAuthEntries struct {
	Date    string      `json:"date"`
	Entries []authEntry `json:"entries"`
}

// slice for containing all dated entries
type allEntries []datedAuthEntries

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

// maps a function to the entries underneath a datedAuthEntries struct
func (dae *datedAuthEntries) apply(f func(authEntry) authEntry) []authEntry {
	applied := make([]authEntry, 0)
	for _, ae := range dae.Entries {
		applied = append(applied, f(ae))
	}
	return applied
}

// filters the entries based on certain functional criteria
func (dae *datedAuthEntries) filter(f func(authEntry) bool) []authEntry {
	filtered := make([]authEntry, 0)
	for _, ae := range dae.Entries {
		if f(ae) {
			filtered = append(filtered, ae)
		}
	}
	return filtered
}

func (ae allEntries) print() {
	for idx, dae := range ae {
		// don't print the date if there are no entries
		if len(dae.Entries) > 0 {
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
				fmt.Println(ae.Location.composeLocationString())
				color.Set(color.FgYellow)
				fmt.Print("Attempts: ")
				color.Unset()
				fmt.Println(ae.Count)
				color.Set(color.FgYellow)
				fmt.Print("Usernames: ")
				color.Unset()
				fmt.Println(strings.Join(ae.Users, ", "))
			}
			// don't print a newline after the last date
			if idx != len(ae)-1 {
				fmt.Println()
			}
		}
	}
}

func (ae allEntries) printJSON() {
	bytes, _ := json.MarshalIndent(ae, "", "    ")
	fmt.Println(string(bytes))
}
