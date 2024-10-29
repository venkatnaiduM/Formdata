package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	fmt.Println("Server is running on :8080")
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
		c.String(http.StatusInternalServerError, "Failed to insert data")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully", "data": formData})

}

func parseAge(ageStr string) int {
	age, _ := strconv.Atoi(ageStr)
	return age
}

func parseSalary(salaryStr string) int {
	salary, _ := strconv.Atoi(salaryStr)
	return salary
}
