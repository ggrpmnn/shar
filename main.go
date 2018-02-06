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
)

func init() {
	flag.BoolVar(&debugOn, "d", false, "enable debug output")
	flag.BoolVar(&jsonOut, "j", false, "output results in JSON format")
	flag.StringVar(&filename, "f", "/var/log/auth.log", "auth log file to parse")
	flag.IntVar(&threshold, "n", 0, "minimum number of attempts required to collect info")
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

	if jsonOut {
		debug("outputting JSON")
		attempts.jsonPrint()
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
