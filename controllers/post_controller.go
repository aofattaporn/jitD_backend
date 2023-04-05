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

// check isLike from post like ref
func checkIsLikePost(postLike []*models.Like, userID string) bool {
	for _, l := range postLike {
		if l.UserRef.ID == userID {
			return true
		}
	}
	return false
}

func checkIsBookMark(postID string, bookMark []*firestore.DocumentRef) bool {

	for _, dr := range bookMark {
		if dr.ID == postID {
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

	userData := models.User{}

	// get all id to use
	userID := c.Request.Header.Get("id")

	// get all post by limit 10 result
	allDocSnap, err := client.Collection("Post").Limit(10).Documents(ctx).GetAll()
	if err != nil {
		return
	}

	// get user data
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		return
	}

	mapstructure.Decode(userDoc.Data(), &userData)

	var user models.User
	// Decode user data from Firestore document
	bookmarkRefs, ok := userDoc.Data()["Bookmark"].([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve bookmarked post references. Bookmark field not found in the user document.",
		})
		return
	}

	mapstructure.Decode(userDoc.Data(), &user)

	if len(bookmarkRefs) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully retrieved bookmarked posts",
			"data":    []models.PostResponse{},
		})
		return
	}

	// Get bookmarked posts from Firestore
	mapstructure.Decode(bookmarkRefs, &user.BookMark)

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
		postRes.IsBookmark = checkIsBookMark(doc.Ref.ID, userData.BookMark)
		posts = append(posts, postRes)
	}

	// return json status code 200
	c.JSON(http.StatusOK, posts)
}

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
		postRes.IsLike = checkIsLikePost(post.LikesRef, userID)
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
	for _, doc := range postLikeDocs {

		post := models.Post{}
		postRes.Category = []string{}
		mapstructure.Decode(doc.Data(), &post)

		// map data to seend to fronend
		postRes.UserId = doc.Ref.ID
		postRes.PostId = doc.Ref.ID
		postRes.Content = post.Content
		postRes.Date = post.Date
		postRes.IsPublic = post.IsPublic
		postRes.CountLike = len(post.LikesRef)
		postRes.CountComment = len(post.Comment)
		postRes.Category = post.Category
		postRes.IsLike = checkIsLikePost(post.LikesRef, userID)
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

// service post by like
func GetPostByLikeIndividual(c *gin.Context) {

	// declare data object
	posts := []models.PostResponse{}
	postRes := models.PostResponse{}
	userData := models.User{}

	// get all id to use
	userID := c.Request.Header.Get("id")

	//
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
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		return
	}

	mapstructure.Decode(userDoc.Data(), &userData)

	// Get the posts corresponding to the extracted post references
	postDocs, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get posts",
		})
		return
	}
	mapstructure.Decode(userDoc.Data(), &userData)

	// loop data snap and decode data to post respone
	for _, doc := range postDocs {

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
		postRes.IsBookmark = checkIsBookMark(doc.Ref.ID, userData.BookMark)
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
	userData := models.User{}

	// get all post by have a userID == userID
	allPostdocSnap, err := client.Collection("Post").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get infomation",
		})
	}

	// get user data
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		return
	}
	mapstructure.Decode(userDoc.Data(), &userData)

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
		postRes.IsBookmark = checkIsBookMark(postDoc.Ref.ID, userData.BookMark)

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
		IsPublic: postUpdate.IsPublic,
		Category: postUpdate.Category,
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

// get by by cateegory
func GetPostByCategorry(c *gin.Context) {

	// create a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all id to use
	userID := c.Request.Header.Get("id")

	// declare object to usee
	category := c.Param("category")
	postsRes := []models.PostResponse{}
	postRes := models.PostResponse{}

	//geet all post by category
	postsDoc, err := client.Collection("Post").Where("Category", "array_contains", category).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to get all post by cateegory"})
		return
	}

	// loop data snap and decode data to post respone
	for _, doc := range postsDoc {

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

		postsRes = append(postsRes, postRes)
	}

	// return json status code 200
	c.JSON(http.StatusOK, postsRes)

}
