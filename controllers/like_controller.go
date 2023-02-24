package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// LikePost creates a user like on a post
func LikePost(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare id to use in this function
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// declare object to use in this function
	var post models.Post

	// crete like object for save to db
	like := models.Like{
		UserRef: client.Collection("User").Doc(userID),
		PostRef: client.Collection("Post").Doc(postID),
		Date:    time.Now().UTC(),
	}

	likeExis, errExist := client.Collection("Like").Where("UserRef", "==", client.Collection("User").Doc(userID)).Where("PostRef", "==", client.Collection("Post").Doc(postID)).Documents(ctx).GetAll()
	if errExist != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to access like collection",
		})
		return
	}

	if len(likeExis) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "you already like this post",
		})
		return
	}

	// run transaction
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		// Create a like
		_, _, err := client.Collection("Like").Add(ctx, like)
		if err != nil {
			return err
		}

		// get post data by post id to map
		postRef := client.Collection("Post").Doc(postID)
		postSnap, err := tx.Get(postRef)
		if err != nil {
			return err
		}

		// map data to post
		err = postSnap.DataTo(&post)
		if err != nil {
			return err
		}

		// append to post like ref and set to DDB
		post.LikesRef = append(post.LikesRef, &like)
		err = tx.Set(postRef, post)
		if err != nil {
			return err
		}

		return nil
	})

	// check err from transaction
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

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare id to use in this function
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// declare object to use in this function
	post := models.Post{}

	// get referenc of postID and userID
	userRef := client.Collection("User").Doc(userID)
	postRef := client.Collection("Post").Doc(postID)

	// get post data from post ref and map to post
	postDocsnap, errPost := postRef.Get(ctx)
	if errPost != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", errPost),
		})
		return
	}

	// save data to post
	mapstructure.Decode(postDocsnap.Data(), &post)

	// save like data to post.like
	if likeDoc, err := postDocsnap.DataAt("LikesRef"); err == nil {
		if err := mapstructure.Decode(likeDoc, &post.LikesRef); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to decode like of post",
			})
			return
		}
	}

	// check if the user has liked the post
	likeQuery := client.Collection("Like").Where("UserRef", "==", userRef).Where("PostRef", "==", postRef)
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

			// remove all like from like collection
			for i, _ := range likeDocs {
				likeRef := likeDocs[i].Ref
				// delete the user like
				err = tx.Delete(likeRef)
				if err != nil {
					return err
				}
			}

			// remove lke from post data
			for i, like := range post.LikesRef {
				if like.UserRef.ID == userID && like.PostRef.ID == postID {
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

	// check err from backend
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Something went wrong: %v", err),
		})
		return
	}

	// returrn json to frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
	})
}

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

func LikePost2(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare id to use in this function
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

	// declare object to use in this function
	var post models.Post

	// crete like object for save to db
	like := models.Like{
		UserRef: client.Collection("User").Doc(userID),
		PostRef: client.Collection("Post").Doc(postID),
		Date:    time.Now().UTC(),
	}

	fmt.Printf("postID: %v\n", postID)

	// save to post collection
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": "cant to get post by this id ก",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)

	// add to new data
	_, _, err = client.Collection("Like").Add(ctx, like)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "cant to set like to like collection",
		})
		return
	}

	post.LikesRef = append(post.LikesRef, &like)
	_, err = client.Collection("Post").Doc(postID).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "cant to set like to like collection",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post liked successfully",
	})
}
