package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// get stress test
func GetTestStress(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve test data
	var testStress models.Test
	settestStress := []models.Test{}

	// get all doument of test stress
	testData, err := client.Collection("TestStress").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to add stress test",
		})
		return
	}

	//
	for i, _ := range testData {
		testStress = models.Test{}
		mapstructure.Decode(testData[i].Data(), &testStress)
		settestStress = append(settestStress, testStress)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "add set of test stress success",
	})
}

// create set of test
func CreateSetTestStress(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// create set of test
	mocktestStress := []models.Test{
		{
			Number:       1,
			QuestionText: "question 1",
			Text:         [4]string{"1", "2", "3", "4"},
		},
		{
			Number:       2,
			QuestionText: "question 2",
			Text:         [4]string{"1", "2", "3", "4"},
		},
		{
			Number:       3,
			QuestionText: "question 3",
			Text:         [4]string{"1", "2", "3", "4"},
		},
	}

	// add to data base
	for index, _ := range mocktestStress {
		_, _, err := client.Collection("TestStress").Add(ctx, mocktestStress[index])
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
