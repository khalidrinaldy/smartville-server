package helper

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/messaging"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func SendNotification(registrationToken string, title, notificationMsg string) (string,error) {
	//Initialize App
	opt := option.WithCredentialsFile("firebase-key.json")
	config := &firebase.Config{ProjectID: "smartville-fcm"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return "error initializing app", err
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return "error getting Messaging client", err
	}

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	// message := &messaging.Message{
	// 	Data: map[string]string{
	// 		"message": notificationMsg,
	// 	},
	// 	Token: registrationToken,
	// }

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body: notificationMsg,
		},
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		return "error while sending notification", err
	}

	// Response is a message ID string.
	log.Printf("Successfully sent notification: %s", response)
	return fmt.Sprintf("Successfully sent notification: %s", response), nil
}
