package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
<<<<<<< HEAD

	mapstructure "github.com/mitchellh/mapstructure"

	"net/http"

	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
)

func GetAllPost(c *gin.Context) {

	posts := []models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("Post").Documents(ctx)
	for {
		post := models.Post{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "Not found")
		}
		mapstructure.Decode(doc.Data(), &post)
		posts = append(posts, post)
	}

	fmt.Println(posts)
	c.JSON(http.StatusOK, posts)
}

func GetPostById(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := client.Collection("Post").Add(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	_, err := client.Collection("Post").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant delete",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := client.Collection("Post").Doc(id).Set(ctx, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant update",
		})
	}
	c.JSON(http.StatusOK, post)
}

func GetPostByUserId(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}

func GetPostByCategoryId(c *gin.Context) {
	id := c.Param("id")
	post := models.Post{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Post").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}
	mapstructure.Decode(dsnap.Data(), &post)
	c.JSON(http.StatusOK, post)
}
=======
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

// ------------- unused -------------

// update user

// delete user
>>>>>>> 3a98ee4586c78a8b00c2340751b5fe9fced8ddb1
