package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/getlantern/enproxy"
	"github.com/getlantern/flashlight/proxy"
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

	log.Printf("UpstreamHost: %s", flashlightClient.UpstreamHost)
	enproxyConfig := flashlightClient.BuildEnproxyConfig()

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network string, addr string) (net.Conn, error) {
				conn := &enproxy.Conn{
					Addr:   addr,
					Config: enproxyConfig,
				}
				err := conn.Connect()
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	}

	r, err := client.Head("http://www.google.com/humans.txt")
	if err != nil {
		resp.WriteHeader(500)
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
