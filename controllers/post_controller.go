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

// ?check isLike from post like ref
func checkIsLikePost(userID string, postLike []*models.Like) bool {
	for _, l := range postLike {
		if l.UserRef.ID == userID {
			return true
		}
	}
	return false
}

// ?check isBookmark from post like ref
func checkIsBookMark(postID string, bookMark []*firestore.DocumentRef) bool {
	for _, dr := range bookMark {
		if dr.ID == postID {
			return true
		}
	}
	return false
}

// ?getUserData returns the user data for the given user ID
func getUserData(client *firestore.Client, userID string) (models.User, error) {
	userDoc, err := client.Collection("User").Doc(userID).Get(context.Background())
	if err != nil {
		return models.User{}, err
	}
	var userData models.User
	err = userDoc.DataTo(&userData)
	if err != nil {
		return models.User{}, err
	}
	return userData, nil
}

// ?getAllPosts returns all posts sorted by date in descending order
func getAllPosts(client *firestore.Client, userID string, userData models.User) ([]models.PostResponse, error) {
	var postResponses []models.PostResponse

	postsSnap, err := client.Collection("Post").OrderBy("Date", firestore.Desc).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range postsSnap {
		var post models.Post
		err := doc.DataTo(&post)
		if err != nil {
			return nil, err
		}

		// convertPostToResponse
		postResponse := convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponses = append(postResponses, postResponse)
	}
	return postResponses, nil
}

// ?convertPostToResponse converts a post to a post response object
func convertPostToResponse(post models.Post, userID string, userData models.User, postID string) models.PostResponse {
	postResponse := models.PostResponse{
		PostId:       postID,
		Content:      post.Content,
		Date:         post.Date,
		IsPublic:     post.IsPublic,
		Category:     post.Category,
		CountLike:    len(post.LikesRef),
		CountComment: len(post.Comment),
		IsLike:       checkIsLikePost(userID, post.LikesRef),
		IsBookmark:   checkIsBookMark(postID, userData.BookMark),
	}
	if post.UserID != nil {
		postResponse.UserId = post.UserID.ID
	}
	return postResponse
}

// ----------------------------------------------------------

func GetAllPostHomePage(c *gin.Context) {

	// declare data object
	postsResByDate := []models.PostResponse{}
	postResByLike := []models.PostResponse{}
	postRes := models.PostResponse{}

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// get all post by limit 10 result
	userData := models.User{}

	// get user data from userID
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user i",
		})
		return
	}
	mapstructure.Decode(userDoc.Data(), &userData)

	// get data homepage case1 - orderbydate
	postDateDoc, err := client.Collection("Post").Limit(10).OrderBy("Date", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get liked posts",
		})
		return
	}

	for _, doc := range postDateDoc {

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
		postRes.IsLike = checkIsLikePost(userID, post.LikesRef)
		postRes.IsBookmark = checkIsBookMark(doc.Ref.ID, userData.BookMark)
		postsResByDate = append(postsResByDate, postRes)
	}

	// get data homepage case2 - orderLike
	likeDocs, err := client.Collection("Like").Where("Date", ">", time.Now().AddDate(0, 0, -3)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get liked posts",
		})
		return
	}

	// Extract the post references from the liked documents
	postRefs := make([]*firestore.DocumentRef, len(likeDocs))
	for i, likeDoc := range likeDocs {
		like := models.Like{}
		if err := likeDoc.DataTo(&like); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to decode liked document",
			})
			return
		}
		postRefs[i] = like.PostRef
	}

	postLikeDocs, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get posts",
		})
		return
	}

	postRes = models.PostResponse{}

	fmt.Println(postLikeDocs)

	for _, v := range postLikeDocs {
		fmt.Println(v.Ref.ID)
	}

	for _, doc := range postLikeDocs {

		post := models.Post{}
		postRes.Category = []string{}

		if doc.Data() == nil {
			continue
		}
		mapstructure.Decode(doc.Data(), &post)

		postRes.UserId = post.UserID.ID
		postRes.PostId = doc.Ref.ID
		postRes.Content = post.Content
		postRes.Date = post.Date
		postRes.IsPublic = post.IsPublic
		postRes.CountLike = len(post.LikesRef)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.IsLike = checkIsLikePost(userID, post.LikesRef)
		postRes.IsBookmark = checkIsBookMark(doc.Ref.ID, userData.BookMark)
		postResByLike = append(postResByLike, postRes)
	}

	if len(postResByLike) < 10 {
		for _, post := range postsResByDate {
			postResByLike = append(postResByLike, post)
		}
	}

	// return json status code 200
	c.JSON(http.StatusOK, map[string]interface{}{
		"postDate":       postsResByDate,
		"postLike":       postResByLike,
		"postReccommend": postsResByDate,
	})
}

