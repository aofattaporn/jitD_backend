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

	snap, err := client.Collection("User").Documents(ctx).GetAll()
	if err != nil {
		return
	}

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

// Create User
func CreateUser(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	user := models.User{}

	// Call BindJSON to bind the received JSON body
	if err := c.BindJSON(&user); err != nil {
		log.Fatalln(err)
		return
	}

	// assign rmpthy object
	user.Posts = []*firestore.DocumentRef{}
	user.Comments = []*firestore.DocumentRef{}
	user.Likes = []*firestore.DocumentRef{}

	// add data to document
	_, _, err := client.Collection("User").Add(ctx, user)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message": "create data success",
		})
	}
}

// Create User
func NamingPet(c *gin.Context) {

	id := c.Param("id")
	pet := models.User{}

	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.BindJSON(&pet); err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(pet.PetName)

	_, err := client.Collection("User").Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "PetName",
			Value: pet.PetName,
		},
	})

	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message": "naming data success",
		})
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
