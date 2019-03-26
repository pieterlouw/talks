package goingserverless

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// AddExpenseDB creates a new expense and adds it to the database.
func AddExpenseDB(w http.ResponseWriter, r *http.Request) {

	var expenseReceived Expense

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, fmt.Sprintf("invalid HTTP method: %s", r.Method))
		return
	}

	// process HTTP body and decode into internal Expense type.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %v", err))
		return
	}
	if err := r.Body.Close(); err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to close request body: %v", err))
		return
	}
	if err := json.Unmarshal(body, &expenseReceived); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, fmt.Sprintf("failed to decode request body: %v", err))
		return
	}

	ctx := r.Context()

	_, wr, err := dbClient.Collection("expenses").Add(ctx, expenseReceived)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to save expense: %v", err))
		return
	}

	log.Println(wr)

	writeJSON(w, http.StatusAccepted, fmt.Sprintf("expense has been logged successfully"))
}

// gcloud functions deploy AddExpense --runtime=go111 --trigger-http --entry-point=AddExpenseDB
