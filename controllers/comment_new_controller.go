package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// C (create comment)
func NewCreateComment(c *gin.Context) {

	// declare variable
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	post := models.Post2{}
	comment := models.Comment2{}
	// commentRes := models.CommentResponse{}

	// get post data
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	// get data form body
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
		})
		return
	}

	comment.UserId = client.Collection("User").Doc(userID).ID
	comment.Like = []*models.LikeComment{}
	comment.Date = time.Now().UTC()
	comment.CommentID = uuid.NewString()
	post.Comment = append(post.Comment, &comment)

	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to add a comment",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// R (read a comment)
func NewGetComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	post := models.Post2{}
	// comments := []models.Comment2{}/
	commentRes := models.CommentResponse{}

	// get Post by id
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	if len(post.Comment) < 1 {
		c.JSON(http.StatusOK, post.Comment)
		return
	}

	for _, element := range post.Comment {
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = client.Collection("User").Doc(userID).ID
		commentRes.PostId = postDoc.Ref.ID
		commentRes.CountLike = len(element.Like)
	}
	c.JSON(http.StatusOK, commentRes)
}

// R (read a comment)
func NewGetAllComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	post := models.Post2{}
	commentsRes := []models.CommentResponse{}
	commentRes := models.CommentResponse{}

	// get Post by id
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	for _, element := range post.Comment {
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = client.Collection("User").Doc(userID).ID
		commentRes.PostId = postDoc.Ref.ID
		commentRes.CountLike = len(element.Like)
		commentRes.Date = element.Date
		commentsRes = append(commentsRes, commentRes)
	}
	c.JSON(http.StatusOK, commentsRes)
}

func NewUpdateComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	commentID := c.Param("comment_id")
	post := models.Post2{}
	comment := models.Comment2{}

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
		})
		return
	}

	// get Post by id
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	// find index
	for _, element := range post.Comment {

		if element.CommentID == commentID {
			element.Content = comment.Content
			print(element.Content)
			print("------------")
			break
		}
	}

	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}

	c.JSON(http.StatusOK, post.Comment[0])

}
