package main
// just a quick web server to show
// the host and path of a request. Used
// to verify the proxy and redirects are
// working correctly

import(
    "io"
    "fmt"
    "log"
    "net/http"
)


func handler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, fmt.Sprintf("host: %v\npath: %v", r.Host, r.URL))
}
func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Server started")
}
