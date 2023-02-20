package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
)

// LikePost creates a user like on a post
func LikePost(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	var post models.Post

	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// get user ref
	userRef, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find user ID",
		})
		return
	}

	// get post ref
	postRef, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find user ID",
		})
		return
	}
	mapstructure.Decode(postRef.Data(), &post)

	// check post have lke
	for _, l := range post.LikesRef {
		if l.UserRef.Path == userRef.Ref.Path {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": "User alredy like this post",
			})
			return
		}
	}

	like := models.Like{
		UserRef: userRef.Ref,
		PostRef: postRef.Ref,
		Date:    time.Now().UTC(),
	}

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		// Create a like
		_, _, err := client.Collection("Like").Add(ctx, like)
		if err != nil {
			return err
		}

		// Update post
		postRef := client.Collection("Post").Doc(postID)
		postSnap, err := tx.Get(postRef)
		if err != nil {
			return err
		}

		err = postSnap.DataTo(&post)
		if err != nil {
			return err
		}

		post.LikesRef = append(post.LikesRef, &like)
		err = tx.Set(postRef, post)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post liked successfully",
	})
}

// UnlikePost removes a user like from a post
func UnlikePost(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post := models.Post{}

	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// get user ref
	userRef := client.Collection("User").Doc(userID)

	// get post ref
	postRef := client.Collection("Post").Doc(postID)
	postDocsnap, errPost := postRef.Get(ctx)
	mapstructure.Decode(postDocsnap.Data(), &post)
	if errPost != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", errPost),
		})
		return
	}

	// check if the user has liked the post
	likeQuery := client.Collection("Like").Where("UserRef", "==", userRef).Where("PostRef", "==", postRef).Limit(1)
	likeDocs, err := likeQuery.Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	// delete the user like and update the post
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if len(likeDocs) > 0 {
			likeRef := likeDocs[0].Ref

			// delete the user like
			err = tx.Delete(likeRef)
			if err != nil {
				return err
			}

			for i, likeRef := range post.LikesRef {
				if likeRef.PostRef.Path == likeRef.PostRef.Path {
					post.LikesRef = append(post.LikesRef[:i], post.LikesRef[i+1:]...)
					break
				}
			}

			err = tx.Set(postRef, post)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
	})
}

// Like Comment
func LikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}

	userID := c.Request.Header.Get("id")
	commentID := c.Param("comment_id")

	// get user ref
	userRef := client.Collection("User").Doc(userID)

	// get post ref
	commentDocsnap, err := client.Collection("Comment").Doc(commentID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find user ID",
		})
		return
	}
	mapstructure.Decode(commentDocsnap.Data(), &comment)

	// check post have lke
	for _, l := range comment.Like {
		if l.UserRef.Path == userRef.Path {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": "User alredy like this post",
			})
			return
		}
	}

	like := models.LikeComment{
		UserRef:    userRef,
		CommentRef: commentDocsnap.Ref,
		Date:       time.Now().UTC(),
	}

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		comment.Like = append(comment.Like, &like)
		err = tx.Set(commentDocsnap.Ref, comment)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment liked successfully",
	})
}

// Like Comment
func UnLikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}

	userID := c.Request.Header.Get("id")
	commentID := c.Param("comment_id")

	// get user ref
	userRef := client.Collection("User").Doc(userID)

	// get post ref
	commentDocsnap, err := client.Collection("Comment").Doc(commentID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find user ID",
		})
		return
	}
	mapstructure.Decode(commentDocsnap.Data(), &comment)

	// check if the user has liked the post

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		for i, likeRef := range comment.Like {
			if likeRef.UserRef != userRef && likeRef.CommentRef != commentDocsnap.Ref {
				comment.Like = append(comment.Like[:i], comment.Like[i+1:]...)
				break
			}
		}

		err = tx.Set(commentDocsnap.Ref, comment)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment UnLiked successfully",
	})
}
