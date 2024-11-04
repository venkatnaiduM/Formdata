package main

import (
	"context"
	"log"
	"time"

	"database/controller"
	"database/services"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := services.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	r := gin.Default()
	r.LoadHTMLFiles("form.html")

	r.GET("/", controller.ServeForm)
	r.POST("/submit", func(c *gin.Context) {
		controller.SubmitHandler(client, c)
	})
	r.POST("/delete", func(c *gin.Context) {
		controller.DeleteData(client, c)
	})
	r.POST("/update", func(c *gin.Context) {
		controller.UpdateDetails(client, c)
	})

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
