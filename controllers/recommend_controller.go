package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type Post struct {
	ID      int
	Content string
	UserID  int
}

func saveCountCategory(client *firestore.Client, categorys []string, userID string, categoryID string, ctx context.Context) {

	var err error
	// TODO : check if exist

	// Check if the document exists
	docRef := client.Collection("User").Doc(categoryID)
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		// Handle error
		fmt.Println("Failed to fetch document:", err)
		return
	}

	if !docSnapshot.Exists() {
		// Create the document if it doesn't exist
		_, err = docRef.Set(ctx, models.CategoryReccommend{
			UserID:              client.Collection("User").Doc(userID),
			Date:                time.Now().UTC(),
			CountStudy:          0,
			CountWork:           0,
			CountMentalHealth:   0,
			CountLifeProblem:    0,
			CountRelationship:   0,
			CountFamily:         0,
			CountPhysicalHealth: 0,
		})
		if err != nil {
			// Handle error
			fmt.Println("Failed to create document:", err)
			return
		}

		// update user feild
		_, err = client.Collection("User").Doc(userID).Update(ctx, []firestore.Update{
			{
				Path:  "CategoryID",
				Value: docRef.ID,
			},
		})
		if err != nil {
			// Handle error
			fmt.Println("Failed to create document:", err)
			return
		}
	}

	// TODO : update post
	for _, cat := range categorys {
		if cat == "การเรียน" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountStudy",
					Value: firestore.Increment(1),
				},
			})
		} else if cat == "การงาน" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountWork",
					Value: firestore.Increment(1),
				},
			})

		} else if cat == "สุขภาพจิต" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountMentalHealth",
					Value: firestore.Increment(1),
				},
			})

		} else if cat == "ปัญหาชีวิต" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountLifeProblem",
					Value: firestore.Increment(1),
				},
			})

		} else if cat == "ความสัมพันธ์" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountRelationship",
					Value: firestore.Increment(1),
				},
			})

		} else if cat == "ครอบครัว" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountFamily",
					Value: firestore.Increment(1),
				},
			})

		} else if cat == "สุขภาพร่างกาย" {
			_, err = client.Collection("User").Doc(categoryID).Update(ctx, []firestore.Update{
				{
					Path:  "CountPhysicalHealth",
					Value: firestore.Increment(1),
				},
			})

		}
	}

	fmt.Println(err)

}

func ReccomendPost(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// TODO : Get like informaton and anlysis
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

	// postLikeDocs, err := client.GetAll(ctx, postRefs)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Failed to get posts",
	// 	})
	// 	return
	// }

	// TODO : calculate on all post

	// TODO :  1 - tranfer all post to spot

	// TODO :  2 - callculate all spot

	// // Generate some example posts
	// var posts []*Post
	// for i := 1; i <= 10; i++ {
	// 	post := &Post{ID: i, Content: fmt.Sprintf("Post %d content", i), UserID: i % 3}
	// 	posts = append(posts, post)
	// }

	// // // Build a k-d tree from the feature vectors of the posts
	// var tree *kdtree.KDTree

	// for i, _ := range posts {
	// 	fmt.Print(i)
	// 	// tree.Insert(points.NewPoint([]float64{12, 4, 6}, i))
	// }

	// // // Find the k-nearest neighbors of a new post
	// pointNewPost := points.NewPoint([]float64{12, 4, 6}, 0)

	// k := 3
	// neighbors := tree.KNN(pointNewPost, k)

	// // // Recommend the posts based on the k-nearest neighbors
	// var recommendedPosts []*Post
	// fmt.Println("==================================")
	// for _, neighbor := range neighbors {
	// 	post := neighbor.Dimensions()
	// 	fmt.Println(post)
	// 	// recommendedPosts = append(recommendedPosts, post)
	// }
	// fmt.Println("==================================")

	// fmt.Printf("Recommended posts: %v\n", recommendedPosts)
}
