package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// tracks the login attempts (per-IP, single day)
type authEntry struct {
	IP    string   `json:"ip"`
	Count int      `json:"count"`
	Users []string `json:"usernames"`
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

func (dae datedAuthEntries) print() {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Date: " + dae.Date)
	color.Unset()
	for _, ae := range dae.Entries {
		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("IP: %s\n", ae.IP)
		color.Unset()
		color.Set(color.FgYellow)
		fmt.Print("Num. attempts: ")
		color.Unset()
		fmt.Printf("%d\n", ae.Count)
		color.Set(color.FgYellow)
		fmt.Print("Usernames: ")
		color.Unset()
		fmt.Printf("%s\n", strings.Join(ae.Users, ", "))
	}
	fmt.Println()
}

func (ae allEntries) print() {
	for _, dae := range ae {
		dae.print()
	}
}

func (ae allEntries) jsonPrint() {
	_ = json.NewEncoder(os.Stdout).Encode(&ae)
}
