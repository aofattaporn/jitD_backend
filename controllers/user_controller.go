package controllers

import (
	"context"
	configs "jitD/configs"
	models "jitD/models"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// * Create a User
func CreateUser(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	userID := c.Request.Header.Get("id")

	// create user document
	user := models.User{
		UserID:       userID,
		PetName:      "my bear",
		PetHP:        100,
		Point:        0,
		IsAdmin:      false,
		FCMToken:     "",
		BookMark:     []*firestore.DocumentRef{},
		Notification: []*models.Notification{},
		RegisterDate: time.Now().UTC(),
	}

	// create a user infomation
	if _, err := client.Collection("User").Doc(userID).Set(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create user document",
		})
		return
	}

	// TODO: create a test infomation
	_, _, err := client.Collection("ResultStress").Add(ctx, models.TestResualt{
		UserID:   client.Collection("User").Doc(userID),
		TestDate: time.Now().UTC(),
		TestName: "Test Stress",
		Point:    0,
		Result:   "No data",
		Desc:     "please to do this test",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create user document",
		})
		return
	}

	// TODO: create a quest infomation
	_, _, err = client.Collection("ResultConsult").Add(ctx, models.TestResualt{
		UserID:   client.Collection("User").Doc(userID),
		TestDate: time.Now().UTC(),
		TestName: "Test Consult",
		Point:    0,
		Result:   "No data",
		Desc:     "please to do this test",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to create user document",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "create a user succsess fully",
	})
}

// *signin a user
func SignIn(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()
	userID := c.Request.Header.Get("id")

	// add data to document
	_, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "you can not acess data",
		})
	}

	// response status to clieint
	c.JSON(http.StatusOK, gin.H{
		"message": "signin user success",
	})
}

// *signin a user for check admin role
func SignInAdmin(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	userID := c.Request.Header.Get("id")

	// add data to document
	_, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "you can not acess data",
		})
	}

	// response status to clieint
	c.JSON(http.StatusOK, gin.H{
		"message": "signin user success",
	})
}

// *signin a user with google
func SignInGoogle(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// delcare id for use in function
	userID := c.Request.Header.Get("id")

	// get a user
	_, user_err := client.Collection("User").Doc(userID).Get(ctx)
	if user_err != nil {
		if status.Code(user_err) == codes.NotFound {
			_, err := client.Collection("User").Doc(userID).Set(ctx, models.User{
				UserID:       userID,
				PetName:      "my bear",
				PetHP:        100,
				Point:        0,
				IsAdmin:      false,
				FCMToken:     "",
				BookMark:     []*firestore.DocumentRef{},
				Notification: []*models.Notification{},
				RegisterDate: time.Now().UTC(),
			})

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Cant to set data",
				})
				return
			} else {
				// TODO: create a test infomation
				_, _, err := client.Collection("ResultStress").Add(ctx, models.TestResualt{
					UserID:   client.Collection("User").Doc(userID),
					TestDate: time.Now().UTC(),
					TestName: "Test Stress",
					Point:    0,
					Result:   "No data",
					Desc:     "please to do this test",
				})

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "unable to create user document",
					})
					return
				}

				// TODO: create a quest infomation
				_, _, err = client.Collection("ResultConsult").Add(ctx, models.TestResualt{
					UserID:   client.Collection("User").Doc(userID),
					TestDate: time.Now().UTC(),
					TestName: "Test Consult",
					Point:    0,
					Result:   "No data",
					Desc:     "please to do this test",
				})
				c.JSON(http.StatusCreated, gin.H{
					"message": "create succcess",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something wronng",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Data found",
		})
		return
	}
}

// *Get all user
func GetAllUser(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()

	// retrieve all users
	userDocs, err := client.Collection("User").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to get user documents",
		})
		return
	}

	// extract user data from documents
	users := make([]models.User, len(userDocs))
	for i, doc := range userDocs {
		if err := doc.DataTo(&users[i]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "unable to extract user data",
			})
			return
		}
	}

	// respond to client
	c.JSON(http.StatusOK, users)
}

// *Get user by id
func GetUserById(c *gin.Context) {
	// Get user ID from the header
	userID := c.Request.Header.Get("id")

	// Create Firestore client and context
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()

	// get user document by ID
	docRef := client.Collection("User").Doc(userID)
	userSnap, err := docRef.Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unable to retrieve user information",
		})
		return
	}

	// extract user data from document
	var user models.User
	if err := userSnap.DataTo(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to extract user data",
		})
		return
	}

	// construct user response
	userRes := models.UserResponse{
		UserID:   user.UserID,
		PetName:  user.PetName,
		PetHP:    user.PetHP,
		Point:    user.Point,
		BookMark: make([]string, len(user.BookMark)),
	}
	for i, ref := range user.BookMark {
		userRes.BookMark[i] = ref.ID
	}

	// respond to client
	c.JSON(http.StatusOK, userRes)
}

// *naming a pet
func NamingPet(c *gin.Context) {

	// Get user ID from the header
	userID := c.Request.Header.Get("id")

	// Parse the pet name from the request body
	var requestBody struct {
		PetName string `json:"petName"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	// Update the user's pet name in the database
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	docRef := client.Collection("User").Doc(userID)
	if _, err := docRef.Update(ctx, []firestore.Update{
		{
			Path:  "PetName",
			Value: requestBody.PetName,
		},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update pet name",
		})
		return
	}

	// Return the new pet name
	c.JSON(http.StatusOK, gin.H{
		"petName": requestBody.PetName,
	})
}

// *Delete User
func DeleteUser(c *gin.Context) {

	id := c.Request.Header.Get("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("User").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	c.JSON(http.StatusOK, dsnap.UpdateTime)
}
