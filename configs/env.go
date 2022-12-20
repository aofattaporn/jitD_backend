package configs

import (
	"context"
	"log"

	// "google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "jitd-application"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
			  log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}