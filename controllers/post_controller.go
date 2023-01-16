package controllers

import (
	"context"
	configs "jitD/configs"
	models "jitD/models"
	"net/http"

	// "strings"
	"time"

	"github.com/gin-gonic/gin"
)

// service create post
func CreatePost(c *gin.Context) {
	// id := c.Param("id")
	// user := models.User{}
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// service push post id

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

	_, _, err := client.Collection("Post").Add(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	c.JSON(http.StatusOK, post)
}
