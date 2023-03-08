package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pin comment
func PinComment(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	postID := c.Param("post_id")
	commentID := c.Param("comment_id")

	// get post from post collection
	post := models.Post{}
	postData, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
			"error":   err.Error(),
		})
		return
	}

	// mapping fata to post struc model
	if err := postData.DataTo(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to decode user data",
			"error":   err.Error(),
		})
		return
	}

	// find that's want to set index
	for index, _ := range post.Comment {
		if post.Comment[index].CommentID == commentID {
			post.Comment[index].IsPin = true
			break
		}
	}

	// set data to post collection
	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
			"error":   err.Error(),
		})
		return
	}

	// Return success response with bookmarked posts
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved bookmarked posts",
	})
}

// pin comment
func UnPinComment(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	postID := c.Param("post_id")
	commentID := c.Param("comment_id")

	// get post from post collection
	post := models.Post{}
	postData, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
			"error":   err.Error(),
		})
		return
	}

	// mapping fata to post struc model
	if err := postData.DataTo(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to decode user data",
			"error":   err.Error(),
		})
		return
	}

	// find that's want to set index
	for index, _ := range post.Comment {
		if post.Comment[index].CommentID == commentID {
			post.Comment[index].IsPin = false
			break
		}
	}

	// set data to post collection
	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
			"error":   err.Error(),
		})
		return
	}

	// Return success response with bookmarked posts
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved bookmarked posts",
	})
}
