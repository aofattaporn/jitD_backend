package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"

	mapstructure "github.com/mitchellh/mapstructure"

	"net/http"

	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
)

func GetAllPost(c *gin.Context) {

	posts := []models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("Post").Documents(ctx)
	for {
		post := models.Post{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "Not found")
		}
		mapstructure.Decode(doc.Data(), &post)
		posts = append(posts, post)
	}

	fmt.Println(posts)
	c.JSON(http.StatusOK, posts)
}

func GetPostById(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := client.Collection("Post").Add(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	_, err := client.Collection("Post").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant delete",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := client.Collection("Post").Doc(id).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}
	c.JSON(http.StatusOK, post)
}

func GetPostByUserId(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}

func GetPostByCategoryId(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}
