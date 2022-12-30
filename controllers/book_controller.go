package controllers

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gin-gonic/gin"
	mapstructure "github.com/mitchellh/mapstructure"
	// "google.golang.org/api/iterator"
	configs "jitD/configs"
	models "jitD/models"
	"log"
	"net/http"
)

// Retrive all book
func GetAllBook(c *gin.Context) {

	// create client
	books := []models.Book{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	//get all books
	docRef, err := client.Collection("Book").Documents(ctx).GetAll()
	if err!= nil {
      log.Fatal(err)
	}

	for _, data := range docRef {
		book := models.Book{}
      mapstructure.Decode(data.Data(), &book)
      books = append(books, book)
	}

	// return data
	c.JSON(http.StatusOK, books)
}

// Retrive book by id
func GetBookById(c *gin.Context) {

	// create client and set variable
	id := c.Param("id")
	book := models.Book{}
	ctx := context.Background()
	client := configs.CreateClient(ctx)

	//get all books
	dsnap, err := client.Collection("Book").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cant create",
		})
	}

	// mapping and return data
	mapstructure.Decode(dsnap.Data(), &book)
	c.JSON(http.StatusOK, book)
}

// get seller by id
func GetSellerById(c *gin.Context) {

	id := c.Param("id")
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer client.Close()

	dsnap, err := client.Collection("Book").Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found 1",
		})
		return
	}

	sell, err := dsnap.DataAt("sallers")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "NotFound a Collection",
		})
	}

	var s []*firestore.DocumentRef
	mapstructure.Decode(sell, &s)

	var seller models.Seller
	var sellers []models.Seller
	snaps, err := client.GetAll(ctx, s)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found 2",
		})
	}
	for i := 0; i < len(snaps); i++ {
		mapstructure.Decode(snaps[i].Data(), &seller)
		sellers = append(sellers, seller)
	}

	c.JSON(http.StatusOK, sellers)
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
