# Go-ing Serverless

This is a repo with slides and demo code for a talk I did at the [Joburg Golang Meetup](https://www.meetup.com/Joburg-Golang-Group/events/259877159/) in March 2019.

The demo's are deployed to [Google Cloud Functions](https://cloud.google.com/functions) so to run the code you would need a Google Cloud Project, instructions to setup a project for Google Cloud Functions can be found [here](https://cloud.google.com/functions/docs/quickstart).

The first demo is a function that randomize a slice of predefined Game of Thrones quotes.

The second demo has more meat. It's an expense tracker where you can log expenses (which is saved in Firestore) and then there are a trigger function that listens to changes on the Firestore document to track monthly expenses and warn (by email using Sendgrid) or alert (by SMS using Twilio) if 50% or 90% of the specified threshold has been reached.

## TODO

- [ ] Create a HTML/JS frontend.
- [ ] Accept gps location of location and run through a Maps/Places API.
- [ ] Upload receipt(s) and let a trigger function run OCR using GCP Cloud Vision API.

## Resources

https://medium.com/@feloy/firebase-saving-scores-into-firestore-with-go-functions-b128fd8c425
https://medium.com/google-cloud/google-cloud-functions-for-go-57e4af9b10da
https://medium.com/@hiranya911/firebase-database-interactions-from-go-121d19217639
https://thenewstack.io/go-the-programming-language-of-the-cloud/
https://medium.com/google-cloud/firebase-developing-serverless-functions-in-go-963cb011265d
https://www.martinfowler.com/articles/serverless.html

https://medium.com/google-cloud/using-google-cloud-vision-api-with-golang-830e70323de7