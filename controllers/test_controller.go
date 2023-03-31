package controllers

import (
	"context"
	configs "jitD/configs"
	"jitD/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// get stress test
func GetTest(c *gin.Context, testName string) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// retrieve test data
	var testStress models.Test
	settestStress := []models.Test{}

	// get all doument of test stress
	testData, err := client.Collection(testName).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get test stress",
			"error":   err.Error(),
		})
		return
	}

	//
	for i, _ := range testData {
		testStress = models.Test{}
		mapstructure.Decode(testData[i].Data(), &testStress)
		settestStress = append(settestStress, testStress)
	}

	c.JSON(http.StatusOK, settestStress)
}

func CalTestStress(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get point from param
	point, err := strconv.Atoi(c.Param("point"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get poiint int this param",
			"error":   err.Error(),
		})
		return
	}

	result := ""
	desc := ""

	// get id that's to use
	userID := c.Request.Header.Get("id")
	testRes := models.TestResualtResponse{}

	// check point
	if point <= 5 {
		result = "คุณมีความเครียดน้อยว่าปกติ"
		desc = "อาจเป็นเพราะคุณมีชีวิตที่เรียบง่าย ไม่ค่อยมีเรื่องให้ต้องตื่นเต้น และไม่ค่อยกระตือรือร้น"
	} else if point <= 17 {
		result = "คุณมีความเครียดในระดับปกติ"
		desc = "คุณสามารถจัดการกับความเครียดที่เกิดขึ้นในชีวิตประจำวันได้ดี และสามารถปรับตัว ปรับใจ ให้เข้ากับสถานการณ์ต่าง ๆ ได้อย่างถูกต้องเหมาะสม คุณควรพยายามคงระดับความเครียดในระดับนี้ต่อไปให้นาน ๐ "
	} else if point <= 25 {
		result = "คุณมีความเครียดสูงกว่าระดับปกติเล็กน้อย"
		desc = "แสดงว่าคุณอาจกำลังมีปัญหาบางอย่างที่ทำให้ไม่สบายใจอยู่ อาจทำให้มีอาการผิดปกติทางร่างกาย จิตใจ และพฤติกรรมเล็กน้อย คุณควรผ่อนคลายความเครียด"
	} else {
		result = "คุณมีความเครียดสูงกว่าปกติ"
		desc = "แสดงว่าคุณอาจมีปัญหาบางอย่างในชีวิตที่ยังหาทางแก้ไขไม่ได้ ทำให้มีอาการผิดปกติทางร่างกาย จิตใจ และพฤติกรรมอย่างเห็นได้ชัด คุณควรฝึกเทคนิคเฉพาะในการคลายเครียด"
	}

	test := models.TestResualt{
		UserID:   client.Collection("User").Doc(userID),
		TestDate: time.Now().UTC(),
		TestName: "Test Stress",
		Point:    point,
		Result:   result,
		Desc:     desc,
	}

	// set data in user
	testData, err := client.Collection("ResultStress").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if len(testData) == 0 {
		client.Collection("ResultStress").Add(ctx, test)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to add test result",
				"error":   err.Error(),
			})
			return
		}
	} else {
		client.Collection("ResultStress").Doc(testData[0].Ref.ID).Set(ctx, test)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to set test result",
				"error":   err.Error(),
			})
			return
		}
	}

	mapstructure.Decode(test, &testRes)
	testRes.TestDate = test.TestDate
	c.JSON(http.StatusOK, testRes)
}

func GetTestResult(c *gin.Context, testName string) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	userID := c.Request.Header.Get("id")
	testRes := models.TestResualtResponse{}

	testData, err := client.Collection(testName).Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find result by this userID",
		})
		return
	}

	if len(testData) == 0 {
		c.JSON(http.StatusOK, models.TestResualtResponse{
			TestDate: time.Now().UTC(),
			TestName: "Test Stress",
			Point:    0,
			Result:   "No data",
			Desc:     "please to do this test",
		})
		return
	} else {
		mapstructure.Decode(testData[0].Data(), &testRes)
		c.JSON(http.StatusOK, testRes)
		return
	}

}

func CalTestConsult(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get point from param
	point, err := strconv.Atoi(c.Param("point"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to get poiint int this param",
			"error":   err.Error(),
		})
		return
	}

	result := ""
	desc := ""

	// get id that's to use
	userID := c.Request.Header.Get("id")
	testRes := models.TestResualtResponse{}

	// check point
	if point <= 5 {
		result = "คุณมีความเครียดน้อยว่าปกติ"
		desc = "อาจเป็นเพราะคุณมีชีวิตที่เรียบง่าย ไม่ค่อยมีเรื่องให้ต้องตื่นเต้น และไม่ค่อยกระตือรือร้น"
	} else if point <= 17 {
		result = "คุณมีความเครียดในระดับปกติ"
		desc = "คุณสามารถจัดการกับความเครียดที่เกิดขึ้นในชีวิตประจำวันได้ดี และสามารถปรับตัว ปรับใจ ให้เข้ากับสถานการณ์ต่าง ๆ ได้อย่างถูกต้องเหมาะสม คุณควรพยายามคงระดับความเครียดในระดับนี้ต่อไปให้นาน ๐ "
	} else if point <= 25 {
		result = "คุณมีความเครียดสูงกว่าระดับปกติเล็กน้อย"
		desc = "แสดงว่าคุณอาจกำลังมีปัญหาบางอย่างที่ทำให้ไม่สบายใจอยู่ อาจทำให้มีอาการผิดปกติทางร่างกาย จิตใจ และพฤติกรรมเล็กน้อย คุณควรผ่อนคลายความเครียด"
	} else {
		result = "คุณมีความเครียดสูงกว่าปกติ"
		desc = "แสดงว่าคุณอาจมีปัญหาบางอย่างในชีวิตที่ยังหาทางแก้ไขไม่ได้ ทำให้มีอาการผิดปกติทางร่างกาย จิตใจ และพฤติกรรมอย่างเห็นได้ชัด คุณควรฝึกเทคนิคเฉพาะในการคลายเครียด"
	}

	test := models.TestResualt{
		UserID:   client.Collection("User").Doc(userID),
		TestDate: time.Now().UTC(),
		TestName: "Test Stress",
		Point:    point,
		Result:   result,
		Desc:     desc,
	}

	// set data in user
	testData, err := client.Collection("ResultStress").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if len(testData) == 0 {
		client.Collection("ResultStress").Add(ctx, test)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to add test result",
				"error":   err.Error(),
			})
			return
		}
	} else {
		client.Collection("ResultStress").Doc(testData[0].Ref.ID).Set(ctx, test)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "cant to set test result",
				"error":   err.Error(),
			})
			return
		}
	}

	mapstructure.Decode(test, &testRes)
	testRes.TestDate = test.TestDate
	c.JSON(http.StatusOK, testRes)
}
