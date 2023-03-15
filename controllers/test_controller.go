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

// create set of test
func CreateSetTestStress(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	allChoice := [4]models.Choice{
		{
			Number: 1,
			Text:   "ไม่เคยเลย",
			Value:  0,
		},
		{
			Number: 2,
			Text:   "เป็นครั้งคราว",
			Value:  1,
		},
		{
			Number: 3,
			Text:   "เป็นบ่อย",
			Value:  2,
		}, {
			Number: 4,
			Text:   "เป็นประจำ",
			Value:  3,
		},
	}

	// create set of test
	mocktestStress := []models.Test{
		{
			Number:       1,
			QuestionText: "นอนไม่หลับเพราะคิดมากหรือกังวลใจ",
			Choices:      allChoice,
		},
		{
			Number:       2,
			QuestionText: "รู้สึกหงุดหงิด รำคาญใจ",
			Choices:      allChoice,
		},
		{
			Number:       3,
			QuestionText: "ทำอะไรไม่ได้เลยเพราะประสาทตึงเครียด",
			Choices:      allChoice,
		},
		{
			Number:       4,
			QuestionText: "มีความวุ่นวายใจ",
			Choices:      allChoice,
		},
		{
			Number:       5,
			QuestionText: "ไม่อยากพบปะผู้คน",
			Choices:      allChoice,
		},
		{
			Number:       6,
			QuestionText: "ปวดหัวข้างเดียว หรือปวดบริเวณขมับทั้ง 2 ข้าง",
			Choices:      allChoice,
		},
		{
			Number:       7,
			QuestionText: "รู้สึกไม่มีความสุขและเศร้าหมอง",
			Choices:      allChoice,
		},
		{
			Number:       8,
			QuestionText: "รู้สึกหมดหวังในชีวิต",
			Choices:      allChoice,
		},
		{
			Number:       9,
			QuestionText: "รู้สึกชีวิตตนเองไม่มีคุณค่า",
			Choices:      allChoice,
		},
		{
			Number:       10,
			QuestionText: "กระวนกระวายอยู่ตลอดเวลา",
			Choices:      allChoice,
		},
		{
			Number:       11,
			QuestionText: "รู้สึกว่าตนเองไม่มีสมาธิ",
			Choices:      allChoice,
		},
		{
			Number:       12,
			QuestionText: "รู้สึกเพลียจนไม่มีแรงจะทำอะไร",
			Choices:      allChoice,
		},
		{
			Number:       13,
			QuestionText: "รู้สึกเหนื่อยหน่ายไม่อยากทำอะไร",
			Choices:      allChoice,
		},
		{
			Number:       14,
			QuestionText: "มีอาการหัวใจเต้นแรง",
			Choices:      allChoice,
		},
		{
			Number:       15,
			QuestionText: "เสียงสั่น ปากสั่น หรือมือสั่นเวลาไม่พอใจ",
			Choices:      allChoice,
		},
		{
			Number:       16,
			QuestionText: "รู้สึกกลัวผิดพลาดในการทำสิ่งต่าง ๆ",
			Choices:      allChoice,
		},
		{
			Number:       17,
			QuestionText: "ปวดหรือเกร็งกล้ามเนื้อบริเวณท้ายทอย หลัง หรือไหล่",
			Choices:      allChoice,
		},
		{
			Number:       18,
			QuestionText: "ตื่นเต้นง่ายกับเหตุการณ์ที่ไม่คุ้นเคย",
			Choices:      allChoice,
		},
		{
			Number:       19,
			QuestionText: "มึนงงหรือเวียนศีรษะ",
			Choices:      allChoice,
		},
		{
			Number:       20,
			QuestionText: "ความสุขทางเพศลดลง",
			Choices:      allChoice,
		},
	}

	// add to data base
	for index, _ := range mocktestStress {
		_, err := client.Collection("TestStress").Doc(strconv.Itoa(index+1)).Set(ctx, mocktestStress[index])
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

func GetTestStressResult(c *gin.Context) {

	// declare instance of fiirestore
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	// get id that's to use
	userID := c.Request.Header.Get("id")
	testRes := models.TestResualtResponse{}

	testData, err := client.Collection("ResultStress").Where("UserID", "==", client.Collection("User").Doc(userID)).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant to find result by this userID",
			"error":   err.Error(),
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
