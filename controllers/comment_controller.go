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

	// get post doucument by post id
	postRef := client.Collection("Post").Doc(post_id)
	dsnap, err2 := postRef.Get(ctx)
	if err2 != nil {
		log.Println(err2)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error finding post"})
		return
	}
	mapstructure.Decode(dsnap.Data(), &post)

	// get user document by user iid
	userRef := client.Collection("User").Doc(user_id)
	dsnap2, err3 := userRef.Get(ctx)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error finding post"})
		return
	}
	mapstructure.Decode(dsnap2.Data(), &user)

	// assigning comment object to comment model
	comment.Content = string(comment.Content)
	comment.UserId = userRef
	comment.PostId = postRef
	comment.Like = []*firestore.DocumentRef{}
	comment.Date = time.Now().UTC()

	// start the transaction
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		// add document [ commnt colleection ]
		commentRef, _, err := client.Collection("Comment").Add(ctx, comment)
		if err != nil {
			return err
		}

		// update field comment [ post collection ]
		post.Comment = append(post.Comment, commentRef)
		_, err = postRef.Set(ctx, post)
		if err != nil {
			return err
		}

		// update field comment [ user collection ]
		user.Comments = append(user.Comments, commentRef)
		_, err = userRef.Set(ctx, user)
		if err != nil {
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

// service Comment by Post id
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

		// commentRes.UserId = comment.UserId.ID
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

	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	id := c.Request.Header.Get("id")
	_, err := client.Collection("Comment").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete comment successfully",
	})
}

// ยากชิบหายเลยยยย อัพเดพเนี่ย
func UpdateComment(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	id := c.Request.Header.Get("id")
	comment := models.Comment{}

	c.BindJSON(&comment)
	_, err := client.Collection("Comment").Doc(id).Set(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update comment successfully",
	})
}

// func UpdateCommentByAof(c *gin.Context) {

// 	comment_id := c.Param("comment_id")
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	fmt.Printf("time.Now().UTC():%v\n", time.Now().Format(time.RFC3339))

// 	currentTime := time.Now().Format(time.RFC3339)
// 	currentDateTime, err := time.Parse(time.RFC3339, currentTime)

// 	commentDoc, err := client.Collection("Comment").Doc(comment_id).Get(c.Request.Context())
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Cant to find commentid"})
// 		return
// 	}

// 	var comment models.Comment
// 	if err := commentDoc.DataTo(&comment); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"Cant to find commentid": err.Error()})
// 		return

// 	}

// 	var updatedComment models.Comment
// 	if err := c.ShouldBindJSON(&updatedComment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	if updatedComment.Content != "" {
// 		comment.Content = updatedComment.Content
// 	}

// 	if !updatedComment.Date.IsZero() {
// 		comment.Date = updatedComment.Date
// 	}

// 	comment.IsPublic = updatedComment.IsPublic
// 	comment.Date = currentDateTime
// 	comment.Like = updatedComment.Like
// 	if len(comment.Like) == 0 {
// 		comment.Like = []*firestore.DocumentRef{}
// 	}

// 	if _, err := client.Collection("Comment").Doc(comment_id).Set(c.Request.Context(), comment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Cant to find commentid",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Update comment successfully",
// 	})

// }

func LikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	id := c.Request.Header.Get("id")
	comment := models.Comment{}
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	// comment.Like = append(comment.Like, id)
	_, err = client.Collection("Comment").Doc(id).Set(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Like comment successfully",
	})
}

func DislikeComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	id := c.Request.Header.Get("id")
	comment := models.Comment{}
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	// comment.Like = append(comment.Like, id)
	_, err = client.Collection("Comment").Doc(id).Set(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Dislike comment successfully",
	})
}

func GetCommentByUserID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func GetCommentByPostIDAndUserID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func GetCommentByPostIDAndUserIDAndCommentID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func GetCommentByPostIDAndCommentID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func GetCommentByUserIDAndCommentID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}

func GetCommentByCommentID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &comment)
	c.JSON(http.StatusOK, comment)
}
