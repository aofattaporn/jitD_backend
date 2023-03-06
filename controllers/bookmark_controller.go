package controllers

import (
	"context"
	configs "jitD/configs"

	"github.com/gin-gonic/gin"
)

// LikePost creates a user like on a posts
func AddBookmark(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare id to use in this function
	userID := c.Request.Header.Get("id")
	postID := c.Param("post_id")

}
