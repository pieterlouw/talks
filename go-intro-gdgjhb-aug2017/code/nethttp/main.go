package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// HelloHandler is the handler function for the /hello route
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Requested URL: %s", r.URL.Path) // logging
	fmt.Fprintf(w, "Hello GDG!")                   // Fprintf expects a http.ResponseWriter type
}

func main() {
	http.HandleFunc("/hello", HelloHandler) // setup routing
	http.ListenAndServe(":8080", nil)       // start serving on port 8080, this call is blocking
}
