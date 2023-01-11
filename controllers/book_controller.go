package controllers

import (
	"context"
	"fmt"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

// Retrive all book
func GetAllBook(c *gin.Context) {

	books := []models.Book{}
	ctx := context.Background()
	client := configs.CreateClient(c)

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
		print(doc.Data())
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

// get seller by id
func GetSellerById(c *gin.Context) {

	// id := c.Param("id")
	// book := models.Book{}
	seller_x := []models.Seller{}
	xxx := models.Seller{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()

	doc, err := client.Collection("Book").Doc("lmW8BLciqMyUMLglxAZ9").Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found 1",
		})
		return
	}

	sell, err := doc.DataAt("sallers")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found 2",
		})
		return
		//log.Fatal(err)
	}

	type Seller_list []struct {
		ID string `json:"ID"`
	}

	var s Seller_list

	mapstructure.Decode(sell, &s)

	for _, s := range s {
		dsnap, err := client.Collection("Seller").Doc(s.ID).Get(ctx)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Not found 2",
			})
			return
			//log.Fatal(err)
		}

		mapstructure.Decode(dsnap.Data(), &xxx)
		seller_x = append(seller_x, xxx)

	}

	c.JSON(http.StatusOK, seller_x)

}

// Add a books
func AddBook(c *gin.Context) {

	// create conteext
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	book := models.Book{}

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&book); err != nil {
		return
	}

	_, _, err := client.Collection("Book").Add(ctx, book)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message": "create data success",
		})
	}
}

// delete a book
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	dsnap, err := client.Collection("Book").Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete fail",
		})
	} else {
		print(dsnap)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete success",
		})
	}
}

// update a book
func UpdateBbook(c *gin.Context) {
}
