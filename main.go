package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	pathHostTarget = map[string]string{
		// change required_path into a path that client need to access when forwading to targeted host
		// change target_schema_host to a target host with the schema
		// example:
		// {
		//  "article": "https://article.com"
		//}
		"required_path": "target_schema_host",
	}

	// change this for default forward request if no path and host target configured
	defaultPath = "article"
	defaultHost = ""
)

func overrideRequest(r *http.Request) {
	pathSplitted := strings.Split(r.URL.Path, "/")

	firstPath := defaultPath
	if len(pathSplitted) > 1 && pathSplitted[1] != "" {
		firstPath = pathSplitted[1]
	}

	host := defaultHost
	if pathHostTarget[firstPath] != "" {
		host = pathHostTarget[firstPath]
	}
	target, _ := url.Parse(host)

	r.Host = target.Host
	r.URL.Host = r.Host
	r.URL.Scheme = target.Scheme
}

func main() {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost",
	})
	proxy.Director = overrideRequest
	http.Handle("/", proxy)
	if err := http.ListenAndServe(":7723", nil); err != nil {
		log.Fatal(err)
	}
}