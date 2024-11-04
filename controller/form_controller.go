package controller

import (
	"context"
	"database/constants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// var err error

func ServeForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func SubmitHandler(client *mongo.Client, c *gin.Context) {
	formData := FormData{
		FirstName:   c.PostForm("first_name"),
		MiddleName:  c.PostForm("middle_name"),
		Age:         parseAge(c.PostForm("age")),
		Location:    c.PostForm("location"),
		Email:       c.PostForm("email"),
		Salary:      parseSalary(c.PostForm("salary")),
		Designation: c.PostForm("designation"),
	}
	if formData.FirstName == "" || formData.Age <= 0 || formData.Location == "" || formData.Email == "" || formData.Salary <= 0 || formData.Designation == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "All fields must be filled out correctly."})
		return
	}
	collection := client.Database(constants.DatabaseName).Collection(constants.CollectionName)
	_, err := collection.InsertOne(context.Background(), formData)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully", "data": formData})

}

func DeleteData(client *mongo.Client, c *gin.Context) {

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
	collection := client.Database(constants.DatabaseName).Collection(constants.CollectionName)

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

func UpdateDetails(client *mongo.Client, c *gin.Context) {
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
	if formData.FirstName == "" || formData.Age <= 0 || formData.Location == "" || formData.Email == "" || formData.Salary <= 0 || formData.Designation == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "All fields must be filled out correctly."})
		return
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

	collection := client.Database(constants.DatabaseName).Collection(constants.DatabaseName)
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
