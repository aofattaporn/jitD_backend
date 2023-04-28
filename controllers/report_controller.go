package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// service get my post
func GetAllReport(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all object variable
	reportRes := models.ReportResponse{}
	reportsRes := []models.ReportResponse{}

	// get all post by have a userID == userID
	allReportSnap, err := client.Collection("Report").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get all report",
		})
	}

	// convert data each post snap to post object
	for _, reportDoc := range allReportSnap {
		mapstructure.Decode(reportDoc.Data(), &reportRes)
		reportRes.ReportID = reportDoc.Ref.ID
		reportsRes = append(reportsRes, reportRes)
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, reportsRes)
}

// service get my post
func AddReport(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// declare all object variable
	report := models.Report{}
	userID := c.Request.Header.Get("id")

	// mappig data body to object post
	if err := c.BindJSON(&report); err != nil {
		log.Fatalln(err)
		return
	}

	// get all post by have a userID == userID
	_, _, err := client.Collection("Report").Add(ctx, models.Report{
		UserID:     userID,
		ReportName: report.ReportName,
		ReportType: report.ReportType,
		Date:       time.Now().UTC(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant get all report",
		})
	}

	// return data to frontend status 200
	c.JSON(http.StatusOK, gin.H{
		"message": "add report success",
	})
}

// DeletePostReportByPostID is a function that handles the deletion of a report by post ID
func DeleteReport(c *gin.Context) {
	// Declare an instance of firestore and a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get the post ID from the URL parameter
	reportID := c.Param("report_id")

	// Delete the report from the "Report" collection in Firestore with a matching post ID
	_, err := client.Collection("Report").Doc(reportID).Delete(ctx)

	// Handle errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to delete report",
		})
		return
	}

	// Return a success message to the frontend with a 200 status code
	c.JSON(http.StatusOK, gin.H{
		"message": "Report deleted successfully",
	})
}

// UpdateReport is a function that handles the updating of a report
func UpdateReport(c *gin.Context) {
	// Declare an instance of firestore and a context
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// Get the report ID from the URL parameter
	reportID := c.Param("id")

	// Declare a variable to hold the updated report data
	updatedReport := models.Report{}

	// Map the data from the request body to the updatedReport object
	if err := c.BindJSON(&updatedReport); err != nil {
		log.Fatalln(err)
		return
	}

	// Update the report in the "Report" collection in Firestore with the updated data
	_, err := client.Collection("Report").Doc(reportID).Set(ctx, updatedReport)

	// Handle errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update report",
		})
		return
	}

	// Return a success message to the frontend with a 200 status code
	c.JSON(http.StatusOK, gin.H{
		"message": "Report updated successfully",
	})
}
