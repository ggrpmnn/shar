package main

import (
	"bufio"
	"regexp"

	"log"
	"os"
)

var debugOn = false

func main() {
	file, err := os.Open("./auth.log.sample")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	debug("auth file loaded")

	attempts := make(allEntries, 0)
	// example auth log line for invalid entries: "Feb  1 19:02:48 grpi sshd[8749]: Invalid user pi from 202.120.42.141"
	rx := regexp.MustCompile(`(\w+\s+\d)+\s+(\d{2}:\d{2}:\d{2})\s+grpi sshd\[\d+\]: Invalid user (.*) from (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := rx.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		// matches[0]=full string, [1]=date, [2]=time, [3]=user, [4]=IP
		dateFound := false
		for _, dae := range attempts {
			if dae.date == matches[1] {
				dateFound = true
				if dae.entries.exists(matches[4]) {
					tmpEntry := dae.entries[matches[4]]
					tmpEntry.addUser(matches[3])
					dae.entries[matches[4]] = tmpEntry
				} else {
					dae.entries[matches[4]] = authEntry{count: 1, users: []string{matches[3]}}
				}
			}
		}
		if dateFound == false {
			newDate := datedAuthEntries{date: matches[1], entries: make(authEntryList, 0)}
			newDate.entries[matches[4]] = authEntry{count: 1, users: []string{matches[3]}}
			attempts = append(attempts, newDate)
		}
		dateFound = false
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
