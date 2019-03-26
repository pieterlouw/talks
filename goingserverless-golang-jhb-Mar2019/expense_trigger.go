package goingserverless

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/functions/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type monthlyExpenseTotals struct {
	Transactions int `json:"transactions" firestore:"transactions"`
	Amount       int `json:"amount" firestore:"amount"`
}

// This function uses environment variables. see https://cloud.google.com/functions/docs/env-var

// ExpenseChangeTrigger is triggered by a change in the `expenses` collection.
func ExpenseChangeTrigger(ctx context.Context, e FirestoreEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("Old value: %+v", e.OldValue)
	log.Printf("New value: %+v", e.Value)
	log.Printf("Name: %+v", e.Value.Name)

	// Get data from Expense
	expense, err := NewFromFirestoreValue(e.Value)
	if err != nil {
		return err
	}
	log.Printf("%+v", expense)

	// TODO(pieterlouw): workout change delta if document has been updated
	//old, err := NewFromFirestoreValue(e.OldValue)
	//if err != nil {
	//		return err
	//	}
	/*delta := 0
	if expense.TotalAmount != old.TotalAmount {
		delta = expense.TotalAmount - old.TotalAmount
	}*/

	// TODO(pieterlouw): check if timestamp was updated to different month

	//get year and month from timestamp
	docName := expense.Timestamp[:7] //yyyy-mm

	//lookup monthly expense document to update totals
	totalsDoc := dbClient.Doc(fmt.Sprintf("monthly_totals/%s", docName))

	exists := true
	docsnap, err := totalsDoc.Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			exists = false
		} else {
			return fmt.Errorf("failed to get document snapshot: %v", err)
		}
	}

	var totals monthlyExpenseTotals
	if exists {
		if err := docsnap.DataTo(&totals); err != nil {
			return fmt.Errorf("failed to convert document to value: %v", err)
		}
		fmt.Println("Monthly Expenses found: ", totals)

		// update totals
		totals.Transactions++
		totals.Amount += expense.TotalAmount

		//update db
		wr, err := totalsDoc.Set(ctx, totals)
		if err != nil {
			return fmt.Errorf("failed to set totals: %v", err)
		}

		log.Println("Update onthly expense", wr)
	} else {
		wr, err := dbClient.Collection("monthly_totals").Doc(docName).Create(ctx, monthlyExpenseTotals{
			Transactions: 1,
			Amount:       expense.TotalAmount,
		})
		if err != nil {
			return fmt.Errorf("failed to create totals: %v", err)
		}

		log.Println("Added a new monthly expense", wr)
	}

	// check if alert need to be sent
	// 50% of monthly threshold, send email
	// 90% of monthly threshold, send sms
	// Monthly threshold, email address and sms number are set in env_vars

	monthlyThreshold, err := strconv.Atoi(os.Getenv("MONTHLY_THRESHOLD"))
	if err != nil {
		return err
	}

	// TODO(pieterlouw): remember if alert is already sent
	if totals.Amount > (monthlyThreshold*90)/100 {
		//send sms
		sms := os.Getenv("ALERT_SMS")

		text := fmt.Sprintf("ALERT! You have spend 90%% or more of your monthly threshold of %d.", monthlyThreshold)
		log.Printf("%s. (sending sms to %s)", text, sms)

		if err := sendAlertSMS(sms, text); err != nil {
			return fmt.Errorf("failed to send alert sms to %s: %v", sms, err)
		}

	} else if totals.Amount > (monthlyThreshold*50)/100 {
		//send email
		email := os.Getenv("WARNING_EMAIL")

		text := fmt.Sprintf("WARNING! You have spend 50%% or more of your monthly threshold of %d.", monthlyThreshold)
		log.Printf("%s. (sending email to %s)", text, email)

		if err := sendWarningEmail(email, text); err != nil {
			return fmt.Errorf("failed to send warning email to %s: %v", email, err)
		}

	}

	return nil
}

// gcloud functions deploy ExpenseChangeTrigger --runtime go111 --set-env-vars MONTHLY_THRESHOLD=100000,WARNING_EMAIL=your-email@host.emails,ALERT_SMS=+27821234567  --trigger-event providers/cloud.firestore/eventTypes/document.write --trigger-resource "projects/your-project-id-here/databases/(default)/documents/expenses/{pushId}"
