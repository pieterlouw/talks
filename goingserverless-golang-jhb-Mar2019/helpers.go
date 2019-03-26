package goingserverless

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"log"

	"cloud.google.com/go/functions/metadata"

	"gopkg.in/sendgrid/sendgrid-go.v2"
	"bitbucket.org/ckvist/twilio/twirest"
)

// writeJSON writes a JSON Content-Type header and a JSON-encoded object to the
// http.ResponseWriter.
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	// Indent the JSON so it's easier to read.
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.WriteHeader(code)

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}

type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// getStringValue extracts a string value from a Firestore value
func (v FirestoreValue) getStringValue(name string) (string, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Error extracting value %s from %+v", name, fields)
	}
	value, ok := mapped["stringValue"].(string)
	if !ok {
		return "", fmt.Errorf("Error extracting value %s from %+v", name, fields)
	}
	return value, nil
}

// getIntegerValue extracts an integer value from a Firestore value
func (v FirestoreValue) getIntegerValue(name string) (int, error) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("Error extracting value %s from %+v", name, fields)
	}
	strValue, ok := mapped["integerValue"].(string)
	if !ok {
		return 0, fmt.Errorf("Error extracting value %s from %+v", name, fields)
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getDocumentID(ctx context.Context) (string, error) {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return "", err
	}
	parts := strings.Split(meta.Resource.Name, "/")
	if len(parts) == 0 {
		return "", errors.New("Error getting ID from context")
	}
	return parts[len(parts)-1], nil
}

func getCreateTime(ctx context.Context) (time.Time, error) {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return meta.Timestamp, nil
}

// NewFromFirestoreValue returns a new Expense, from values in FirestoreValue
func NewFromFirestoreValue(v FirestoreValue) (*Expense, error) {
	timestamp, err := v.getStringValue("timestamp")
	if err != nil {
		return nil, err
	}
	location, err := v.getStringValue("location")
	if err != nil {
		return nil, err
	}
	totalAmount, err := v.getIntegerValue("total_amount")
	if err != nil {
		return nil, err
	}
	paymentMethod, err := v.getStringValue("payment_method")
	if err != nil {
		return nil, err
	}
	comments, err := v.getStringValue("comments")
	if err != nil {
		return nil, err
	}

	// TODO(pieterlouw): add function to get slice values
	/*tags, err := v.getSliceValues("tags")
	if err != nil {
		return nil, err
	}*/
	return &Expense{
		Timestamp:     timestamp,
		Location:      location,
		TotalAmount:   totalAmount,
		PaymentMethod: paymentMethod,
		Comments:      comments,
	}, nil
}

func sendWarningEmail(email string, text string) error {

	m := sendgrid.NewMail()
	m.AddTo(email)
	m.SetSubject("50% Spending threshold reached")
	m.SetHTML(text)
	m.SetFrom("yourconscience@expenses.money")

	if err := sendgridClient.Send(m); err != nil {
		return err
	}

	return nil
}

func sendAlertSMS(number string, text string) error {
	msg := twirest.SendMessage{
		Text: text,
		From: twilioNumber,
		To:   number,
	}

	rsp, err := twilioClient.Request(msg)
	if err != nil {
		return err
	}

	log.Println("Twilio response:", rsp)

	return nil
}