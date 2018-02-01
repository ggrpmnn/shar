package main

import (
	"bufio"
	"regexp"

	"log"
	"os"
)

var debugOn = false

func main() {
	file, err := os.Open("/var/log/auth.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	debug("auth file loaded")

	attempts := make(allEntries, 0)
	// example auth log line for invalid entries: "Feb  1 19:02:48 grpi sshd[8749]: Invalid user pi from 202.120.42.141"
	rx := regexp.MustCompile(`(\w+\s+\d+\s+\d{2}:\d{2}:\d{2})\s+grpi sshd\[\d+\]: Invalid user (.*) from (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := rx.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		// matches[0]=full string, [1]=date, [2]=user, [3]=IP
		dateFound := false
		for _, dae := range attempts {
			if dae.date == matches[1] {
				dateFound = true
				if dae.entries.exists(matches[3]) {
					tmpEntry := dae.entries[matches[3]]
					tmpEntry.addUser(matches[2])
					dae.entries[matches[3]] = tmpEntry
				} else {
					dae.entries[matches[3]] = authEntry{count: 1, users: []string{matches[2]}}
				}
			}
		}
		if dateFound == false {
			newDate := datedAuthEntries{date: matches[2], entries: make(authEntryList, 0)}
			newDate.entries[matches[3]] = authEntry{count: 1, users: []string{matches[2]}}
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	debug("finished parsing log file")

	attempts.print()
}

func debug(msg string) {
	if debugOn {
		log.Println(msg)
	}
}
