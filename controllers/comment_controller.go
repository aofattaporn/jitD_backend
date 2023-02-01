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

	// Adding Comment data
	comment.Date.Format(time.RFC3339)
	currentTime := time.Now().Format(time.RFC3339)
	currentDateTime, err := time.Parse(time.RFC3339, currentTime)

	comment.Date = currentDateTime
	comment.Like = []*firestore.DocumentRef{}
	fmt.Printf("comment: %v\n", comment)
	commentRef, _, err := client.Collection("Comment").Add(ctx, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Can't to create comment",
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

	comments := []models.Comment{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("Comment").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			log.Fatalln(err)
		}
		comment := models.Comment{}
		mapstructure.Decode(doc.Data(), &comment)
		comments = append(comments, comment)
	}

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

// //New Create Comment
// func handleComments(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		handleGetComments(w, r)
// 	case "POST":
// 		handleAddComment(w, r)
// 	default:
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

// func handleGetComments(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(comment)
// }

// func handleAddComment(w http.ResponseWriter, r *http.Request) {
// 	var comment Comment
// 	err := json.NewDecoder(r.Body).Decode(&comment)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	comments = append(comments, comment)
// 	fmt.Fprintf(w, "Comment added successfully")
// }

// func GetAllComment(c *gin.Context) {

// 	comments := []models.Comment{}
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	iter := client.Collection("Comment").Documents(ctx)
// 	for {
// 		comment := models.Comment{}
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 			c.JSON(http.StatusNotFound, "Not found")
// 		}
// 		mapstructure.Decode(doc.Data(), &comment)
// 		comments = append(comments, comment)
// 	}

// 	fmt.Println(comments)
// 	c.JSON(http.StatusOK, comments)
// }

// func GetCommentById(c *gin.Context) {
// 	id := c.Param("id")
// 	comment := models.Comment{}
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	dsnap, err := client.Collection("Comment").Doc(id).Get(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "cant create",
// 		})
// 	}
// 	mapstructure.Decode(dsnap.Data(), &comment)
// 	c.JSON(http.StatusOK, comment)
// }

// func CreateComment(c *gin.Context) {
// 	comment := models.Comment{}
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	if err := c.ShouldBindJSON(&comment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	_, _, err := client.Collection("Comment").Add(ctx, comment)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "cant create",
// 		})
// 	}
// 	c.JSON(http.StatusOK, comment)
// }

// func UpdateComment(c *gin.Context) {
// 	id := c.Param("id")
// 	comment := models.Comment{}
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	if err := c.ShouldBindJSON(&comment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	_, err := client.Collection("Comment").Doc(id).Set(ctx, comment)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "cant update",
// 		})
// 	}
// 	c.JSON(http.StatusOK, comment)
// }

// func DeleteComment(c *gin.Context) {
// 	id := c.Param("id")
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	_, err := client.Collection("Comment").Doc(id).Delete(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "cant delete",
// 		})
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "deleted",
// 	})
// }
