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
	user := models.User{}
	var errors []error

	// mapping data form body
	if err := c.BindJSON(&user); err != nil {
		errors = append(errors, err)
	}

	// assign empthy object
	user.Posts = []*firestore.DocumentRef{}
	user.Comments = []*firestore.DocumentRef{}
	user.Likes = []*firestore.DocumentRef{}

	// hinde
	if err := c.BindHeader(&user); err != nil {
		errors = append(errors, err)
	}

	user_id := c.Request.Header.Get("id")

	// add data to document
	fmt.Printf("header: %v\n", user_id)
	_, err := client.Collection("User").Doc(user_id).Set(ctx, user)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create data success",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message": "createdata success",
		})
	}
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
			"message": "you can acess data",
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
	user_id := c.Request.Header.Get("id")
	user := models.User{}
	user.Point = 0
	user.Posts = []*firestore.DocumentRef{}
	user.Comments = []*firestore.DocumentRef{}
	user.Likes = []*firestore.DocumentRef{}

	_, user_err := client.Collection("User").Doc(user_id).Get(ctx)
	if user_err != nil {
		if status.Code(user_err) == codes.NotFound {
			_, err := client.Collection("User").Doc(user_id).Set(ctx, user)
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

	users := []models.UserResponse{}
	userRes := models.UserResponse{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve all user
	snap, err := client.Collection("User").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// maping data to user model
	for _, element := range snap {
		/// mapdata and count list
		user := models.User{}
		mapstructure.Decode(element.Data(), &user)
		userRes.UserId = element.Ref.ID
		userRes.PetName = user.PetName
		userRes.Point = user.Point
		userRes.CountPosts = len(user.Posts)
		userRes.CountComments = len(user.Comments)
		userRes.CountLikes = len(user.Likes)

		users = append(users, userRes)
	}
	c.JSON(http.StatusOK, users)
}

// Get user by id
func GetUserById(c *gin.Context) {

	print("sdfsdfdsfdf")

	id := c.Request.Header.Get("id")
	user := models.User{}
	userRes := models.UserResponse{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve user by id
	docref := client.Collection("User").Doc(id)
	dsnap, err := docref.Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}
	// maping to response
	userRes.UserId = dsnap.Ref.ID
	userRes.PetName = user.PetName
	userRes.Point = user.Point
	userRes.CountPosts = len(user.Posts)
	userRes.CountComments = len(user.Comments)
	userRes.CountLikes = len(user.Likes)

	// maping data to user model
	mapstructure.Decode(dsnap.Data(), &userRes)

	// check quest in user
	// if have to map
	//
	// not have to initial
	//

	InitializeDailyQuests(id, docref, ctx)
	c.JSON(http.StatusOK, userRes)

}

func InitializeDailyQuests(userID string, userRef *firestore.DocumentRef, ctx context.Context) error {

	// Get the user's daily quest document
	dailyQuestRef := userRef.Collection("dailyQuest").Doc("today")
	dailyQuestDoc, err := dailyQuestRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get daily quest document: %v", err)
	}

	// If the daily quest document doesn't exist, create it
	if !dailyQuestDoc.Exists() {
		_, err := dailyQuestRef.Set(ctx, map[string]interface{}{
			"quest1": map[string]interface{}{
				"progress1": 0,
				"progress2": 0,
				"progress3": 0,
			},
			"quest2": map[string]interface{}{
				"progress1": 0,
				"progress2": 0,
				"progress3": 0,
			},
			"quest3": map[string]interface{}{
				"progress1": 0,
				"progress2": 0,
				"progress3": 0,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create daily quest document: %v", err)
		}
	}
	return nil
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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "naming data success",
		})
		return
	}
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
