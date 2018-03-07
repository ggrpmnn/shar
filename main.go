package main

import (
	"flag"

	"log"
	"os"
)

var (
	debugOn   bool
	filename  string
	jsonOut   bool
	threshold int
	address   string
	user      string
)

func init() {
	flag.BoolVar(&debugOn, "d", false, "enables debug output")
	flag.BoolVar(&jsonOut, "j", false, "outputs results in JSON format")
	flag.StringVar(&filename, "f", "/var/log/auth.log", "indicates auth log file to parse")
	flag.IntVar(&threshold, "n", 0, "limits output to entries that have at least n login attempts")
	flag.StringVar(&address, "i", "", "limits output to entries that originate from the specified IP address")
	flag.StringVar(&user, "u", "", "limits output to entries that are logging in as the specified user")
}

func main() {
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	debug("auth file loaded: %s", filename)

	attempts := parseSSHAttempts(file)
	debug("finished parsing log file")

	// output parsed data to debug
	debug("data: %+v", attempts)

	// filter the results based on flags
	for idx := range attempts {
		// count filter
		if threshold > 0 {
			filtered := attempts[idx].Filter(func(ae authEntry) bool {
				return ae.Count >= threshold
			})
			attempts[idx].Entries = filtered
		}
		// IP
		if address != "" {
			filtered := attempts[idx].Filter(func(ae authEntry) bool {
				return ae.IP == address
			})
			attempts[idx].Entries = filtered
		}
		if user != "" {
			filtered := attempts[idx].Filter(func(ae authEntry) bool {
				for _, name := range ae.Users {
					if name == user {
						return true
					}
				}
				return false
			})
			attempts[idx].Entries = filtered
		}
	}
	debug("filtered data: %+v", attempts)

	if jsonOut {
		debug("outputting JSON")
		attempts.printJSON()
	} else {
		debug("outputting plaintext")
		attempts.print()
	}

	debug("operation complete")
}

func debug(fmt string, a ...interface{}) {
	if debugOn {
		log.Printf(fmt, a...)
	}
}
