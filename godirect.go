package main

// TODO: handle http and https
// TODO: set permanent or not, maybe an option?
// TODO: app to manage config

import (
	"os"
	"fmt"
	"log"
	"regexp"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/http/httputil"
)

var config Config

type Redirects map[string]string

type Host struct {
	Redirects Redirects
	Proxy     string
}

type Config struct {
	CommandPort             int
	HttpPort                int
        HttpsPort               int
        HttpsCertFile           string
        HttpsKeyFile            string
	DefaultRedirectHostName string
	Hosts                   map[string]Host
}

func readConfig() (c Config) {
    // reads config from a file named config.json in the same path as where
    // the executable is called.
    // TODO: make location of config and file name an optional parameter
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	err := json.Unmarshal(file, &c)
	if err != nil {
		fmt.Println("Error parsing config file: ", err)
		os.Exit(1)
	}
	return

}

func handler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	// check the host name, if we don't have an entry in the
	// config for this host, use the global redirect
	if h, ok := config.Hosts[host]; ok {
		// if there is a redirect for this path, do it
		redirect := fmt.Sprintf("%s", r.URL)
		if redir, ok := h.Redirects[redirect]; ok {
			// this just checks to see if the configured redirect is
			// absolute or not.
			log.Println(redir)
			external, err := regexp.MatchString("^.*://*", redir)
			if err != nil {
				log.Fatal(err)
			}
			if external {
				log.Println("doing external redirect for", redir)
				http.Redirect(w, r, redir, http.StatusFound)
			} else {
				log.Println("doing internal redirect for", redir)
				http.Redirect(w, r, fmt.Sprintf("http://%s%s", r.Host, redir), http.StatusFound)
			}
		} else {
			u, err := url.Parse(fmt.Sprintf("http://%s", h.Proxy))
			if err != nil {
				log.Fatal(err)
			}
			h := httputil.NewSingleHostReverseProxy(u)
			h.ServeHTTP(w, r)
		}
	} else {
		// this is the global redirect
		http.Redirect(w, r, fmt.Sprintf("http://%v", config.DefaultRedirectHostName), http.StatusFound)
	}

}

func main() {
	config = readConfig()
	http.HandleFunc("/", handler)

	// TODO: https listener as well?
	go func() {
            log.Println("Starting server")
            if err := http.ListenAndServe(fmt.Sprintf(":%v", config.HttpPort), nil); err != nil {
		log.Fatal(err)
            }
        }()
        if config.HttpsPort > 0 {
            go func(){
                log.Println("Starting TLS server")
                if err := http.ListenAndServeTLS(fmt.Sprintf(":%v", config.HttpsPort), config.HttpsCertFile, config.HttpsKeyFile, nil); err != nil {
                    log.Fatal(err)
                }
            }()
        }
        select {}
}
