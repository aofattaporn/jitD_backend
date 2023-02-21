package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"context"
// 	"fmt"
// 	configs "jitD/configs"
// 	models "jitD/models"
// 	"net/http"
// 	"time"

// 	"cloud.google.com/go/firestore"
// 	"github.com/gin-gonic/gin"
// 	mapstructure "github.com/mitchellh/mapstructure"
// )

// // LikePost creates a user like on a post
// func LikePost(c *gin.Context) {
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)
// 	var post models.Post

// 	userID := c.Request.Header.Get("id")
// 	postID := c.Param("post_id")

// 	// get post ref
// 	postRef, err := client.Collection("Post").Doc(postID).Get(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Cannot find user ID",
// 		})
// 		return
// 	}
// 	mapstructure.Decode(postRef.Data(), &post)

// 	like := models.Like{
// 		UserRef: client.Collection("User").Doc(userID),
// 		PostRef: postRef.Ref,
// 		Date:    time.Now().UTC(),
// 	}

// 	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

// 		// Create a like
// 		_, _, err := client.Collection("Like").Add(ctx, like)
// 		if err != nil {
// 			return err
// 		}

// 		// Update post
// 		postRef := client.Collection("Post").Doc(postID)
// 		postSnap, err := tx.Get(postRef)
// 		if err != nil {
// 			return err
// 		}

// 		err = postSnap.DataTo(&post)
// 		if err != nil {
// 			return err
// 		}

// 		post.LikesRef = append(post, &like)
// 		err = tx.Set(postRef, post)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": fmt.Sprintf("Something went wrong: %v", err),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Post liked successfully",
// 	})
// }

// // UnlikePost removes a user like from a post
// func UnlikePost(c *gin.Context) {
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)
// 	post := models.Post{}

// 	userID := c.Request.Header.Get("id")
// 	postID := c.Param("post_id")

// 	// get user ref
// 	userRef := client.Collection("User").Doc(userID)

// 	// get post ref
// 	postRef := client.Collection("Post").Doc(postID)
// 	postDocsnap, errPost := postRef.Get(ctx)
// 	mapstructure.Decode(postDocsnap.Data(), &post)
// 	if errPost != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": fmt.Sprintf("Something went wrong: %v", errPost),
// 		})
// 		return
// 	}

// 	// check if the user has liked the post
// 	likeQuery := client.Collection("Like").Where("UserRef", "==", userRef).Where("PostRef", "==", postRef).Limit(1)
// 	likeDocs, err := likeQuery.Documents(ctx).GetAll()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": fmt.Sprintf("Something went wrong: %v", err),
// 		})
// 		return
// 	}

// 	// delete the user like and update the post
// 	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
// 		if len(likeDocs) > 0 {
// 			likeRef := likeDocs[0].Ref

// 			// delete the user like
// 			err = tx.Delete(likeRef)
// 			if err != nil {
// 				return err
// 			}

// 			for i, likeRef := range post.LikesRef {
// 				if likeRef.PostRef.Path == likeRef.PostRef.Path {
// 					post.LikesRef = append(post.LikesRef[:i], post.LikesRef[i+1:]...)
// 					break
// 				}
// 			}

// 			err = tx.Set(postRef, post)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": fmt.Sprintf("Something went wrong: %v", err),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Post unliked successfully",
// 	})
// }

// Like Comment
func LikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post := models.Post{}
	index := -1

	userID := c.Request.Header.Get("id")
	commentID := c.Param("comment_id")
	postID := c.Param("post_id")

	// get post ref
	postDocsnap, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find post ID",
		})
		return
	}
	if err := postDocsnap.DataTo(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error decoding post data",
		})
		return
	}
	// loop for find post.comment
	for i, c2 := range post.Comment {
		if c2.CommentID == commentID {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "comment id does not exist",
		})
		return
	}

	for _, c2 := range post.Comment[index].Like {
		if c2.CommentID == commentID {
			c.JSON(http.StatusOK, gin.H{
				"message": "like aleady exist",
			})
			return
		}
	}

	post.Comment[index].Like = append(post.Comment[index].Like, &models.LikeComment{
		UserID:    userID,
		CommentID: commentID,
		Date:      time.Now().UTC(),
	})

	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot to set post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment liked successfully",
	})
}

func UnLikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	commentID := c.Param("comment_id")
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// Get post document
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find post ID",
		})
		return
	}

	// Decode post document data to post struct
	var post models.Post
	if err := postDoc.DataTo(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot decode post data",
		})
		return
	}

	// Find comment index in post's comments slice
	index := -1
	for i, comment := range post.Comment {
		if comment.CommentID == commentID {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find comment ID",
		})
		return
	}

	// Find user's like for the comment
	likeIndex := -1
	for i, like := range post.Comment[index].Like {
		if like.CommentID == commentID && like.UserID == userID {
			likeIndex = i
			break
		}
	}
	if likeIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find user's like for the comment",
		})
		return
	}

	// Remove user's like from comment's likes slice
	post.Comment[index].Like = append(post.Comment[index].Like[:likeIndex], post.Comment[index].Like[likeIndex+1:]...)

	// Update post document with modified post struct
	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot set post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment unliked successfully",
	})
}
