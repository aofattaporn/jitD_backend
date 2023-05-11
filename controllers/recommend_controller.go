package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"github.com/mitchellh/mapstructure"
)

// ? All structure
type PostInteger struct {
	// Content             float64
	// Date                float64
	IsPublic            float64
	Comment             float64
	LikesCount          float64
	CountStudy          float64
	CountWork           float64
	CountMentalHealth   float64
	CountLifeProblem    float64
	CountRelationship   float64
	CountFamily         float64
	CountPhysicalHealth float64
}

type CountCategory struct {
	CateggoryName string
	Count         float64
}

type Data struct {
	PostId  string
	PostRef *firestore.DocumentRef
}

// ? normalization
// Example implementation of normalizing Content
func normalizeContent(content string) float64 {
	// Example implementation: length of content divided by a constant factor
	const factor = 10.0
	normalizedContent := float64(len(content)) / factor
	return normalizedContent
}

// Example implementation of normalizing Date
func normalizeDate(date time.Time) float64 {
	// Example implementation: convert the date to Unix timestamp and normalize it
	normalizedDate := float64(date.Unix())
	return normalizedDate
}

// Example implementation of normalizing IsPublic
func normalizeIsPublic(isPublic bool) float64 {
	// Example implementation: assign a value of 1.0 for true, 0.0 for false
	if isPublic {
		return 1.0
	}
	return 0.0
}

// Example implementation of normalizing Comment
func normalizeComment(commentCount int) float64 {
	// Example implementation: divide the comment count by a constant factor
	const factor = 100.0
	normalizedComment := float64(commentCount) / factor
	return normalizedComment
}

// Example implementation of normalizing LikesCount
func normalizeLikesCount(likesCount int) float64 {
	// Example implementation: divide the likes count by a constant factor
	const factor = 1000.0
	normalizedLikesCount := float64(likesCount) / factor
	return normalizedLikesCount
}

// Example implementation of normalizing the count attributes
func normalizeCount(catName string, category []string) float64 {
	for _, cat := range category {
		if strings.Contains(cat, catName) {
			return 1.0
		}
	}
	return 0.0
}

// ? avg function

func avgIsPublic(posts []models.PostResponse) float64 {
	countIsPublic := 0
	for _, pr := range posts {
		if pr.IsPublic {
			countIsPublic++
		}
	}

	if countIsPublic >= len(posts)/2 {
		return 1
	} else {
		return 0
	}
}

func avgCountComment(posts []models.PostResponse) float64 {
	countComment := 0
	for _, pr := range posts {
		countComment += pr.CountComment
	}

	return normalizeComment(countComment / len(posts))
}

func avgCountLike(posts []models.PostResponse) float64 {
	countLike := 0
	for _, pr := range posts {
		countLike += pr.CountLike
	}

	return normalizeLikesCount(countLike / len(posts))
}

func avgCountCategory(posts []models.PostResponse) []string {
	categoryCounts := make(map[string]int)

	// Count the occurrences of each category in the posts slice
	for _, post := range posts {
		for _, category := range post.Category {
			categoryCounts[category]++
		}
	}

	// Create a slice of CateTest objects to hold the category counts
	type CateTest struct {
		CatName string `json:"catName"`
		Count   int    `json:"count"`
	}
	var categories []CateTest

	// Convert the categoryCounts map to a slice
	for cat, count := range categoryCounts {
		categories = append(categories, CateTest{
			CatName: cat,
			Count:   count,
		})
	}

	// Sort the categories slice by count in descending order
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Count > categories[j].Count
	})

	// Extract the top 3 categories with the highest counts
	topCategories := make([]string, 0)
	for i := 0; i < len(categories) && i < 3; i++ {
		topCategories = append(topCategories, categories[i].CatName)
	}

	return topCategories
}

func convertAveragePost(posts []models.PostResponse) []float64 {

	categories := avgCountCategory(posts)

	postFloat := PostInteger{
		IsPublic:            avgIsPublic(posts),
		Comment:             avgCountComment(posts),
		LikesCount:          avgCountLike(posts),
		CountStudy:          normalizeCount("การเรียน", categories),
		CountWork:           normalizeCount("การงาน", categories),
		CountMentalHealth:   normalizeCount("สุขภาพจิต", categories),
		CountLifeProblem:    normalizeCount("ปัญหาชีวิต", categories),
		CountRelationship:   normalizeCount("ความสัมพันธ์", categories),
		CountFamily:         normalizeCount("ครอบครัว", categories),
		CountPhysicalHealth: normalizeCount("สุขภาพร่างกาย", categories),
	}

	return []float64{
		postFloat.IsPublic,
		postFloat.Comment,
		postFloat.LikesCount,
		postFloat.CountStudy,
		postFloat.CountWork,
		postFloat.CountMentalHealth,
		postFloat.CountLifeProblem,
		postFloat.CountRelationship,
		postFloat.CountFamily,
		postFloat.CountPhysicalHealth,
	}
}

