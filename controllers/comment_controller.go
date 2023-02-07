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

	id := c.Request.Header.Get("id")
	post := models.Post{}
	dsnap, err2 := client.Collection("Post").Doc(id).Get(ctx)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	post.Comment = append(post.Comment, commentRef)
	setData, _ := client.Collection("Post").Doc(id).Set(ctx, post)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find postid",
		})
	}
	c.JSON(http.StatusOK, setData)

	//------Updating to User
	// id := c.Request.Header.Get("id")
	// user := models.User{}
	// dsnap, err2 := client.Collection("User").Doc(id).Get(ctx)
	// 	if err2 != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Cant to find userid",
	// 	})
	// }
	// mapstructure.Decode(dsnap.Data(), &user)
	// user.Comment = append(user.Comment, commentRef)
	// setData, _ := client.Collection("User").Doc(id).Set(ctx, user)
	// if err2 != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Cant to find userid",
	// 	})
	// }
	// 		//----------- return data ---------------
	// 		c.JSON(http.StatusOK, setData)

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

	// iter := client.Collection("Comment").Documents(ctx)
	// for {
	// 	doc, err := iter.Next()
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	comment := models.Comment{}
	// 	mapstructure.Decode(doc.Data(), &comment)
	// 	comments = append(comments, comment)
	// }

}

func GetMyComment(c *gin.Context) {
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

// ---------- another way to get comment -------------

// func GetComment(c *gin.Context) {
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)
// 	comment := models.Comment{}
// 	id := c.Request.Header.Get("id")
// 	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Cant to find commentid",
// 		})
// 	}
// 	mapstructure.Decode(dsnap.Data(), &comment)
// 	c.JSON(http.StatusOK, comment)
// }

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

func GetCommentByPostID(c *gin.Context) {
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
