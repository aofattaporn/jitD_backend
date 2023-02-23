package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func UpdateProgressQuest(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// delcrae object data
	user := models.User{}

	// get a user
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to find by user id"})
		return
	}

	// mapping data
	mapstructure.Decode(userDoc.Data(), &user)

	quest := models.DailyQuestProgress{}
	// Decode user's daily quests from the document
	if dailyQuests, err := userDoc.DataAt("DailyQuests"); err == nil {
		if err := mapstructure.Decode(dailyQuests, &quest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to decode daily quests",
			})
			return
		}
	}

	if quest.Quests != nil {
		for index, _ := range quest.Quests {
			if quest.Quests[index].QuestName == "LikeQuest" {
				// updat progess quest
				if quest.Quests[index].Completed {
					return
				} else {
					quest.Quests[index].Progress += 1

					// set new reward
					if quest.Quests[index].IsGetPoint {
						quest.Quests[index].Reward = quest.Quests[index].Progress * 5
					} else {
						quest.Quests[index].Reward = 5
					}

					// set completed
					if quest.Quests[index].Progress == quest.Quests[index].MaxProgress {
						quest.Quests[index].Completed = true
					}
					quest.Quests[index].IsGetPoint = true

				}
				break
			}
		}
	}

	mapstructure.Decode(quest, &user.DailyQuests)

	fmt.Printf("user.Point: %v\n", user.Point)

	_, err = client.Collection("User").Doc(userID).Set(ctx, user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "update success"})
		return
	}
}

func GetPointFromQuest(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// delcrae object data
	user := models.User{}

	// get a user
	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to find by user id"})
		return
	}

	// mapping data
	mapstructure.Decode(userDoc.Data(), &user)

	pointGet := 0
	// Decode user's daily quests from the document
	if dailyQuests, err := userDoc.DataAt("DailyQuests"); err == nil {
		if err := mapstructure.Decode(dailyQuests, &user.DailyQuests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to decode daily quests",
			})
			return
		}
	}

	// get point by name quest
	for index, _ := range user.DailyQuests.Quests {
		if user.DailyQuests.Quests[index].QuestName == "LikeQuest" && user.DailyQuests.Quests[index].IsGetPoint == true {
			pointGet = user.DailyQuests.Quests[index].Reward
			user.DailyQuests.Quests[index].IsGetPoint = false
			break
		}
	}

	// set new point and hp
	user.Point += pointGet

	// set to db
	_, err = client.Collection("User").Doc(userID).Set(ctx, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
		return
	}

}
