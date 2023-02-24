package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"
	"time"

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
	userID := c.Request.Header.Get("id")

	// mapping data form body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get user data ",
		})
	}

	user.PetHP = 0
	// add data to document
	_, err := client.Collection("User").Doc(userID).Set(ctx, user)
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
	user_id := c.Request.Header.Get("id")
	user := models.User{}
	user.Point = 0

	_, user_err := client.Collection("User").Doc(user_id).Get(ctx)
	if user_err != nil {
		if status.Code(user_err) == codes.NotFound {
			user.PetHP = 0
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
	user := models.User{}
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

	// Decode user's daily quests from the document
	if dailyQuests, err := dsnap.DataAt("DailyQuests"); err == nil {
		if err := mapstructure.Decode(dailyQuests, &user.DailyQuests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to decode daily quests",
			})
			return
		}
	}

	// Check and update user's daily quests
	if updatedUser, err := checkQuest(user, client, ctx, userID); err == nil {
		user = updatedUser
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update daily quests",
		})
		return
	}

	// Return user response
	c.JSON(http.StatusOK, userRes)
}

func checkQuest(user models.User, client *firestore.Client, ctx context.Context, userID string) (models.User, error) {

	today := time.Now().Day()

	if user.DailyQuests == nil || user.DailyQuests.QuestDate.Day() != today {
		var questName = [3]string{"PostQuest", "CommentQuest", "LikeQuest"}
		var element []*models.Quest
		for i, _ := range questName {
			element = append(element, &models.Quest{
				QuestName:      questName[i],
				Progress:       0,
				MaxProgress:    3,
				Reward:         5,
				IsGetPoint:     false,
				Completed:      false,
				LastCompletion: time.Now().UTC(),
			})
		}

		user.DailyQuests = &models.DailyQuestProgress{
			QuestDate: time.Now().UTC(),
			Quests:    element,
		}

		// set data
		user.DailyQuests.QuestDate = time.Now().UTC()

		// save to db
		_, err := client.Collection("User").Doc(userID).Set(ctx, user)
		if err != nil {
			return user, err
		}

		return user, nil

	} else {
		return user, nil
	}

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
