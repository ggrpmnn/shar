package main

import (
	"bufio"
	"log"
	"os"
)

// AuthEntry keeps track of the login attempts on a per-IP basis
type AuthEntry struct {
	count int
	users []string
}

type entries map[string]AuthEntry

func main() {
	file, err := os.Open("/var/log/auth.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.Println("auth file loaded")

	attempts := make(entries, 0)
	// example auth log line for invalid entries: "Feb  1 19:02:48 grpi sshd[8749]: Invalid user pi from 202.120.42.141"
	rx := regexp.MustCompile(`(\w{3}\s+\d+\s+\d{2}:\d{2}\d{2})\s+grpi sshd[\d+]: Invalid user (.*) from (\d{1,3}\.\d{1,3}\.d{1,3})`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := rx.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		// matches[0]=full string, [1]=date, [2]=user, [3]=IP
		if attempts.exists(matches[3]) {
			attempts[matches[3]].addUser(matches[2])
		} else {
			new := AuthEntry{count: 1, users: []string{matches[2]}}
			attempts[matches[3]] = new
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("finished parsing log file")

	attempts.print()
}

// returns true if the IP string exists in the given map
func (e entries) exists(ip string) bool {
	_, ok := e[ip]
	return ok
}

func (e entries) print() {
	for ip, ae := range entries {
		fmt.Printf("IP: %s, Attempt Count: %d, Users: %s\n", ip, ae.count, strings.Join(ae.users, ", "))
	}
}

// adds the username to the list for the given IP AuthEntry struct
func (ae *AuthEntry) addUser(user string) {
	ae.count++
	for _, un := range ae.users {
		if un == user {
			return
		}
	}
	ae.users = append(ae.users, user)
}
