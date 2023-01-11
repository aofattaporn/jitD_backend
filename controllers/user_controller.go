package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"

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

	// section create data
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&user); err != nil {
		log.Fatalln(err)
		return
	}

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
