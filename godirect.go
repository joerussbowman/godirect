package main

// TODO: handle http and https
// TODO: set permanent or not, maybe an option?
// TODO: app to manage config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
        "strings"
        "regexp"
)

var config Config

type Redirects map[string]string

type Host struct {
	Redirects Redirects
        Proxy   string
}

type Config struct {
        HttpPort    int
	AdminHostName string
        DefaultRedirectHostName string
	Hosts         map[string]Host
}

func readConfig() (c Config) {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
        err := json.Unmarshal(file, &c)
        if err != nil {
            fmt.Println("Error parsing config file: ",err)
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

	log.Println("Starting server")
        // TODO: https listener as well?
	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.HttpPort), nil); err != nil {
		log.Fatal(err)
	}
}
