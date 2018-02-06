package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// tracks the login attempts (per-IP, single day)
type authEntry struct {
	ip    string
	count int
	users []string
}

// associates authEntryList objects with a particular date
type datedAuthEntries struct {
	date    string
	entries []authEntry
}

// slice for containing all dated entries
type allEntries []datedAuthEntries

// adds the username to the list for the given IP AuthEntry struct
func (ae *authEntry) addUser(user string) {
	ae.count++
	for _, un := range ae.users {
		if un == user {
			return
		}
	}
	ae.users = append(ae.users, user)
}

// returns true if the IP string exists in the given map
func (dae *datedAuthEntries) exists(ip string) (int, bool) {
	for idx, ae := range dae.entries {
		if ae.ip == ip {
			return idx, true
		}
	}
	return 0, false
}

func (dae datedAuthEntries) print() {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Date: " + dae.date)
	color.Unset()
	for _, ae := range dae.entries {
		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("IP: %s\n", ae.ip)
		color.Unset()
		color.Set(color.FgYellow)
		fmt.Print("Num. attempts: ")
		color.Unset()
		fmt.Printf("%d\n", ae.count)
		color.Set(color.FgYellow)
		fmt.Print("Usernames: ")
		color.Unset()
		fmt.Printf("%s\n", strings.Join(ae.users, ", "))
	}
	fmt.Println()
}

func (ae allEntries) print() {
	for _, dae := range ae {
		dae.print()
	}
}

func (ae allEntries) jsonPrint() {
	for _, dae := range ae {
		bytes, _ := json.Marshal(dae.entries)
		fmt.Println(string(bytes))
	}
}
