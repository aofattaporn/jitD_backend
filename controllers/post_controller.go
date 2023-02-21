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
	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
)

// service create post
func CreatePost(c *gin.Context) {

	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post := models.Post{}
	userID := c.Request.Header.Get("id")

	// Call BindJSON to bind the received JSON body
	if err := c.BindJSON(&post); err != nil {
		log.Fatalln(err)
		return
	}

	_, _, err := client.Collection("Post").Add(ctx, models.Post{
		UserID:   client.Collection("User").Doc(userID),
		Content:  post.Content,
		Date:     time.Now().UTC(),
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
	})
}

// service get all post
func GetAllPost(c *gin.Context) {
	posts := []models.PostResponse{}

	// post := models.Post{}
	postRes := models.PostResponse{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	snap, err := client.Collection("Post").Limit(10).Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, element := range snap {

		post := models.Post{}
		postRes.Category = []string{}
		mapstructure.Decode(element.Data(), &post)
		mapstructure.Decode(post, &postRes)

		postRes.UserId = post.UserID.ID
		postRes.PostId = element.Ref.ID
		postRes.CountLike = len(post.LikesRef)
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
	userID := c.Request.Header.Get("id")

	postRes := models.PostResponse{}
	postsRes := []models.PostResponse{}
	post := models.Post{}

	doc, err := client.Collection("Post").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	for _, element := range doc {
		mapstructure.Decode(element.Data(), &post)
		mapstructure.Decode(post, &postRes)
		postRes.UserId = post.UserID.ID
		postRes.PostId = element.Ref.ID
		postRes.Date = post.Date
		postRes.CountComment = len(post.Comment)
		postRes.CountLike = len(post.LikesRef)
		postRes.Category = post.Category

		postsRes = append(postsRes, postRes)
	}

	c.JSON(http.StatusOK, postsRes)
}

func DeleteMyPost(c *gin.Context) {
	postId := c.Param("post_id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	_, err := client.Collection("Post").Doc(postId).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	c.JSON(200, gin.H{"message": "deleete success"})
}

func UpdatePost(c *gin.Context) {
	postId := c.Param("post_id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	post := models.Post{}
	postUpdate := models.Post{}

	// Call BindJSON to bind the received JSON body
	if err := c.BindJSON(&postUpdate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	postDoc, err := client.Collection("Post").Doc(postId).Get(ctx)
	mapstructure.Decode(postDoc.Data(), &post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	_, err = client.Collection("Post").Doc(postId).Set(ctx, models.Post{
		UserID:   post.UserID,
		Content:  postUpdate.Content,
		Date:     post.Date,
		IsPublic: post.IsPublic,
		Category: post.Category,
		Comment:  post.Comment,
		LikesRef: post.LikesRef,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	c.JSON(200, gin.H{"message": "update post success"})
}

// service get my post
func GetPostByKeyword(c *gin.Context) {
	// create a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	userID := c.Request.Header.Get("id")
	keyword := c.Param("keyword")

	// Find user document and update search history
	userRef := client.Collection("User").Doc(userID)
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		userSnap, err := tx.Get(userRef)
		if err != nil {
			return err
		}
		userData := userSnap.Data()
		userHistory := userData["HistorySearch"].([]interface{})
		userHistory = append(userHistory, keyword)
		if len(userHistory) > 5 {
			userHistory = userHistory[len(userHistory)-5:]
		}
		userData["HistorySearch"] = userHistory
		return tx.Set(userRef, userData)
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Find posts containing the keyword
	query := client.Collection("Post").Where("Content", ">=", keyword).Limit(10)
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
}
