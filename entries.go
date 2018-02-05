package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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
		color.Set(color.FgBlue, color.Bold)
		fmt.Printf("IP: %s\n", ip)
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
}

func (ael authEntryList) jsonPrint() {
	for ip, ae := range ael {
		fmt.Printf("\t\t\t\"ip\": \"%s\",", ip)
		fmt.Printf("\t\t\t\"count\": \"%d\",", ae.count)
		fmt.Printf("\t\t\t\"usernames\": \"[%s]\"", strings.Join(ae.users, ", "))
	}
}

func (dae datedAuthEntries) print() {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Date: " + dae.date)
	color.Unset()
	dae.entries.print()
	fmt.Println()
}

func (dae datedAuthEntries) jsonPrint() {
	fmt.Println("\t{")
	fmt.Printf("\t\t\"date\": \"%s\",\n", dae.date)
	fmt.Printf("\t\t\"entries\": {")
	dae.entries.jsonPrint()
	fmt.Println("\t\t}")
	fmt.Println("\t}")
}

func (ae allEntries) print() {
	for _, dae := range ae {
		dae.print()
	}
}

func (ae allEntries) jsonPrint() {
	fmt.Println("{")
	for _, dae := range ae {
		dae.jsonPrint()
	}
	fmt.Println("}")
}
