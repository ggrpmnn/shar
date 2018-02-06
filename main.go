package main

import (
	"bufio"
	"flag"
	"regexp"

	"log"
	"os"
)

var (
	debugOn  bool
	filename string
	jsonOut  bool
)

func init() {
	flag.BoolVar(&debugOn, "d", false, "enable debug output")
	flag.BoolVar(&jsonOut, "j", false, "output results in JSON format")
	flag.StringVar(&filename, "f", "/var/log/auth.log", "auth log file to parse")
}

func main() {
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	debug("auth file loaded: %s", filename)

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
		for idx, dae := range attempts {
			if dae.date == matches[1] {
				dateFound = true
				jdx, ok := dae.exists(matches[4])
				if ok {
					debug("updating existing IP entry: '%s'", matches[4])
					tmp := dae.entries[jdx]
					tmp.addUser(matches[3])
					dae.entries[jdx] = tmp
				} else {
					debug("appending new IP: '%s'", matches[4])
					tmp := attempts[idx]
					tmp.entries = append(tmp.entries, authEntry{ip: matches[4], count: 1, users: []string{matches[3]}})
					attempts[idx] = tmp
				}
			}
		}
		if dateFound == false {
			debug("adding new date: '%s'", matches[1])
			newDate := datedAuthEntries{date: matches[1], entries: make([]authEntry, 0)}
			newDate.entries = append(newDate.entries, authEntry{ip: matches[4], count: 1, users: []string{matches[3]}})
			attempts = append(attempts, newDate)
		}
		dateFound = false
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	debug("finished parsing log file")

	// output parsed data to debug
	debug("data: %+v", attempts)

	if jsonOut {
		debug("outputting JSON")
		attempts.jsonPrint()
	} else {
		attempts.print()
	}
}

func debug(fmt string, a ...interface{}) {
	if debugOn {
		log.Printf(fmt, a...)
	}
}
