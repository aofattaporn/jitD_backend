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

func GetUserID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	user := client.Collection("user").Documents(ctx)

	fmt.Println(user.GetAll())

	// return data
	c.JSON(http.StatusOK, user)
}
