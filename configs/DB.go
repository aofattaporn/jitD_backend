package configs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func ConnectFireStore(ctx context.Context) firebase.App {

	opt := option.WithCredentialsFile("./jitd-application-firebase-adminsdk-ee9le-ebfdb5c4c5.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return *app
	} else {
		fmt.Print("Connection success", app)
		return *app
	}
}

func CreateClient(ctx context.Context) *firestore.Client {

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

func CreateCheckAuth(ctx context.Context, c *gin.Context, x string) {

	opt := option.WithCredentialsFile("./jitd-application-firebase-adminsdk-ee9le-ebfdb5c4c5.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	client, err := app.Auth(ctx)
	token, err := client.VerifyIDToken(ctx, x)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "cant access data",
		})
		c.Abort()
	} else {
		fmt.Printf("token.UID: %v\n", token.UID)
		c.Request.Header.Add("id", token.UID)
		c.Next()
	}
}

func Verify(c *gin.Context) {

	type testHeader struct {
		Bareers string `header:"Bareers"`
	}
	h := testHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
	}

	splitToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")
	reqToken := splitToken[1]

	fmt.Printf("reqToken: %v\n", reqToken)

	CreateCheckAuth(c, c, reqToken)
}
