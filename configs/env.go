package configs

import (
	"context"
	"log"
	"fmt"


	// "google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

)

func CreateClient(ctx context.Context) *firestore.Client {

	opt := option.WithCredentialsFile("./jitd-application-firebase-adminsdk-ee9le-ebfdb5c4c5.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}


	// Sets your Google Cloud Platform project ID.
	// projectID := "jitd-application"

	client, err := app.Firestore(ctx) 
	if err != nil {
			  log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}