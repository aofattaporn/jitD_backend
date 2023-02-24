package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func GetMyQuest(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// delcrae object data
	quest := models.DailyQuestProgress{}
	today := time.Now().UTC().Day()

	fmt.Println(uuid.NewUUID())
	fmt.Println(uuid.NewRandom())

	// check diary queest
	questDoc, err := client.Collection("DialyQuest").Where("UserID", "==", client.Collection("User").Doc(userID)).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
		return
	}

	// mapping data to object quest
	if len(questDoc) < 1 {

		// create quest
		setQuest := []*models.Quest{}
		questName := []string{"PostQuest", "CommentQuest", "LikeQuest"}
		for _, name := range questName {
			setQuest = append(setQuest, &models.Quest{
				QuestName:      name,
				Progress:       0,
				MaxProgress:    3,
				Reward:         5,
				IsGetPoint:     false,
				Completed:      false,
				LastCompletion: time.Now().UTC(),
			})
		}

		// set to diary quest
		diaryQuest := models.DailyQuestProgress{
			UserID:    client.Collection("User").Doc(userID),
			QuestDate: time.Now().UTC(),
			Quests:    setQuest,
		}

		// set data and return
		_, err := client.Collection("DialyQuest").Doc(string(uuid.NewString())).Set(ctx, diaryQuest)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
			return
		}

		c.JSON(http.StatusOK, diaryQuest)
		return

		// case quest equal than 0 (no data)
	} else if len(questDoc) == 1 {

		// case quest equal than 1
		questID := questDoc[0].Ref.ID
		mapstructure.Decode(questDoc[0].Data(), &quest)

		// check data quest
		if quest.QuestDate.Day() != today {

			// create quest
			setQuest := []*models.Quest{}
			questName := []string{"PostQuest", "CommentQuest", "LikeQuest"}
			for _, name := range questName {
				setQuest = append(setQuest, &models.Quest{
					QuestName:      name,
					Progress:       0,
					MaxProgress:    3,
					Reward:         5,
					IsGetPoint:     false,
					Completed:      false,
					LastCompletion: time.Now().UTC(),
				})
			}

			// set to diary quest
			diaryQuest := models.DailyQuestProgress{
				UserID:    client.Collection("User").Doc(userID),
				QuestDate: time.Now().UTC(),
				Quests:    setQuest,
			}

			// set data and return
			_, err := client.Collection("DialyQuest").Doc(questID).Set(ctx, diaryQuest)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
				return
			}

			c.JSON(http.StatusNotFound, diaryQuest)
			return
		} else {
			c.JSON(http.StatusOK, quest)
			return
		}

	} else {

		// case quest more than 1
		c.JSON(http.StatusBadRequest, gin.H{"message": "this userID have quest more than 1 pls contect BE for rresolve this ploblem"})
		return
	}

}

func UpdateProgressQuest(c *gin.Context, questName string) {

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
	if dailyQuests, err := userDoc.DataAt(questName); err == nil {
		if err := mapstructure.Decode(dailyQuests, &quest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to decode daily quests",
			})
			return
		}
	}

	if quest.Quests != nil {
		for index, _ := range quest.Quests {
			if quest.Quests[index].QuestName == questName {
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

	// mapstructure.Decode(quest, &user.DailyQuests)

	fmt.Printf("user.Point: %v\n", user.Point)

	_, err = client.Collection("User").Doc(userID).Set(ctx, user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "update success"})
		return
	}
}

// func GetPointFromQuest(c *gin.Context) {

// 	// declare instance of firestore
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)

// 	// get all id to use
// 	userID := c.Request.Header.Get("id")

// 	// delcrae object data
// 	user := models.User{}

// 	// get a user
// 	userDoc, err := client.Collection("User").Doc(userID).Get(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "cant to find by user id"})
// 		return
// 	}

// 	// mapping data
// 	mapstructure.Decode(userDoc.Data(), &user)

// 	pointGet := 0
// 	// Decode user's daily quests from the document
// 	if dailyQuests, err := userDoc.DataAt("DailyQuests"); err == nil {
// 		if err := mapstructure.Decode(dailyQuests, &user.DailyQuests); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Failed to decode daily quests",
// 			})
// 			return
// 		}
// 	}

// 	// get point by name quest
// 	for index, _ := range user.DailyQuests.Quests {
// 		if user.DailyQuests.Quests[index].QuestName == "LikeQuest" && user.DailyQuests.Quests[index].IsGetPoint == true {
// 			pointGet = user.DailyQuests.Quests[index].Reward
// 			user.DailyQuests.Quests[index].IsGetPoint = false
// 			break
// 		}
// 	}

// 	// set new point and hp
// 	user.Point += pointGet

// 	// set to db
// 	_, err = client.Collection("User").Doc(userID).Set(ctx, user)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
// 		return
// 	}

// }
