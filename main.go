package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormData struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	Age         int    `json:"age"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Salary      int    `json:"salary"`
	Designation string `json:"designation"`
}

type Delete_Data struct {
	Id string `json:"_id"`
}

var client *mongo.Client

func main() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	r.LoadHTMLFiles("form.html")
	r.GET("/", serveForm)
	r.POST("/submit", submitHandler)
	r.POST("/delete", DeleteData)
	r.POST("/update", UpdateDetails)
	r.Run(":8080")
}

func serveForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func submitHandler(c *gin.Context) {
	formData := FormData{
		FirstName:   c.PostForm("first_name"),
		MiddleName:  c.PostForm("middle_name"),
		Age:         parseAge(c.PostForm("age")),
		Location:    c.PostForm("location"),
		Email:       c.PostForm("email"),
		Salary:      parseSalary(c.PostForm("salary")),
		Designation: c.PostForm("designation"),
	}

	collection := client.Database("venkat_naidu").Collection("categories")
	_, err := collection.InsertOne(context.Background(), formData)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully", "data": formData})

}

func DeleteData(c *gin.Context) {
	idStr := c.PostForm("_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.M{"_id": id}

	collection := client.Database("venkat_naidu").Collection("categories")

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No document found with the given ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data deleted successfully", "id": idStr})
}

func UpdateDetails(c *gin.Context) {
	idStr := c.PostForm("_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	formData := FormData{
		FirstName:   c.PostForm("first_name"),
		MiddleName:  c.PostForm("middle_name"),
		Age:         parseAge(c.PostForm("age")),
		Location:    c.PostForm("location"),
		Email:       c.PostForm("email"),
		Salary:      parseSalary(c.PostForm("salary")),
		Designation: c.PostForm("designation"),
	}

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{
		"first_name":  formData.FirstName,
		"middle_name": formData.MiddleName,
		"age":         formData.Age,
		"location":    formData.Location,
		"email":       formData.Email,
		"salary":      formData.Salary,
		"designation": formData.Designation,
	}}

	collection := client.Database("venkat_naidu").Collection("categories")
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No document found with the given ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully", "data": formData})
}

func parseAge(ageStr string) int {
	age, _ := strconv.Atoi(ageStr)
	return age
}

func parseSalary(salaryStr string) int {
	salary, _ := strconv.Atoi(salaryStr)
	return salary
}
