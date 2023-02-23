package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
)

// check isLike from post like ref
func checkIsLikePost(postLike []*models.Like, userID string) bool {
	for _, l := range postLike {
		if l.UserRef.ID == userID {
			return true
		}
	}
	return false
}

// service get all post
func GetAllPost(c *gin.Context) {

	// declare data object
	posts := []models.PostResponse{}
	postRes := models.PostResponse{}

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// get all post by limit 10 result
	allDocSnap, err := client.Collection("Post").Limit(10).Documents(ctx).GetAll()
	if err != nil {
		return
	}

	// loop data snap and decode data to post respone
	for _, doc := range allDocSnap {

		post := models.Post{}
		postRes.Category = []string{}
		mapstructure.Decode(doc.Data(), &post)

		// map data to seend to fronend
		postRes.UserId = post.UserID.ID
		postRes.PostId = doc.Ref.ID
		postRes.Content = post.Content
		postRes.Date = post.Date
		postRes.IsPublic = post.IsPublic
		postRes.CountLike = len(post.LikesRef)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.IsLike = checkIsLikePost(post.LikesRef, userID)

		posts = append(posts, postRes)
	}

	// return json status code 200
	c.JSON(http.StatusOK, posts)
}

// service create post
func CreatePost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all object variable
	post := models.Post{}
	today := time.Now().UTC()

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// mappig data body to object post
	if err := c.BindJSON(&post); err != nil {
		log.Fatalln(err)
		return
	}

	// TODO: Maybe it's hsould to check post

	// add object post to firestore DB
	posRef, _, err := client.Collection("Post").Add(ctx, models.Post{
		UserID:   client.Collection("User").Doc(userID),
		Content:  post.Content,
		Date:     today,
		IsPublic: post.IsPublic,
		Category: post.Category,
		Comment:  []*models.Comment2{},
		LikesRef: []*models.Like{},
	})
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to create post",
		})
		return
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, models.PostResponse{
		UserId:       userID,
		PostId:       posRef.ID,
		Content:      post.Content,
		Date:         today,
		IsPublic:     post.IsPublic,
		Category:     post.Category,
		CountComment: len(post.Comment),
		CountLike:    len(post.LikesRef),
		IsLike:       false,
	})
}

// service get my post
func GetMyPost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// declare all object variable
	postRes := models.PostResponse{}
	postsRes := []models.PostResponse{}
	post := models.Post{}

	// get all post by have a userID == userID
	allPostdocSnap, err := client.Collection("Post").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// convert data each post snap to post object
	for _, postDoc := range allPostdocSnap {
		mapstructure.Decode(postDoc.Data(), &post)

		postRes.UserId = post.UserID.ID
		postRes.PostId = postDoc.Ref.ID
		postRes.Content = post.Content
		postRes.Date = post.Date
		postRes.IsPublic = post.IsPublic
		postRes.CountLike = len(post.LikesRef)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.IsLike = checkIsLikePost(post.LikesRef, userID)

		postsRes = append(postsRes, postRes)
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, postsRes)
}

func DeleteMyPost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	postId := c.Param("post_id")

	// deelete post by user idd
	_, err := client.Collection("Post").Doc(postId).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	// response status 200ok to client
	c.JSON(200, gin.H{
		"message": "deleete success",
		"postId":  postId,
	})
}

func UpdatePost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")

	// declare all object to use
	post := models.Post{}
	postUpdate := models.Post{}

	// mapping data to object data
	if err := c.BindJSON(&postUpdate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	// get post By id and update for update some field
	postDoc, err := client.Collection("Post").Doc(postID).Get(ctx)
	mapstructure.Decode(postDoc.Data(), &post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	dataUopdate := models.Post{
		UserID:   post.UserID,
		Content:  postUpdate.Content,
		Date:     post.Date,
		IsPublic: post.IsPublic,
		Category: post.Category,
		Comment:  post.Comment,
		LikesRef: post.LikesRef,
	}
	_, err = client.Collection("Post").Doc(postID).Set(ctx, dataUopdate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, models.PostResponse{
		UserId:       dataUopdate.UserID.ID,
		PostId:       postID,
		Content:      dataUopdate.Content,
		Date:         dataUopdate.Date,
		IsPublic:     dataUopdate.IsPublic,
		Category:     dataUopdate.Category,
		CountComment: len(dataUopdate.Comment),
		CountLike:    len(dataUopdate.LikesRef),
		IsLike:       checkIsLikePost(dataUopdate.LikesRef, userID),
	})
}

// service get my post
func GetPostByKeyword(c *gin.Context) {
	// create a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	// userID := c.Request.Header.Get("id")
	keyword := c.Param("keyword")

	// Find posts containing the keyword
	query := client.Collection("Post").Where("Content", ">=", keyword).Where("Content", "<=", keyword+"\uf8ff")

	iter := query.Documents(ctx)
	defer iter.Stop()

	// Process query results and build response
	postsRes := make([]models.PostResponse, 0, 10)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		var post models.Post
		err = doc.DataTo(&post)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		postRes := models.PostResponse{
			PostId:       doc.Ref.ID,
			UserId:       post.UserID.ID,
			Content:      post.Content,
			Category:     post.Category,
			Date:         post.Date,
			CountLike:    len(post.LikesRef),
			CountComment: len(post.Comment),
		}
		postsRes = append(postsRes, postRes)
	}

	// return the results
	c.JSON(http.StatusOK, postsRes)

	// TODO: set to user data
}
