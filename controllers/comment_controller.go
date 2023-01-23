package controllers

import (
	"context"
	"fmt"
	"jitD/models"
	"net/http"

	"github.com/gin-gonic/gin"

	configs "jitD/configs"
	"log"

	mapstructure "github.com/mitchellh/mapstructure"

	"google.golang.org/api/iterator"
)

func GetAllComment(c *gin.Context) {

	comments := []models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("Comment").Documents(ctx)
	for {
		comment := models.Comment{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "Not found")
		}
		mapstructure.Decode(doc.Data(), &comment)
		comments = append(comments, comment)
	}

	fmt.Println(comments)
	c.JSON(http.StatusOK, comments)
}

func GetCommentById(c *gin.Context) {
	id := c.Param("id")
	comment := models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func CreateComment(c *gin.Context) {
	comment := models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := client.Collection("Comment").Add(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	c.JSON(http.StatusOK, comment)
}

func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	comment := models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := client.Collection("Comment").Doc(id).Set(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	_, err := client.Collection("Comment").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant delete",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
