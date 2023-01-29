package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"
	"sort"
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
	id := c.Request.Header.Get("id")

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

	dsnap, err2 := client.Collection("User").Doc(id).Get(ctx)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
	}
	post.UserID = dsnap.Ref
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
	user := models.User{}
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
	postRes := models.PostResponse{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	snap, err := client.Collection("Post").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, element := range snap {
		post := models.Post{}
		mapstructure.Decode(element.Data(), &post)
		mapstructure.Decode(post, &postRes)

		postRes.UserId = post.UserID.ID
		postRes.PostId = element.Ref.ID
		postRes.CountLike = len(post.Like)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.Date = post.Date

		posts = append(posts, postRes)
	}
	c.JSON(http.StatusOK, posts)
}

// service get my post
func GetMyPost(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	//----------- finding my id user ---------------
	id := c.Request.Header.Get("id")
	user := models.User{}

	dsnap, err := client.Collection("User").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	postResf := user.Posts
	postRes := models.PostResponse{}
	postsRes := []models.PostResponse{}
	post := models.Post{}

	postData, typeerr := dsnap.DataAt("Posts")
	if typeerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cant get get type information",
		})
	}
	mapstructure.Decode(postData, &postResf)

	X, _ := client.GetAll(ctx, postResf)
	for _, element := range X {
		mapstructure.Decode(element.Data(), &post)
		mapstructure.Decode(post, &postRes)
		postRes.UserId = id
		postRes.PostId = element.Ref.ID
		postRes.Date = post.Date
		postRes.CountComment = len(post.Comment)
		postRes.CountLike = len(post.Like)

		postsRes = append(postsRes, postRes)
	}

	c.JSON(http.StatusOK, postsRes)
}

// update user

func DeleteMyPost(c *gin.Context) {
	post_id := c.Param("post_id")
	user_id := c.Request.Header.Get("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	userData, err := client.Collection("User").Doc(user_id).Get(ctx)
	user := models.User{}
	mapstructure.Decode(userData.Data(), &user)
	user.Posts = RemoveCopy(user.Posts, post_id)

	_, setErr := client.Collection("User").Doc(user_id).Set(ctx, user)
	if setErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete cant to delete data",
		})
		return
	}

	// delete document post
	_, err2 := client.Collection("Post").Doc(post_id).Delete(ctx)
	if err2 != nil {
		log.Printf("An error has occurred: %s", err)
	} else {

		c.JSON(http.StatusOK, gin.H{
			"message": "delete post success",
		})
	}
}

func RemoveCopy(slice []*firestore.DocumentRef, id string) []*firestore.DocumentRef {
	var idx int
	for index, dr := range slice {
		if dr.ID == id {
			idx = index
			break
		}
	}
	slice = append(slice[:idx], slice[idx+1:]...)
	return slice
}

func binarySearch(slice []*firestore.DocumentRef, id string) []*firestore.DocumentRef {

	i := sort.Search(len(slice), func(i int) bool { return id <= slice[i].ID })
	if i < len(slice) && slice[i].ID == id {
		fmt.Printf("Found %s at index %d in %v.\n", id, i, slice)
	} else {
		fmt.Printf("Did not find %s in %v.\n", id, slice)
	}
	slice = append(slice[:i], slice[i+1:]...)
	return slice
}

// ------------- unused -------------

// delete user
