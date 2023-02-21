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
	user := models.User{}
	post := models.Post{}
	user_id := c.Request.Header.Get("id")

	// Call BindJSON to bind the received JSON body
	if err := c.BindJSON(&post); err != nil {
		log.Fatalln(err)
		return
	}

	// get userRef
	userRef, user_err := client.Collection("User").Doc(user_id).Get(ctx)
	if user_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cant to find userid",
		})
		return
	}
	mapstructure.Decode(userRef.Data(), &user)

	//----------- adding post data to Posts ---------------
	post.Date = time.Now().UTC()

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		post.UserID = userRef.Ref
		post.Date = time.Now().UTC()
		post.Comment = []*models.Comment2{}
		post.LikesRef = []*models.Like{}
		postRef, _, post_err := client.Collection("Post").Add(ctx, post)
		if post_err != nil {
			return post_err
		}

		// update post feild [ useer collection ]
		user.Posts = append(user.Posts, postRef)
		_, user_err_update := client.Collection("User").Doc(user_id).Set(ctx, user)
		if user_err_update != nil {
			return user_err_update
		}
		return nil
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
		postRes.CountLike = len(post.LikesRef)
		postRes.Category = post.Category

		postsRes = append(postsRes, postRes)
	}

	c.JSON(http.StatusOK, postsRes)
}

func DeleteMyPost(c *gin.Context) {
	post_id := c.Param("post_id")
	user_id := c.Request.Header.Get("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	if err := deletePostAndRemoveReference(ctx, client, post_id, user_id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleete success"})
}

func UpdatePost(c *gin.Context) {
	post_id := c.Param("post_id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	fmt.Printf("time.Now().UTC(): %v\n", time.Now().Format(time.RFC3339))

	//----------- adding post data to Posts ---------------
	currentTime := time.Now().Format(time.RFC3339)
	currentDateTime, err := time.Parse(time.RFC3339, currentTime)

	// Get the post document
	postDoc, err := client.Collection("Post").Doc(post_id).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	// Convert the post document to a Post struct
	var post models.Post
	if err := postDoc.DataTo(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Get the updated post data from the request body
	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update only the specified fields and retain the original values for the rest
	if updatedPost.UserID != nil {
		post.UserID = updatedPost.UserID
	}
	if updatedPost.Content != "" {
		post.Content = updatedPost.Content
	}
	if !updatedPost.Date.IsZero() {
		post.Date = updatedPost.Date
	}
	post.IsPublic = updatedPost.IsPublic
	post.Category = updatedPost.Category
	post.Date = currentDateTime
	post.LikesRef = updatedPost.LikesRef
	if len(post.Comment) == 0 {
		post.Comment = []*models.Comment2{}
	}
	if len(post.LikesRef) == 0 {
		post.LikesRef = []*models.Like{}
	}

	// Update the post document in the database
	if _, err := client.Collection("Post").Doc(post_id).Set(c.Request.Context(), post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func deletePostAndRemoveReference(ctx context.Context, fsClient *firestore.Client, postID, userID string) error {
	err := fsClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get the post document
		postDoc, err := tx.Get(fsClient.Doc("Post/" + postID))
		if err != nil {
			return err
		}

		// Get the user document
		userDoc, err := tx.Get(fsClient.Doc("User/" + userID))
		if err != nil {
			return err
		}

		// Get the references from the user document as a map
		referencesMap := make(map[string]interface{})
		if err := userDoc.DataTo(&referencesMap); err != nil {
			return err
		}

		references, ok := referencesMap["Posts"].([]interface{})
		if !ok {
			return fmt.Errorf("References field is not an array")
		}

		// Convert the references array to []*firestore.DocumentRef
		refs := make([]*firestore.DocumentRef, len(references))
		for i, ref := range references {
			refs[i] = ref.(*firestore.DocumentRef)
		}

		// Find the index of the reference to the post
		var index int
		for i, ref := range refs {
			if ref.ID == postID {
				index = i
				break
			}
		}

		// Remove the reference from the references array
		refs = append(refs[:index], refs[index+1:]...)

		// Convert the references array back to an array of interface{}
		newReferences := make([]interface{}, len(refs))
		for i, ref := range refs {
			newReferences[i] = ref
		}

		// Update the user document with the modified references array
		referencesMap["Posts"] = newReferences
		if err := tx.Set(userDoc.Ref, referencesMap, firestore.MergeAll); err != nil {
			return err
		}

		// Delete the post document
		if err := tx.Delete(postDoc.Ref); err != nil {
			return err
		}

		return nil
	})
	return err
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