// ? main function for reccommend
func tranformPostToPostFloat64(post models.PostResponse) []float64 {

	postFloat := PostInteger{
		IsPublic:            normalizeIsPublic(post.IsPublic),
		Comment:             normalizeComment(post.CountComment),
		LikesCount:          normalizeLikesCount(post.CountLike),
		CountStudy:          normalizeCount("การเรียน", post.Category),
		CountWork:           normalizeCount("การงาน", post.Category),
		CountMentalHealth:   normalizeCount("สุขภาพจิต", post.Category),
		CountLifeProblem:    normalizeCount("ปัญหาชีวิต", post.Category),
		CountRelationship:   normalizeCount("ความสัมพันธ์", post.Category),
		CountFamily:         normalizeCount("ครอบครัว", post.Category),
		CountPhysicalHealth: normalizeCount("สุขภาพร่างกาย", post.Category),
	}

	return []float64{
		postFloat.IsPublic,
		postFloat.Comment,
		postFloat.LikesCount,
		postFloat.CountStudy,
		postFloat.CountWork,
		postFloat.CountMentalHealth,
		postFloat.CountLifeProblem,
		postFloat.CountRelationship,
		postFloat.CountFamily,
		postFloat.CountPhysicalHealth,
	}
}

func getPostByLike(client firestore.Client, c *gin.Context, ctx context.Context, userID string, userData models.User) []float64 {

	// TODO: get My Post
	// Get the posts liked in the last 7 days
	likeDocs, err := client.Collection("Like").Where("Date", ">", time.Now().AddDate(0, 0, -7)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get liked posts",
		})
		return nil
	}

	if len(likeDocs) == 0 {
		return []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
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
			return nil
		}
		postRefs[i] = like.PostRef
	}

	// Get the posts corresponding to the extracted post references
	postDocs, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get posts",
		})
		return nil
	}

	postResponses := []models.PostResponse{}

	// loop data snap and decode data to post respone
	for _, doc := range postDocs {

		post := models.Post{}
		mapstructure.Decode(doc.Data(), &post)
		postResponse := convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponses = append(postResponses, postResponse)

	}

	// convert to post aveage
	return convertAveragePost(postResponses)

}

func RecommendPost(c *gin.Context) {

	// Get user ID from request header
	userID := c.Request.Header.Get("id")

	// Get Firestore client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get user data
	userData, err := getUserData(client, userID)
	// post := models.Post{}
	postRefs := []*firestore.DocumentRef{}
	postResponsesFinal := []models.PostResponse{}

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

	// ? initail data point
	tree := kdtree.New(nil)

	for _, post := range postResponses {
		// ?Insert
		tree.Insert(points.NewPoint(tranformPostToPostFloat64(post), Data{PostId: post.PostId, PostRef: client.Collection("Post").Doc(post.PostId)}))
	}

	// Balance
	tree.Balance()
	// Iterate over points and get Data names
	// fmt.Println("My Point ***************************")
	// fmt.Println(getPostByLike(*client, c, ctx, userID, userData))
	// fmt.Println("***************************")

	allPoints := tree.KNN(&points.Point{Coordinates: getPostByLike(*client, c, ctx, userID, userData)}, 100)
	for _, point := range allPoints {
		data := point.(*points.Point).Data.(Data)
		postRefs = append(postRefs, data.PostRef)

	}

	postDoc, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, doc := range postDoc {
		var post models.Post
		err := doc.DataTo(&post)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// convertPostToResponse
		postResponse := convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponsesFinal = append(postResponsesFinal, postResponse)
	}
	// return json status code 200
	c.JSON(http.StatusOK, postResponsesFinal)

}

func RecommendPostIndividule(c *gin.Context) []models.PostResponse {

	// Get user ID from request header
	userID := c.Request.Header.Get("id")

	// Get Firestore client
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get user data
	userData, err := getUserData(client, userID)
	// post := models.Post{}
	postRefs := []*firestore.DocumentRef{}
	postResponsesFinal := []models.PostResponse{}

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return []models.PostResponse{}
	}

	// Get all posts
	postResponses, err := getAllPosts(client, userID, userData)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return []models.PostResponse{}
	}

	// ? initail data point
	tree := kdtree.New(nil)

	for _, post := range postResponses {
		// ?Insert
		tree.Insert(points.NewPoint(tranformPostToPostFloat64(post), Data{PostId: post.PostId, PostRef: client.Collection("Post").Doc(post.PostId)}))
	}

	// Balance
	tree.Balance()
	// Iterate over points and get Data names
	// fmt.Println("My Point ***************************")
	// fmt.Println(getPostByLike(*client, c, ctx, userID, userData))
	// fmt.Println("***************************")

	allPoints := tree.KNN(&points.Point{Coordinates: getPostByLike(*client, c, ctx, userID, userData)}, 100)
	for _, point := range allPoints {
		data := point.(*points.Point).Data.(Data)
		postRefs = append(postRefs, data.PostRef)

	}

	postDoc, err := client.GetAll(ctx, postRefs)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return []models.PostResponse{}
	}

	for _, doc := range postDoc {
		var post models.Post
		err := doc.DataTo(&post)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return []models.PostResponse{}
		}

		// convertPostToResponse
		postResponse := convertPostToResponse(post, userID, userData, doc.Ref.ID)
		postResponsesFinal = append(postResponsesFinal, postResponse)
	}
	// return json status code 200
	return postResponsesFinal

}
