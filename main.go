package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// "github.com/googleapis/enterprise-certificate-proxy/client"

	"context"
   "google.golang.org/api/iterator"
	"log"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
)

func createClient(ctx context.Context) *firestore.Client {

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

func main() {

   ctx := context.Background()
   client := createClient(ctx)

   iter := client.Collection("Book").Documents(ctx)
   for {
        doc, err := iter.Next()
        if err == iterator.Done {
                break
        }
        if err != nil {
                log.Fatalf("Failed to iterate: %v", err)
        }
        fmt.Println(doc.Data())
   }


	docsnap := client.Doc("Book/r29kO5eHqZhSq9NDipCN")
	y, err := docsnap.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dataMap := y.Data()

	// initail router
	router := gin.Default()
   fmt.Printf("Success fully")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, dataMap)
	})

	
	router.Run("localhost:3000")
}

