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

	post.Content = ""
	post.Comment = []string{}
	post.Like = []string{}
	post.Date = time.Now()
	post.IsPublic = true

	iddd, _, err3 := client.Collection("Post").Add(ctx, post)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	// service push post id find nd update
	dsnap, err := client.Collection("User").Where("UserID", "==", id).Documents(ctx).GetAll()
	if err != nil {
		print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}
	mapstructure.Decode(dsnap[0].Data(), &user)
	fmt.Printf("user.Posts: %v\n", dsnap[0].Data())

	user.Posts = append(user.Posts, iddd.ID)

	_, err2 := client.Collection("User").Doc(id).Set(ctx, user)
	if err2 != nil {
		print(err)

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
