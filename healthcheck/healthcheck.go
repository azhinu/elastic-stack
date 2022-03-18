package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/tidwall/gjson"
)

// New help text
const (
	usage = `Elastic stack healthchecker
  Make GET request to host.

Usage: healthcheck [options] [service] [Host]

  Services: elastic | kibana | logstash

  HTTP Host: [Default: "https://localhost:9200"]
    Non standart HTTP/HTTPS ports should be set explicitly.

  Options:
`
)

var (
	argURL, argService string
	argUser            = flag.String("u", "remote_monitoring_user", "Basic Auth `username`")
	argPasswd          = flag.String("p", "", "Basic Auth `password`")
	argStatus          = flag.String("s", "green|available", "Valid `status` (green, yellow, red), accept RegExp. \n Kibana send \"available\" status")
)

func init() {
	// Disable TLS verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Override log Prefix
	log.SetFlags(0)
	// Override default help text
	flag.Usage = func() {
		log.Print(usage)
		flag.PrintDefaults()
	}

	// Parse Args
	flag.Parse()

	// Get and check Firts arg (command name)
	argService = flag.Arg(0)
	matched, err := regexp.MatchString(`[Ee]lastic|[Kk]ibana|[Ll]ogstash`, argService)
	if err != nil || !matched {
		log.Fatal("Incorrest service name. Use: elastic | kibana | logstash")
	}

	// Check if second arg (url) exist
	argURL = flag.Arg(1)
	if argURL != "" {
		// Check if URL valid
		matched, err := regexp.MatchString(`^http(s)?://([\w\d\.-])+(:\d+)?$`, argURL)
		if err != nil || !matched {
			log.Fatal("Incorrect HTTP endpoint address. Use healthcheck http://localhost:9002")
		}
	} else {
		// Set default url
		switch argService {
		case "elastic", "Elastic":
			argURL = "https://localhost:9200"
		case "kibana", "Kibana":
			argURL = "http://localhost:5601"
		case "logstash", "Logstash":
			argURL = "http://localhost:9600"
		}
	}

	// Check unexpected Flags and Args
	if flag.Arg(2) != "" {
		log.Fatal("Got unexpected argument: ", flag.Arg(2), "\n", "Usage: healthcheck [options] [service] [HTTP Endpoint]")
	}
}

func main() {
	// Select command
	switch argService {
	case "elastic", "Elastic":
		checkStatus(req(argURL + "/_cat/health?h=status"))
	case "kibana", "Kibana":
		checkStatus(gjson.Get(req(argURL+"/api/status"), "status.overall.level").String())
	case "logstash", "Logstash":
		checkStatus(gjson.Get(req(argURL), "status").String())
	}
}

func req(url string) string {
	// Make HTTP client
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(*argUser, *argPasswd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func checkStatus(status string) {
	matched, err := regexp.MatchString(*argStatus, status)
	if err != nil || !matched {
		log.Fatal("Service status:", status)
	}
	log.Println("OK")
}
