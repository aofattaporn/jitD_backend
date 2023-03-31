package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/mock"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// create set of test
func CreateSetTestStress(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// add to data base
	for index, _ := range mock.MockTestStress {
		_, err := client.Collection("TestStress").Doc(strconv.Itoa(index+1)).Set(ctx, mock.MockTestStress[index])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to add stress test",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "add set of test stress success",
	})

}

// create set of test
func CreateSetTestConsult(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// add to data base
	for index, _ := range mock.MockTestConsult {
		_, err := client.Collection("TestConsult").Doc(strconv.Itoa(index+1)).Set(ctx, mock.MockTestConsult[index])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to add stress test",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "add set of test stress success",
	})

}
