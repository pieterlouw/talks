package main

import (
	"log"
	"net/http"

	fn "github.com/pieterlouw/talks/goingserverless"
)

// This is an example of how your functions can be tested locally using a http.Listener

func main() {
	http.HandleFunc("/", fn.GotQOTD)
	log.Println("GoT QOTD server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
