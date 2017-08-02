package main

import (
	"log"
	"os"
)

// Notifications that define notifications
// It's a key value pair of subject=>message
type Notifications map[string]string

// Sender is a contract for delivering notifications.
type Sender interface {
	Send(Notifications) error
}

// SMSEngine is used to send the notification as an SMS.
type SMSEngine struct {
	Recipients []string
	*log.Logger
	// Fields specific to sending a SMS.
}

// Send knows how to send a SMS notification to all recipients
func (s SMSEngine) Send(n Notifications) error {
	for _, r := range s.Recipients {
		for k := range n { // range over map with just the key (subject) returned
			s.Printf("SMS %s Sent To %s\n", k, r)
		}

	}
	return nil
}

// EmailEngine is used to send notifications as an Email.
type EmailEngine struct {
	Recipients []string
	*log.Logger
	// Fields specific to sending a Email.
}

// Send knows how to send an Email notification to all recipients
func (e EmailEngine) Send(n Notifications) error {
	for _, r := range e.Recipients {
		for k, v := range n { // range over map with  the key and value returned
			e.Printf("Email Subject:%s Sent To %s\nMessage:%s\n", k, r, v)
		}

	}
	return nil
}

// Send delivers notifications through any type that can send notifications
func Send(notifications Notifications, senders ...Sender) (int, error) {
	var count int
	for _, s := range senders {
		err := s.Send(notifications)
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func main() {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)

	var smsEngine SMSEngine //static declare

	smsEngine.Logger = logger
	smsEngine.Recipients = append(smsEngine.Recipients, "0718886543")
	smsEngine.Recipients = append(smsEngine.Recipients, "0658886543")

	var emailEngine = EmailEngine{
		Recipients: []string{"sysadmin@company.com", "itguy@company.com"},
		Logger:     logger,
	}

	notifications := make(Notifications) //translates to make(map[string]string)

	notifications["Help!"] = "Someone deleted the master transactions table!"
	notifications["Just joking"] = "haha..go back to sleep!"

	notificationsSent, err := Send(notifications, smsEngine, emailEngine)
	if err != nil {
		logger.Printf("Error while sending notifications: %v\n", err)
	}

	logger.Printf("Number of notifications sent: %d\n", notificationsSent)
}
