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

// service create post
func CreatePost(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post := models.Post{}

	// Call BindJSON to bind the received JSON body
	if err := c.BindJSON(&post); err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("time.Now().UTC(): %v\n", time.Now().Format(time.RFC3339))

	//----------- adding post data to Posts ---------------
	post.Date.Format(time.RFC3339)
	currentTime := time.Now().Format(time.RFC3339)
	currentDateTime, err := time.Parse(time.RFC3339, currentTime)

	post.Date = currentDateTime
	post.Comment = []*firestore.DocumentRef{}
	post.Like = []*firestore.DocumentRef{}
	fmt.Printf("post: %v\n", post)
	postRef, _, err := client.Collection("Post").Add(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to create post",
		})
	}

	//----------- updating to user ---------------
	id := c.Request.Header.Get("id")
	user := models.User{}
	dsnap, err2 := client.Collection("User").Doc(id).Get(ctx)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}
	mapstructure.Decode(dsnap.Data(), &user)
	user.Posts = append(user.Posts, postRef)
	setData, _ := client.Collection("User").Doc(id).Set(ctx, user)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}

	//----------- return data ---------------
	c.JSON(http.StatusOK, setData)
}

// service get all post
func GetAllPost(c *gin.Context) {
	posts := []models.PostResponse{}

	// post := models.Post{}
	postRes := models.PostResponse{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	snap, err := client.Collection("Post").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, element := range snap {
		// fmt.Println(element.Data())
		id := c.Request.Header.Get("id")
		post := models.Post{}
		mapstructure.Decode(element.Data(), &post)
		mapstructure.Decode(post, &postRes)

		postRes.UserId = id
		postRes.PostId = element.Ref.ID
		postRes.CountLike = len(post.Like)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.Date = post.Date

		posts = append(posts, postRes)
		// fmt.Printf("post: %v\n", post)
	}
	c.JSON(http.StatusOK, posts)
}

// service get my post
func GetMyPost(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	//----------- finding my id user ---------------
	id := c.Param("id")
	user := models.User{}
	dsnap, err2 := client.Collection("User").Where("UserID", "==", id).Documents(ctx).GetAll()
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}

	//----------- finding post from data user ---------------
	// data, err3 := client.Collection("Post").
	mapstructure.Decode(dsnap[0].Data(), &user)
	post := models.Post{}
	posts := []models.Post{}

	X, _ := client.GetAll(ctx, user.Posts)
	for _, element := range X {
		mapstructure.Decode(element.Data(), &post)
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

// ------------- unused -------------

// update user
