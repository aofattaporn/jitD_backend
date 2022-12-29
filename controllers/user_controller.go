package controllers

import (
	"context"
	"fmt"
	mapstructure "github.com/mitchellh/mapstructure"
	configs "jitD/configs"
	models "jitD/models"
	"log"

	"net/http"

	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
)

func GetAllUser(c *gin.Context) {

	users := []models.User{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("User").Documents(ctx)
	for {
		user := models.User{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "Not found")
		}
		mapstructure.Decode(doc.Data(), &user)
		users = append(users, user)
	}

	fmt.Println(users)
	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {

	// set paramether 
	id := c.Param("id")
	user := models.User{}

	// get client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// query data
	dsnap, err := client.Collection("User").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	// map and return data 
	m := dsnap.Data()
	mapstructure.Decode(dsnap.Data(), &user)
	fmt.Printf("Document data: %#v\n", m)
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {

	// create client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// seection create data
	_, _, err := client.Collection("User").Add(ctx, map[string]interface{}{
		"email":    "user123@hotmail.com",
		"password": "1111111111",
		"userName": "attaporrn1234",
	})

	// check nil value
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

func DeleteUser(c *gin.Context) {

	// users := []models.User{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("User").Documents(ctx)
	for {
		// user := models.User{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "Not found")
		}
		doc.Ref.Delete(ctx)
		// mapstructure.Decode(doc.Data(), &user)
		// users = append(users, user)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "delete data success",
	})
}
