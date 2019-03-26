package goingserverless

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Expense is a domain model
type Expense struct {
	Timestamp     string   `json:"timestamp" firestore:"timestamp"`
	TotalAmount   int      `json:"total_amount" firestore:"total_amount"`
	Location      string   `json:"location" firestore:"location"`
	Tags          []string `json:"tags" firestore:"tags"`
	PaymentMethod string   `json:"payment_method" firestore:"payment_method"`
	Comments      string   `json:"comments" firestore:"comments"`
}

const maxBodySize = int64(1048576) // 1mb

// AddExpense creates a new expense
func AddExpense(w http.ResponseWriter, r *http.Request) {

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

	//add to database

	writeJSON(w, http.StatusAccepted, fmt.Sprintf("expense has been logged successfully"))
}
