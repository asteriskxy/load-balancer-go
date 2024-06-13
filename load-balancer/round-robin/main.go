package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Host struct {
	url   string
	alive bool
}

var current_host int = 0

const healthcheck_duration = 2 * time.Minute

var hosts []*Host = []*Host{
	{"http://localhost:8080", true},
	{"http://localhost:8081", true},
	{"http://localhost:8082", true},
}

func setHostStatus(host *Host, alive bool) {
	host.alive = alive
}

func getHost() (*Host, error) {
	for retires := len(hosts); retires >= 0; retires-- {
		current_host = (current_host + 1) % len(hosts)
		if hosts[current_host].alive {
			return hosts[current_host], nil
		}
	}
	return nil, fmt.Errorf("no servers found alive")
}

func lb(r *http.Request) {
	host, err := getHost()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(r.Method, host.url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Host = r.Host
	req.Header = r.Header
	req.RemoteAddr = r.RemoteAddr

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error", err)
		if host.alive {
			setHostStatus(host, false)
		}
	}

	defer res.Body.Close()

	fmt.Printf("Response from host: %s %s\n\n", res.Proto, res.Status)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request received from %s\n", r.RemoteAddr)
	lb(r)
}

func health_check() {
	for {
		fmt.Println("Health check running...")
		for index := 0; index < len(hosts); index++ {
			req, err := http.NewRequest("GET", hosts[index].url+"/health", nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				if hosts[index].alive {
					setHostStatus(hosts[index], false)
				}
			} else {
				if res.StatusCode >= 300 || res.StatusCode < 200 {
					if hosts[index].alive {
						setHostStatus(hosts[index], false)
					}
				} else {
					if !hosts[index].alive {
						setHostStatus(hosts[index], true)
					}
				}
				res.Body.Close()
			}
		}
		time.Sleep(healthcheck_duration)
	}

}

func main() {
	http.HandleFunc("/", handleRequest)
	go health_check()

	fmt.Println("Load balancer running...")

	if err := http.ListenAndServe("127.0.0.1:3030", nil); err != nil {
		log.Fatal(err)
	}
}
