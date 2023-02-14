package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
)

// service Create Comment
func CreateComment(c *gin.Context) {

	// declare a variable
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	user_id := c.Request.Header.Get("id")
	post_id := c.Param("post_id")
	comment := models.Comment{}
	post := models.Post{}
	user := models.User{}

	if err := c.BindJSON(&comment); err != nil {
		log.Fatalln(err)
		return
	}

	postRef := client.Collection("Post").Doc(post_id)
	userRef := client.Collection("User").Doc(user_id)

	// assigning comment object to comment model
	comment.Content = string(comment.Content)
	comment.UserId = userRef
	comment.PostId = postRef
	comment.Like = []*firestore.DocumentRef{}
	comment.Date = time.Now().UTC()

	// get user and post doccument
	userDoc, errUserDoc := userRef.Get(ctx)
	if errUserDoc != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find user id ",
		})
		return
	}
	mapstructure.Decode(userDoc.Data(), &user)
	postDoc, errPostDoc := postRef.Get(ctx)
	if errPostDoc != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}
	mapstructure.Decode(postDoc.Data(), &post)
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		// add document [ comment collection ]
		commentRef, _, err := client.Collection("Comment").Add(ctx, comment)
		if err != nil {
			return err
		}

		// update field comment [ post collection ]
		post.Comment = append(post.Comment, commentRef)
		postData := map[string]interface{}{"Comment": post.Comment}
		if err := tx.Set(postRef, postData, firestore.MergeAll); err != nil {
			return err
		}

		// update field comment [ user collection ]
		user.Comments = append(user.Comments, commentRef)
		userData := map[string]interface{}{"Comments": user.Comments}
		if err := tx.Set(userRef, userData, firestore.MergeAll); err != nil {
			return err
		}

		return nil
	})

	// if trancsaction have a err
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to create comment",
		})
		return
	}

	// none err and respone a status OK
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment created successfully",
	})
}

// service Comment 		  [ by Post id ]
func GetCommentByPostID(c *gin.Context) {

	// initail data
	post := models.Post{}
	commentResf := post.Comment
	comment := models.Comment{}
	commentRes := models.CommentResponse{}
	commentsRes := []models.CommentResponse{}
	post_id := c.Param("post_id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get user document [ user collection ]
	dsnap, err := client.Collection("Post").Doc(post_id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find post id",
		})
		return
	}

	// get only some field [ Comments feild ] -> [ documentRefs[] ]
	commentData, typeerr := dsnap.DataAt("Comment")
	if typeerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find comment id",
		})
		return
	}
	mapstructure.Decode(commentData, &commentResf)

	// map data DB to comment Response
	X, _ := client.GetAll(ctx, commentResf)
	for _, element := range X {
		mapstructure.Decode(element.Data(), &comment)
		commentRes.UserId = comment.UserId.ID
		commentRes.PostId = comment.PostId.ID
		commentRes.CommentId = element.Ref.ID
		commentRes.Content = comment.Content
		commentRes.Date = comment.Date
		commentRes.CountLike = len(comment.Like)
		commentsRes = append(commentsRes, commentRes)
	}

	c.JSON(http.StatusOK, commentsRes)
}

// service update comment [ by comment id ]
func UpdateComment(c *gin.Context) {

	// initail varaible
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	commentID := c.Param("comment_id")

	// binding data
	var comment models.Comment
	c.BindJSON(&comment)

	// get comment document [ by comment id]
	doc, err := client.Collection("Comment").Doc(commentID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find comment with the given ID",
		})
		return
	}

	var originalComment models.Comment
	if err := doc.DataTo(&originalComment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while reading comment data",
		})
		return
	}

	// Modify the desired fields
	originalComment.Content = comment.Content

	// Update the document
	_, err = client.Collection("Comment").Doc(commentID).Set(ctx, originalComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while updating comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
	})
}

// service delete comment [ by comment id ]
func DeleteComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	user_id := c.Request.Header.Get("id")
	comment_id := c.Param("comment_id")
	post_id := c.Param("post_id")

	postRef := client.Collection("Post").Doc(post_id)
	userRef := client.Collection("User").Doc(user_id)
	commentRef := client.Collection("Comment").Doc(comment_id)

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get post and user documents
		postDoc, err := tx.Get(postRef)
		if err != nil {
			return err
		}
		userDoc, err := tx.Get(userRef)
		if err != nil {
			return err
		}

		// Delete comment from post document
		postData := postDoc.Data()
		postCommentRefs := postData["Comment"].([]interface{})
		var updatedPostCommentRefs []*firestore.DocumentRef
		for _, ref := range postCommentRefs {
			docRef := ref.(*firestore.DocumentRef)
			if docRef.ID != comment_id {
				updatedPostCommentRefs = append(updatedPostCommentRefs, docRef)
			}
		}
		postData["Comment"] = updatedPostCommentRefs
		if err := tx.Set(postRef, postData, firestore.MergeAll); err != nil {
			return err
		}

		// Delete comment from user document
		userData := userDoc.Data()
		userCommentRefs := userData["Comments"].([]interface{})
		var updatedUserCommentRefs []*firestore.DocumentRef
		for _, ref := range userCommentRefs {
			docRef := ref.(*firestore.DocumentRef)
			if docRef.ID != comment_id {
				updatedUserCommentRefs = append(updatedUserCommentRefs, docRef)
			}
		}
		userData["Comments"] = updatedUserCommentRefs
		if err := tx.Set(userRef, userData, firestore.MergeAll); err != nil {
			return err
		}

		// Delete comment document
		if err := tx.Delete(commentRef); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete comment successfully",
	})

}

// ------------ service not use in froneend --------------

// service All Comment
func GetAllComment(c *gin.Context) {

	// initail data
	comments := []models.CommentResponse{}
	commentRes := models.CommentResponse{}
	comment := models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all comment [ comment collection ]
	snap, err := client.Collection("Comment").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find comment id",
		})
		return
	}

	// mapping data to commentRes
	for _, element := range snap {
		mapstructure.Decode(element.Data(), &comment)
		commentRes.UserId = comment.UserId.ID
		commentRes.PostId = comment.PostId.ID
		commentRes.CommentId = element.Ref.ID
		commentRes.Content = comment.Content
		commentRes.Date = comment.Date
		commentRes.CountLike = len(comment.Like)

		comments = append(comments, commentRes)
	}
	c.JSON(http.StatusOK, comments)
}

// service Comment Comment
func GetMyComment(c *gin.Context) {

	// initail data
	user := models.User{}
	commentResf := user.Comments
	comment := models.Comment{}
	commentRes := models.CommentResponse{}
	commentsRes := []models.CommentResponse{}
	user_id := c.Request.Header.Get("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get user document [ user collection ]
	dsnap, err := client.Collection("User").Doc(user_id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
		return
	}

	// get only some field [ Comments feild ] -> [ documentRefs[] ]
	commentData, typeerr := dsnap.DataAt("Comments")
	if typeerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(commentData, &commentResf)

	// map data DB to comment Response
	X, _ := client.GetAll(ctx, commentResf)
	for _, element := range X {
		mapstructure.Decode(element.Data(), &comment)
		commentRes.UserId = comment.UserId.ID
		commentRes.PostId = comment.PostId.ID
		commentRes.CommentId = element.Ref.ID
		commentRes.Content = comment.Content
		commentRes.Date = comment.Date
		commentRes.CountLike = len(comment.Like)
		commentsRes = append(commentsRes, commentRes)
	}

	c.JSON(http.StatusOK, commentsRes)
}
