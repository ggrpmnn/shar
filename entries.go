package main

import (
	"fmt"
	"strings"
)

// tracks the login attempts (per-IP, single day)
type authEntry struct {
	count int
	users []string
}

// maps IPs to authEntry attempts
type authEntryList map[string]authEntry

// associates authEntryList objects with a particular date
type datedAuthEntries struct {
	date    string
	entries authEntryList
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
func (ael authEntryList) exists(ip string) bool {
	_, ok := ael[ip]
	return ok
}

func (ael authEntryList) print() {
	for ip, ae := range ael {
		fmt.Printf("IP: %s, Attempt Count: %d, Users: %s\n", ip, ae.count, strings.Join(ae.users, ", "))
	}
}

func (dae datedAuthEntries) print() {
	fmt.Println("Date: " + dae.date)
	dae.entries.print()
}

func (ae allEntries) print() {
	for _, dae := range ae {
		dae.print()
	}
}
