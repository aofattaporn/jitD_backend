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
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
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

			c.JSON(http.StatusOK, diaryQuest)
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

func UpdateProgressQuest(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// get user to object
	dialyQuest := models.DailyQuestProgress{}

	// get a user
	userDoc, err := client.Collection("DialyQuest").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to find by user id"})
		return
	}

	if len(userDoc) > 1 || len(userDoc) == 0 {
		fmt.Println(len(userDoc))
		fmt.Println(client.Collection("User").Doc(userID))

		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to get your quest"})
		return
	}

	// mapping data
	mapstructure.Decode(userDoc[0].Data(), &dialyQuest)

	if dialyQuest.Quests != nil {
		for index, _ := range dialyQuest.Quests {
			if dialyQuest.Quests[index].QuestName == "LikeQuest" {
				// updat progess quest
				if dialyQuest.Quests[index].Completed {
					return
				} else {
					dialyQuest.Quests[index].Progress += 1

					// set new reward
					if dialyQuest.Quests[index].IsGetPoint {
						dialyQuest.Quests[index].Reward = dialyQuest.Quests[index].Progress * 5
					} else {
						dialyQuest.Quests[index].Reward = 5
					}

					// set completed
					if dialyQuest.Quests[index].Progress == dialyQuest.Quests[index].MaxProgress {
						dialyQuest.Quests[index].Completed = true
					}
					dialyQuest.Quests[index].IsGetPoint = true

				}
				break
			}
		}
	}

	_, err = client.Collection("DialyQuest").Doc(userDoc[0].Ref.ID).Set(ctx, dialyQuest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update success"})
		return
	}

	c.JSON(http.StatusOK, dialyQuest)
	return
}
