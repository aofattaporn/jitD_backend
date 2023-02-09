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
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	comment := models.Comment{}

	if err := c.BindJSON(&comment); err != nil {
		log.Fatalln(err)
		return

	}
	fmt.Println("time.Now().UTC(): ", time.Now().Format(time.RFC3339))

	// assigning comment object to comment model
	comment.Content = string(comment.Content)
	comment.Like = []*firestore.DocumentRef{}
	comment.UserId = string(comment.UserId)
	comment.Date.Format(time.RFC3339)
	//currentTime := time.Now().Format(time.RFC3339)

	//add comment to comment collection
	commentRef, _, err := client.Collection("Comment").Add(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to create comment",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create comment success",
		})
	}

	//------- Updateing to post --------------

	post_id := c.Param("post_id")
	post := models.Post{}
	dsnap, err2 := client.Collection("Post").Doc(post_id).Get(ctx)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find postid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	post.Comment = append(post.Comment, commentRef)
	_, err3 := client.Collection("Post").Doc(post_id).Set(ctx, post)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find postid",
		})
	}
	//c.JSON(http.StatusOK, _)

	//------- Updateing to User --------------

	user_id := c.Request.Header.Get("id")
	user := models.User{}
	dataUser, err2 := client.Collection("User").Doc(user_id).Get(ctx)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}
	mapstructure.Decode(dataUser.Data(), &user)
	user.Comments = append(user.Comments, commentRef)
	setDataUser, _ := client.Collection("User").Doc(user_id).Set(ctx, user)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find postid",
		})
	}
	c.JSON(http.StatusOK, setDataUser)

}

func GetAllComment(c *gin.Context) {

	comments := []models.CommentResponse{}
	commentRes := models.CommentResponse{}
	comment := models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	snap, err := client.Collection("Comment").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	for _, element := range snap {
		mapstructure.Decode(element.Data(), &comment)

		commentRes.UserId = comment.UserId
		commentRes.PostId = element.Ref.Parent.ID
		commentRes.Content = comment.Content
		commentRes.CommentId = element.Ref.ID
		commentRes.Date = comment.Date
		commentRes.CountLike = len(comment.Like)

		comments = append(comments, commentRes)
	}
	c.JSON(http.StatusOK, comments)

}

func GetMyComment(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	user := models.User{}
	id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("User").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}

	commentResf := user.Comments
	commentRes := models.CommentResponse{}
	comment := models.Comment{}
	commentsRes := []models.CommentResponse{}

	commentData, typeerr := dsnap.DataAt("Comments")
	if typeerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(commentData, &commentResf)

	X, _ := client.GetAll(ctx, commentResf)
	for _, element := range X {
		mapstructure.Decode(element.Data(), &comment)
		mapstructure.Decode(comment, &commentRes)
		commentRes.UserId = comment.UserId
		commentRes.PostId = element.Ref.Parent.ID
		commentRes.Content = comment.Content
		commentRes.CommentId = element.Ref.ID
		commentRes.Date = comment.Date
		commentRes.CountLike = len(comment.Like)
		commentsRes = append(commentsRes, commentRes)
	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentByPostID(c *gin.Context) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post_id := c.Param("post_id")
	post := models.Post{}
	// comment := models.Comment{}
	// id := c.Request.Header.Get("id")
	dsnap, err := client.Collection("Post").Doc(post_id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find commentid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post.Comment)
}

// ยากแล้ว
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
