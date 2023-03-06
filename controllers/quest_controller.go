package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
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

	// check diary queest
	questDoc, err := client.Collection("DialyQuest").Where("UserID", "==", client.Collection("User").Doc(userID)).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
		return
	}

	// mapping data to object quest
	if len(questDoc) < 1 {

		// case 1 not found quest qcreate quest
		setQuest := []*models.Quest{}
		questName := []string{"PostQuest", "CommentQuest", "LikeQuest"}
		for _, name := range questName {
			setQuest = append(setQuest, &models.Quest{
				QuestName:      name,
				Progress:       0,
				CountGet:       0,
				MaxProgress:    3,
				Reward:         5,
				IsGetPoint:     false,
				Completed:      false,
				LastCompletion: time.Now().UTC(),
			})
		}

		// set diary quest to collection
		diaryQuest := models.DailyQuestProgress{
			UserID:    client.Collection("User").Doc(userID),
			QuestDate: time.Now().UTC(),
			Quests:    setQuest,
		}

		// set diary quest to collection
		_, err := client.Collection("DialyQuest").Doc(string(uuid.NewString())).Set(ctx, diaryQuest)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "cant to update user data quest"})
			return
		}

		// reutn quest that's already create quest
		c.JSON(http.StatusOK, diaryQuest)
		return

		// case quest equal 1
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
					CountGet:       0,
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

func UpdateProgressQuest(c *gin.Context, questName string) {

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

	if len(userDoc) > 1 {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to get your quest"})
		return
	}

	// mapping data
	mapstructure.Decode(userDoc[0].Data(), &dialyQuest)

	if dialyQuest.Quests != nil {
		for index, _ := range dialyQuest.Quests {
			if dialyQuest.Quests[index].QuestName == questName {
				// updat progess quest
				if dialyQuest.Quests[index].Completed {
					return
				} else {

					// case not complte
					dialyQuest.Quests[index].Progress += 1

					// set new reward
					if dialyQuest.Quests[index].IsGetPoint {
						dialyQuest.Quests[index].Reward = (dialyQuest.Quests[index].Progress - dialyQuest.Quests[index].CountGet) * 5
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
}

func GetPointFromQuest(c *gin.Context) {

	// declare instance of firestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get all id to use
	userID := c.Request.Header.Get("id")

	// declare object to use
	dialyQuest := models.DailyQuestProgress{}
	questNameBody := c.Param("questName")
	myPoint, err := strconv.Atoi(c.Param("myPoint"))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
		return
	}
	// check diary quest
	questDoc, err := client.Collection("DialyQuest").Where("UserID", "==", client.Collection("User").Doc(userID)).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
		return
	}

	if len(questDoc) > 1 || len(questDoc) == 0 {
		fmt.Println(len(questDoc))
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to get your quest llll"})
		return
	}

	mapstructure.Decode(questDoc[0].Data(), &dialyQuest)

	getPoint := 0
	// set to diary quest
	for _, q := range dialyQuest.Quests {
		if q.QuestName == questNameBody {

			q.IsGetPoint = false
			getPoint = q.Reward + myPoint
			q.Reward = (q.Progress - q.CountGet) * 5
			q.CountGet += 1
			if q.Progress == q.MaxProgress {
				q.Completed = true
				q.LastCompletion = time.Now().UTC()
			}
			break
		}
	}

	// set diary quest
	_, err = client.Collection("DialyQuest").Doc(questDoc[0].Ref.ID).Set(ctx, dialyQuest)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
		return
	}

	// set to user collection
	_, err = client.Collection("User").Doc(userID).Update(ctx, []firestore.Update{
		{
			Path:  "Point",
			Value: getPoint,
		},
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cant to update user data quest"})
		return
	}

	c.JSON(http.StatusOK, dialyQuest)
}
