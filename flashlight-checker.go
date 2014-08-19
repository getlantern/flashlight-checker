package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/getlantern/flashlight/proxy"
)

const (
	STATUS_GATEWAY_TIMEOUT = 504
)

func main() {
	http.HandleFunc("/", handle)
	port := os.Getenv("PORT")
	log.Printf("About to listen at port: %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func handle(resp http.ResponseWriter, req *http.Request) {
	flashlightClient := &proxy.Client{
		UpstreamHost:       clientIpFor(req),
		UpstreamPort:       443,
		InsecureSkipVerify: true,
	}

	flashlightClient.BuildEnproxyConfig()

	client := &http.Client{
		Transport: &http.Transport{
			Dial: flashlightClient.DialWithEnproxy,
		},
	}

	r, err := client.Head("http://www.google.com/humans.txt")
	if err != nil {
		resp.WriteHeader(STATUS_GATEWAY_TIMEOUT)
	} else {
		defer r.Body.Close()
		resp.WriteHeader(r.StatusCode)
	}
}

func clientIpFor(req *http.Request) string {
	// Client requested their info
	clientIp := req.Header.Get("X-Forwarded-For")
	if clientIp == "" {
		clientIp = strings.Split(req.RemoteAddr, ":")[0]
	}
	// clientIp may contain multiple ips, use the first
	ips := strings.Split(clientIp, ",")
	return strings.TrimSpace(ips[0])
}
