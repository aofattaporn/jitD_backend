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
)

// Get all user
func GetAllUser(c *gin.Context) {

	users := []models.User{}
	user := models.User{}
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
		mapstructure.Decode(element.Data(), &user)
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// Get user by id
func GetUserById(c *gin.Context) {

	id := c.Param("id")
	user := models.User{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve user by id
	dsnap, err := client.Collection("User").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// maping data to user model
	mapstructure.Decode(dsnap.Data(), &user)
	c.JSON(http.StatusOK, user)
}

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

	// add data to document
	_, _, err := client.Collection("User").Add(ctx, user)
	errors = append(errors, err)
	if err != nil {
		for _, err := range errors {
			log.Fatalf("Failed adding alovelace: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant create data",
			})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message": "create data success",
		})
	}
}

// naming a pet
func NamingPet(c *gin.Context) {

	id := c.Param("id")
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
		c.JSON(http.StatusCreated, gin.H{
			"message": "naming data success",
		})
		return
	}
}

// Delete User
func DeleteUser(c *gin.Context) {

	id := c.Param("id")
	user := models.User{}

	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("User").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	m := dsnap.Data()
	mapstructure.Decode(dsnap.Data(), &user)
	fmt.Printf("Document data: %#v\n", m)
	c.JSON(http.StatusOK, user)
}
