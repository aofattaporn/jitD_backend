package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create a User
func CreateUser(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	userID := c.Request.Header.Get("id")

	// add data to document
	_, err := client.Collection("User").Doc(userID).Set(ctx, models.User{
		PetName:       "",
		PetHP:         0,
		Point:         0,
		HistorySearch: []string{},
		BookMark:      []*firestore.DocumentRef{},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to set user to DB ",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "createdata success",
	})
}

// Create a User
func SignIn(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	user_id := c.Request.Header.Get("id")

	// add data to document
	fmt.Printf("header: %v\n", user_id)
	_, err := client.Collection("User").Doc(user_id).Get(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "you can not acess data",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "you can acess data",
		})
	}
}

// Create a User
func SignInGoogle(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// delcare id for use in function
	user_id := c.Request.Header.Get("id")

	// get a user
	_, user_err := client.Collection("User").Doc(user_id).Get(ctx)
	if user_err != nil {
		if status.Code(user_err) == codes.NotFound {
			_, err := client.Collection("User").Doc(user_id).Set(ctx, models.User{
				PetName:       "",
				PetHP:         0,
				Point:         0,
				HistorySearch: []string{},
				BookMark:      []*firestore.DocumentRef{},
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Cant to set data",
				})
				return
			} else {
				c.JSON(http.StatusCreated, gin.H{
					"message": "something wronng",
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

// Get all user
func GetAllUser(c *gin.Context) {

	users := []models.User{}
	// dialyQuest := models.DailyQuestProgress{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve all user
	userSnap, err := client.Collection("User").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// maping data to user model
	for _, element := range userSnap {
		user := models.User{}
		mapstructure.Decode(element.Data(), &user)
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// Get user by id
func GetUserById(c *gin.Context) {
	// Get user ID from the header
	userID := c.Request.Header.Get("id")

	// Initialize user and user response
	userRes := models.UserResponse{}

	// Create Firestore client and context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get user document by ID
	docRef := client.Collection("User").Doc(userID)
	dsnap, err := docRef.Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve user information",
		})
		return
	}

	// Decode user document to user model
	if err := dsnap.DataTo(&userRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to decode user data",
		})
		return
	}
	userRes.UserID = userID

	// Return user response
	c.JSON(http.StatusOK, userRes)
}

// naming a pet
func NamingPet(c *gin.Context) {

	id := c.Request.Header.Get("id")
	pet := models.User{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// mapping data form body
	if err := c.BindJSON(&pet); err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get body from you ",
		})
		return
	}

	// set data to DB
	_, err := client.Collection("User").Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "PetName",
			Value: pet.PetName,
		},
	})

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant complete upddate data ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"petName": pet.PetName,
	})
	return
}

// Delete User
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

// update user
