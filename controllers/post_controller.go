package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"

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

	//----------- adding post data to Posts ---------------
	post.Category = []string{}
	post.Comment = []*firestore.DocumentRef{}
	post.Like = []*firestore.DocumentRef{}
	postRef, _, err := client.Collection("Post").Add(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to create post",
		})
	}

	//----------- updating to user ---------------
	id := c.Param("id")
	user := models.User{}
	dsnap, err2 := client.Collection("User").Where("UserID", "==", id).Documents(ctx).GetAll()
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}
	mapstructure.Decode(dsnap[0].Data(), &user)
	user.Posts = append(user.Posts, postRef)
	setData, _ := client.Collection("User").Doc(dsnap[0].Ref.ID).Set(ctx, user)

	//----------- return data ---------------
	c.JSON(http.StatusOK, setData)
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

// service get all post
func GetAllPost(c *gin.Context) {
	posts := []models.Post{}
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	fmt.Printf("c.Request.Header.Get(\"id\"): %v\n", c.Request.Header.Get("id"))

	snap, err := client.Collection("Post").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, element := range snap {
		fmt.Println(element.Data())
		mapstructure.Decode(element.Data(), &post)
		posts = append(posts, post)
	}
	c.JSON(http.StatusOK, posts)
}

// service deleing post
