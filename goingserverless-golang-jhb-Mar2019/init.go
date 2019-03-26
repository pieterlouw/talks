package goingserverless

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"gopkg.in/sendgrid/sendgrid-go.v2"
	"bitbucket.org/ckvist/twilio/twirest"
)

// dbClient is a Firestore client, reused between function invocations.
var dbClient *firestore.Client

var sendgridClient *sendgrid.SGClient
var twilioClient *twirest.TwilioClient

var twilioNumber string 

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// NOTE(pieterlouw): Rather do lazy-loading if not all code paths use all global objects.
// https://cloud.google.com/functions/docs/bestpractices/tips#do_lazy_initialization_of_global_variables

func init() {
	// initialize firestore
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	dbClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
	
	// NOTE(pieterlouw): Although the API keys below are hard coded, ideally it should be read from a secret store manager like Vault.
	sendgridClient = sendgrid.NewSendGridClientWithApiKey("use-your-sendgrid-api-key-here") 


	twilioNumber =  "twilio-number-her"
	twilioClient = twirest.NewClient("twilio-account-sid-here", "twilio-auth-token-here")
}

