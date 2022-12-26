package controllers

import (
	"context"
	"fmt"
	mapstructure "github.com/mitchellh/mapstructure"
	configs "jitD/configs"
	models "jitD/models"
	"log"

	"net/http"

	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
)

// Retrive all book
func GetAllBook(c *gin.Context) {

	books := []models.Book{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	iter := client.Collection("Book").Documents(ctx)
	for {
		book := models.Book{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, "msg: Not found")
		}
		mapstructure.Decode(doc.Data(), &book)
		books = append(books, book)
	}

	fmt.Println(books)
	c.JSON(http.StatusOK, books)
}

// Retrive book by id
func GetBookById(c *gin.Context) {

	id := c.Param("id")
	book := models.Book{}

	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Book").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	mapstructure.Decode(dsnap.Data(), &book)
	c.JSON(http.StatusOK, book)
}
