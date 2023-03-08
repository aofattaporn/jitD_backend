package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func AddBookmark(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get user data
	var user models.User

	// get id that's to use
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// get user data from User collectiion
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		// handle error
	}
	if err := mapstructure.Decode(userDoc.Data(), &user); err != nil {
		// handle error
	}

	// check post to add bookmark is to be
	postRef := client.Collection("Post").Doc(postID)

	// get data to add
	user.BookMark = append(user.BookMark, postRef)

	// set user to collection
	_, err = client.Collection("User").Doc(userID).Set(ctx, user)
	if err != nil {
		// handle error
	}
}

func GetBookmarks(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	userID := c.Request.Header.Get("id")

	// get user from user colleection
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve user document",
			"error":   err.Error(),
		})
		return
	}

	// Decode user data from Firestore document
	var user models.User
	if err := userDoc.DataTo(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to decode user data",
			"error":   err.Error(),
		})
		return
	}

	// Get bookmarked posts from Firestore
	postRefs := user.BookMark
	postDocs, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve bookmarked posts",
			"error":   err.Error(),
		})
		return
	}

	// Decode post data from Firestore documents
	var posts []models.PostResponse
	for _, doc := range postDocs {
		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to decode post data",
				"error":   err.Error(),
			})
			return
		}
		posts = append(posts, models.PostResponse{
			UserId:       post.UserID.ID,
			PostId:       doc.Ref.ID,
			Content:      post.Content,
			Date:         post.Date,
			IsPublic:     post.IsPublic,
			Category:     post.Category,
			CountComment: len(post.Comment),
			CountLike:    len(post.LikesRef),
			IsLike:       false,
		})
	}

	// Return success response with bookmarked posts
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved bookmarked posts",
		"data":    posts,
	})
}

func AddBookmark2(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// update the array field
	_, err := client.Collection("User").Doc(userID).Update(ctx, []firestore.Update{
		{
			Path:  "Bookmark",
			Value: firestore.ArrayUnion(client.Collection("Post").Doc(postID)),
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update bookmarked posts",
			"error":   err.Error(),
		})
		return
	}

	// reeturn status success
	c.JSON(http.StatusOK, gin.H{
		"message": "add book mark success",
		"data":    "Add post bookmarked successfully",
	})

}

func Remove2(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// update the array field
	_, err := client.Collection("User").Doc(userID).Update(ctx, []firestore.Update{
		{
			Path:  "Bookmark",
			Value: firestore.ArrayRemove(client.Collection("Post").Doc(postID)),
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update bookmarked posts",
			"error":   err.Error(),
		})
		return
	}

	// reeturn status success
	c.JSON(http.StatusOK, gin.H{
		"message": "add book mark success",
		"data":    "Remove post bookmarked successfully",
	})

}
