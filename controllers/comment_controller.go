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
func CreateComment(c *gin.Context) {
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

	var newComment = &models.Comment2{
		CommentID: uuid.NewString(),
		Content:   comment.Content,
		UserId:    client.Collection("User").Doc(userID).ID,
		Like:      []*models.LikeComment{},
		Date:      time.Now().UTC(),
	}

	// Add comment to post
	post.Comment = append(post.Comment, newComment)

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	UpdateProgressQuest(c, "CommentQuest")

	c.JSON(http.StatusOK, &models.CommentResponse{
		UserId:    userID,
		PostId:    postID,
		CommentId: newComment.CommentID,
		Content:   comment.Content,
		CountLike: 0,
		Date:      newComment.Date,
		IsLike:    false,
		IsPin:     false,
	})

}

// R (read all comment by post id)
func GetAllCommentByPostID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")
	post := models.Post{}
	commentsRes := make([]models.CommentResponse, 0, len(post.Comment))

	// get Post by id
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	commentsRes = getCommentsResponse(post, client, userID, postID)
	c.JSON(http.StatusOK, commentsRes)
}

// U (update a comment)
func UpdateComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	postID := c.Param("post_id")
	commentID := c.Param("comment_id")
	userID := c.Request.Header.Get("id")

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

	var indexUpdate int
	// find index
	for index, _ := range post.Comment {
		if post.Comment[index].CommentID == commentID {
			indexUpdate = index
			post.Comment[index].Content = comment.Content
			break
		}
	}

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	c.JSON(http.StatusOK, &models.CommentResponse{
		UserId:    userID,
		PostId:    postID,
		CommentId: commentID,
		Content:   post.Comment[indexUpdate].Content,
		CountLike: len(post.Comment[indexUpdate].Like),
		Date:      post.Comment[indexUpdate].Date,
		IsLike:    false,
		IsPin:     post.Comment[indexUpdate].IsPin,
	})
}

// D (delete a comment)
func DeleteComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	postID := c.Param("post_id")
	commentID := c.Param("comment_id")
	userID := c.Request.Header.Get("id")

	// get Post by id
	post, err := getPost(ctx, client, postID)
	if err != nil {
		errorCheck(c, err, "cant to find post id")
		return
	}

	// get comment delete
	var deleteComment models.CommentResponse
	// find index
	for i, element := range post.Comment {
		if element.CommentID == commentID {
			deleteComment.UserId = userID
			deleteComment.PostId = postID
			deleteComment.CommentId = commentID
			deleteComment.Content = element.Content
			deleteComment.CountLike = len(element.Like)
			deleteComment.Date = element.Date
			deleteComment.IsLike = false

			post.Comment = append(post.Comment[:i], post.Comment[i+1:]...)
			break

		}
	}

	// Update post in Firestore
	if err := updatePost(ctx, client, postID, post); err != nil {
		errorCheck(c, err, "cant to add a comment")
		return
	}

	c.JSON(http.StatusOK, deleteComment)
}

// getPost retrieves the post with the specified ID from Firestore.
func getPost(ctx context.Context, client *firestore.Client, postID string) (*models.Post, error) {
	doc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var post models.Post
	mapstructure.Decode(doc.Data(), &post)

	return &post, nil
}

// updatePost updates the post with the specified ID in Firestore.
func updatePost(ctx context.Context, client *firestore.Client, postID string, post *models.Post) error {
	_, err := client.Collection("Post").Doc(postID).Set(ctx, post)
	return err
}

// errorCheck sends a JSON error response to the client if the specified error is not nil.
func errorCheck(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
	}
}

func getCommentsResponse(post models.Post, client *firestore.Client, userID string, postID string) []models.CommentResponse {
	commentRes := &models.CommentResponse{}
	commentsRes := make([]models.CommentResponse, 0, len(post.Comment))

	for _, element := range post.Comment {
		commentRes = &models.CommentResponse{}
		mapstructure.Decode(element, &commentRes)
		commentRes.UserId = element.UserId
		commentRes.PostId = client.Collection("Post").Doc(postID).ID
		commentRes.CountLike = len(element.Like)
		commentRes.Date = element.Date
		commentRes.IsLike = checkISLike(element.Like, userID)
		commentRes.IsPin = element.IsPin

		commentsRes = append(commentsRes, *commentRes)
	}

	return commentsRes
}

func checkISLike(userLike []*models.LikeComment, userID string) bool {
	for _, lc := range userLike {
		if lc.UserID == userID {
			return true
		}
	}
	return false
}
