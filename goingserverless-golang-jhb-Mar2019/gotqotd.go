package goingserverless

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var gotQuotes = []string{"The North Remembers..",
	"A Lannister always pays his debts.",
	"Winter is coming..",
	"Hold the door!",
	"Valar Dohaeris",
	"Bend the knee!",
	"Valar Morghulis",
	"What is dead may never die.",
	"Fear cuts deeper than swords."}

// GotQOTD is an HTTP Cloud Function.
func GotQOTD(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<h1>%s</h1><h4><i>Refresh to view another quote</i></h4>`, gotQuotes[rand.Intn(len(gotQuotes))])
}

// gcloud functions deploy GotQOTD --runtime=go111 --trigger-http --entry-point=GotQOTD