// *service get all post
func GetAllNewPost(c *gin.Context) {
	// Get user ID from request header
	userID := c.Request.Header.Get("id")

	// Get Firestore client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get user data
	userData, err := getUserData(client, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Get all posts
	postResponses, err := getAllPosts(client, userID, userData)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Return posts
	c.JSON(http.StatusOK, postResponses)
}

// *service post by like
func GetPostByLikeIndividual(c *gin.Context) {

	// declare data object
	postResponse := models.PostResponse{}
	postResponses := []models.PostResponse{}

	// get all id to use
	userID := c.Request.Header.Get("id")

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get the posts liked in the last 7 days
	likeDocs, err := client.Collection("Like").Where("Date", ">", time.Now().AddDate(0, 0, -7)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get liked posts",
		})
		return
	}

	// Extract the post references from the liked documents
	postRefs := make([]*firestore.DocumentRef, len(likeDocs))
	for i, likeDoc := range likeDocs {
		like := models.Like{}
		if err := likeDoc.DataTo(&like); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to decode liked document",
			})
			return
		}
		postRefs[i] = like.PostRef
	}

	// get user data
	userData, err := getUserData(client, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Get the posts corresponding to the extracted post references
	postDocs, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get posts",
		})
		return
	}

	// loop data snap and decode data to post respone
	for _, doc := range postDocs {

		post := models.Post{}
		mapstructure.Decode(doc.Data(), &post)
		postResponse = convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponses = append(postResponses, postResponse)

	}

	// return json status code 200
	c.JSON(http.StatusOK, postResponse)
}

// *service create post
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

	// TODO: update progress dialy quest
	UpdateProgressQuest(c, "PostQuest")

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
		IsBookmark:   false,
	})
}

// *service get my post
func GetMyPost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// declare all object variable
	postResponse := models.PostResponse{}
	postResponses := []models.PostResponse{}

	// get all post by have a userID == userID
	allPostdocSnap, err := client.Collection("Post").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// get user data
	userData, err := getUserData(client, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// convert data each post snap to post object
	for _, postDoc := range allPostdocSnap {
		post := models.Post{}
		mapstructure.Decode(postDoc.Data(), &post)
		postResponse = convertPostToResponse(post, userID, userData, postDoc.Ref.ID)
		postResponses = append(postResponses, postResponse)
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, postResponses)
}

// *service delete my post
func DeleteMyPost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")

	// deelete post by user idd
	_, err := client.Collection("Post").Doc(postID).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	// update the array field
	_, err = client.Collection("User").Doc(userID).Update(ctx, []firestore.Update{
		{
			Path:  "Bookmark",
			Value: firestore.ArrayRemove(client.Collection("Post").Doc(postID)),
		},
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update bookmarked posts",
			"error":   err.Error(),
		})
		return
	}

	// response status 200ok to client
	c.JSON(http.StatusOK, gin.H{
		"message": "deleete success",
		"postId":  postID,
	})
}

// *service update my post
func UpdatePost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	postID := c.Param("post_id")
	userID := c.Request.Header.Get("id")

	// get post by id
	postRef := client.Collection("Post").Doc(postID)
	postDoc, err := postRef.Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	// update post with new data
	var post models.Post
	if err := postDoc.DataTo(&post); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	var postUpdate models.Post
	if err := c.BindJSON(&postUpdate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}

	post.Content = postUpdate.Content
	post.IsPublic = postUpdate.IsPublic
	post.Category = postUpdate.Category

	if _, err := postRef.Set(ctx, post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update post"})
		return
	}

	// return updated post to frontend
	c.JSON(http.StatusOK, models.PostResponse{
		UserId:       post.UserID.ID,
		PostId:       postID,
		Content:      post.Content,
		Date:         post.Date,
		IsPublic:     post.IsPublic,
		Category:     post.Category,
		CountComment: len(post.Comment),
		CountLike:    len(post.LikesRef),
		IsLike:       checkIsLikePost(userID, post.LikesRef),
	})
}

// *Service to get posts by keyword
func GetPostByKeyword(c *gin.Context) {
	// Get Firestore client and context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get keyword from URL parameter
	keyword := c.Param("keyword")

	// declare all object variable
	postResponse := models.PostResponse{}
	postResponses := []models.PostResponse{}

	// Query posts containing the keyword
	query := client.Collection("Post").Where("Content", ">=", keyword).Where("Content", "<=", keyword+"\uf8ff")

	// Get post documents from the query
	postDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if len(postDocs) <= 0 {
		c.JSON(http.StatusOK, postResponses)
		return
	}

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// get user data
	userData, err := getUserData(client, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Map post documents to PostResponse objects
	for _, doc := range postDocs {
		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		postResponse = convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponses = append(postResponses, postResponse)
	}

	// Return the PostResponse objects
	c.JSON(http.StatusOK, postResponses)
}

// *get by by cateegory
func GetPostByCategorry(c *gin.Context) {

	// create a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// declare all object variable
	postResponse := models.PostResponse{}
	postResponses := []models.PostResponse{}

	category := c.Param("category")

	//geet all post by category
	postsDoc, err := client.Collection("Post").Where("Category", "array-contains", category).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to get all post by cateegory"})
		return
	}

	// get user data
	userData, err := getUserData(client, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// loop data snap and decode data to post respone
	for _, doc := range postsDoc {

		post := models.Post{}
		mapstructure.Decode(doc.Data(), &post)
		postResponse = convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponses = append(postResponses, postResponse)
	}

	// return json status code 200
	c.JSON(http.StatusOK, postResponses)
}
