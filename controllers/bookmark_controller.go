package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

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

	var user models.User
	// Decode user data from Firestore document
	bookmarkRefs, ok := userDoc.Data()["Bookmark"].([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve bookmarked post references. Bookmark field not found in the user document.",
		})
		return
	}

	mapstructure.Decode(userDoc.Data(), &user)

	if len(bookmarkRefs) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully retrieved bookmarked posts",
			"data":    []models.PostResponse{},
		})
		return
	}

	// Get bookmarked posts from Firestore
	mapstructure.Decode(bookmarkRefs, &user.BookMark)

	postDocs, err := client.GetAll(ctx, user.BookMark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve bookmarked posts",
			"error":   err.Error(),
		})
		return
	}

	// Decode post data from Firestore documents
	posts := []models.PostResponse{}

	// get user data
	for _, doc := range postDocs {
		fmt.Println("post")

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
			IsLike:       checkIsLikePost(post.LikesRef, userID),
			IsBookmark:   checkIsBookMark(doc.Ref.ID, user.BookMark),
		})
	}

	// Return success response with bookmarked posts
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved bookmarked posts",
		"data":    posts,
	})
}

func AddBookmark(c *gin.Context) {

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

func Remove(c *gin.Context) {

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
