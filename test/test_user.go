package test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"jitD/configs"
// 	"jitD/models"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// )

// func TestCreateUser(t *testing.T) {
// 	// Set up a mock Gin context with a request and response recorder
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	// Create a mock user to pass in the request body
// 	user := models.User{
// 		Name:     "John",
// 		Email:    "john@example.com",
// 		Password: "password",
// 	}

// 	// Convert the user to JSON and set it as the request body
// 	jsonUser, _ := json.Marshal(user)
// 	req, _ := http.NewRequest(http.MethodPost, "/create-user", bytes.NewBuffer(jsonUser))
// 	c.Request = req

// 	// Call the CreateUser function
// 	CreateUser(c)

// 	// Check the response status code
// 	if w.Code != http.StatusCreated {
// 		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
// 	}

// 	// Check that the user was created in the database
// 	ctx := context.Background()
// 	client := configs.CreateClient(ctx)
// 	doc, err := client.Collection("User").Doc("1").Get(ctx)
// 	if err != nil {
// 		t.Errorf("Error getting document: %v", err)
// 	}
// 	var createdUser models.User
// 	err = doc.DataTo(&createdUser)
// 	if err != nil {
// 		t.Errorf("Error converting document data: %v", err)
// 	}
// 	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.Password != user.Password {
// 		t.Errorf("Expected user %v but got %v", user, createdUser)
// 	}
// }
