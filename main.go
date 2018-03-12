package main

import (
	"flag"
	"log"
	"os"
	"regexp"
)

var (
	debugOn   bool
	filename  string
	jsonOut   bool
	threshold int
	address   string
	user      string
	locale    string
	date      string
)

func init() {
	flag.BoolVar(&debugOn, "b", false, "enables debug output")
	flag.BoolVar(&jsonOut, "j", false, "outputs results in JSON format")
	flag.StringVar(&filename, "f", "/var/log/auth.log", "indicates auth log file to parse")
	flag.IntVar(&threshold, "n", 0, "limits output to entries that have at least n login attempts")
	flag.StringVar(&address, "i", "", "limits output to entries that originate from the specified IP address")
	flag.StringVar(&user, "u", "", "limits output to entries that are logging in as the specified user")
	flag.StringVar(&locale, "l", "", "limits output to entries that match the specified location string")
	flag.StringVar(&date, "d", "", "limits output to entries from the specified date (ex. Jan 1)")
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
	debug("raw data: %+v", attempts)

	// filter the results based on flags
	// start by filtering on dates
	if date != "" {
		attempts = applyDateFilter(attempts)
		if attempts == nil {
			log.Printf("found no date matching supplied filter; exiting")
			return
		}
	}
	applyEntryFilters(attempts)
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

// filters the output down to the specified date
func applyDateFilter(dae []datedAuthEntries) []datedAuthEntries {
	for _, day := range dae {
		if day.Date == date {
			return append([]datedAuthEntries{}, day)
		}
	}
	return nil
}

// filter the results for each date's entries based on the provided command-line flags; order of filtering is not
// particularly important (generally, we try to apply the strictest filters first), however,
// the location filter should be last in order to make the fewest requests possible to the IP-API
func applyEntryFilters(dae []datedAuthEntries) {
	for idx := range dae {
		// count filter
		if threshold > 0 {
			filtered := dae[idx].filter(func(ae authEntry) bool {
				return ae.Count >= threshold
			})
			dae[idx].Entries = filtered
		}
		// IP address filter
		if address != "" {
			filtered := dae[idx].filter(func(ae authEntry) bool {
				return ae.IP == address
			})
			dae[idx].Entries = filtered
		}
		// username filter
		if user != "" {
			filtered := dae[idx].filter(func(ae authEntry) bool {
				for _, name := range ae.Users {
					if name == user {
						return true
					}
				}
				return false
			})
			dae[idx].Entries = filtered
		}
		// get IP locations in order to apply location filter
		iac := newIPAPIClient("http://ip-api.com/")
		dae[idx].Entries = dae[idx].apply(func(ae authEntry) authEntry {
			debug("making API request for IP '%s'", ae.IP)
			location, err := iac.locateIP(ae.IP)
			if err != nil {
				log.Printf("error getting location data for IP '%s': %s", ae.IP, err.Error())
			}
			ae.Location = location
			return ae
		})
		// location filter
		if locale != "" {
			filtered := dae[idx].filter(func(ae authEntry) bool {
				rx := regexp.MustCompile(locale)
				return rx.MatchString(ae.Location.composeLocationString())
			})
			dae[idx].Entries = filtered
		}
	}
}

// print debug output if the flag is passed in
func debug(fmt string, a ...interface{}) {
	if debugOn {
		log.Printf(fmt, a...)
	}
}
