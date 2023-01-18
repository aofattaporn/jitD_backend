package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"net/http"

	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
)

// service create post
func CreatePost(c *gin.Context) {
	// iniitail variable
	id := c.Param("id")
	user := models.User{}
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// service create post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Comment = []string{}
	post.Like = []string{}
	post.Date = time.Now()

	//  Adding document post to collection
	iddd, _, err3 := client.Collection("Post").Add(ctx, post)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	// find user colection and set data to DB
	dsnap, err := client.Collection("User").Where("UserID", "==", id).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}

	mapstructure.Decode(dsnap[0].Data(), &user)
	user.Posts = append(user.Posts, iddd)
	fmt.Printf("user.Posts: %v\n", iddd.Parent)

	_, err2 := client.Collection("User").Doc(id).Set(ctx, user)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}
	c.JSON(http.StatusOK, post)
}

func GetAllPost(c *gin.Context) {
	posts := []models.Post{}
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	snap, err := client.Collection("User").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, element := range snap {
		mapstructure.Decode(element.Data(), &post)
		posts = append(posts, post)
	}
	c.JSON(http.StatusOK, posts)
}
