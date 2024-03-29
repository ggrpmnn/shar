package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// parse the auth.log file, looking specifically for failed SSH attempts
func parseSSHAttempts(file *os.File) allEntries {
	attempts := make(allEntries, 0)

	// example auth log line for invalid entries: "Feb  1 19:02:48 grpi sshd[8749]: Invalid user pi from 202.120.42.141"
	rx := regexp.MustCompile(`(\w+\s+\d+) (\d{2}:\d{2}:\d{2}) [A-z]+ sshd\[\d+\]: Invalid user (.*) from ((?:\d{1,3}\.){3}\d{1,3}) port (\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := rx.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		// matches[0]=full string, [1]=date, [2]=time, [3]=user, [4]=IP, [5]=Port
		dateFound := false
		trimDate := strings.Join(strings.Fields(matches[1]), " ")
		for idx, dae := range attempts {
			if dae.Date == trimDate {
				dateFound = true
				jdx, ok := dae.exists(matches[4])
				if ok {
					debug("updating existing IP entry: '%s'", matches[4])
					tmp := dae.Entries[jdx]
					tmp.addUser(matches[3])
					dae.Entries[jdx] = tmp
				} else {
					debug("appending new IP: '%s'", matches[4])
					tmp := attempts[idx]
					tmp.Entries = append(tmp.Entries, authEntry{IP: matches[4], Count: 1, Users: []string{matches[3]}})
					attempts[idx] = tmp
				}
			}
		}
		if dateFound == false {
			debug("adding new date: '%s'", trimDate)
			newDate := datedAuthEntries{Date: trimDate, Entries: make([]authEntry, 0)}
			newDate.Entries = append(newDate.Entries, authEntry{IP: matches[4], Count: 1, Users: []string{matches[3]}})
			attempts = append(attempts, newDate)
		}
		dateFound = false
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return attempts
}
