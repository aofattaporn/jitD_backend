package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// C (create comment)
func NewCreateComment(c *gin.Context) {
	// Get context and client
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()

	// Get post and user data
	postID := c.Param("post_id")
	post, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	userID := c.Request.Header.Get("id")

	// Parse request body
	var comment models.Comment2
	if err := c.BindJSON(&comment); err != nil {
		errorCheck(c, err, "cant to get data from body")
		return
	}

	// Add comment to post
	comment.UserId = client.Collection("User").Doc(userID).ID
	comment.Like = []*models.LikeComment{}
	comment.Date = time.Now().UTC()
	comment.CommentID = uuid.NewString()

	post.Comment = append(post.Comment, &comment)

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	var commentsRes []models.CommentResponse
	for _, element := range post.Comment {
		commentRes := &models.CommentResponse{}
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = client.Collection("User").Doc(userID).ID
		commentRes.PostId = client.Collection("Post").Doc(postID).ID
		commentRes.CountLike = len(element.Like)
		commentRes.Date = element.Date
		commentsRes = append(commentsRes, *commentRes)
	}
	c.JSON(http.StatusOK, commentsRes)
}

// R (read a comment by post id and comment id)
func NewGetCommentByID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	commentRes := models.CommentResponse{}

	// get Post by id
	post, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	if len(post.Comment) < 1 {
		c.JSON(http.StatusOK, post.Comment)
		return
	}

	for _, element := range post.Comment {
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = client.Collection("User").Doc(userID).ID
		commentRes.PostId = client.Collection("Post").Doc(postID).ID
		commentRes.Date = element.Date
		commentRes.CountLike = len(element.Like)
	}
	c.JSON(http.StatusOK, commentRes)
}

// R (read all comment by post id)
func NewGetAllCommentByPostID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	post := models.Post2{}
	commentsRes := make([]models.CommentResponse, 0, len(post.Comment))
	commentRes := &models.CommentResponse{}

	// get Post by id
	postDoc, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	for _, element := range postDoc.Comment {
		commentRes = &models.CommentResponse{}
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = client.Collection("User").Doc(userID).ID
		commentRes.PostId = client.Collection("Post").Doc(postID).ID
		commentRes.CountLike = len(element.Like)
		commentRes.Date = element.Date
		commentsRes = append(commentsRes, *commentRes)
	}
	c.JSON(http.StatusOK, commentsRes)
}

// U (update a comment)
func NewUpdateComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	commentID := c.Param("comment_id")
	comment := models.Comment2{}

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get data from body",
		})
		return
	}

	// get Post by id
	post, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	// find index
	for _, element := range post.Comment {
		if element.CommentID == commentID {
			element.Content = comment.Content
			print(element.Content)
			break
		}
	}

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	c.JSON(http.StatusOK, post.Comment)
}

// D (delete a comment)
func NewDeleteComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	commentID := c.Param("comment_id")

	// get Post by id
	post, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	// find index
	for i, element := range post.Comment {
		if element.CommentID == commentID {
			post.Comment = append(post.Comment[:i], post.Comment[i+1:]...)
			break
		}
	}

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	c.JSON(http.StatusOK, post.Comment)
}

// this code about redundenc about code

// getPost retrieves the post with the specified ID from Firestore.
func getPost(ctx context.Context, client *firestore.Client, postID string) (*models.Post2, error) {
	doc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var post models.Post2
	if err := doc.DataTo(&post); err != nil {
		return nil, err
	}

	return &post, nil
}

// updatePost updates the post with the specified ID in Firestore.
func updatePost(ctx context.Context, client *firestore.Client, postID string, post *models.Post2) error {
	_, err := client.Collection("Post").Doc(postID).Set(ctx, post)
	return err
}

// errorCheck sends a JSON error response to the client if the specified error is not nil.
func errorCheck(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
	}
}
