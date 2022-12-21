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

func ConnectFireStore(ctx context.Context) firebase.App{

	opt := option.WithCredentialsFile("./jitd-application-firebase-adminsdk-ee9le-ebfdb5c4c5.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return *app
	} else{
		fmt.Print("Connection success", app)
		return *app
	}
}

func CreateClient(ctx context.Context) *firestore.Client{

	opt := option.WithCredentialsFile("./jitd-application-firebase-adminsdk-ee9le-ebfdb5c4c5.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	} 

	client, err := app.Firestore(ctx) 
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}