package controller_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/controller"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRouter(client *mongo.Client) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLFiles("html/form.html")
	r.GET("/", controller.ServeForm)
	r.POST("/submit", controller.SubmitHandler(client))
	r.POST("/delete", controller.DeleteData(client))
	r.POST("/update", controller.UpdateDetails(client))
	return r
}

func TestSubmitHandler(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mtest.NewClientOptions()))
	defer mt.Close()

	mt.Run("Submit valid data", func(mt *mtest.T) {
		router := setupRouter(mt.Client)

		body := bytes.NewBufferString("first_name=John&middle_name=Doe&age=30&location=City&email=test@example.com&salary=50000&designation=Engineer")
		req, _ := http.NewRequest(http.MethodPost, "/submit", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", recorder.Code)
		}
	})

	mt.Run("Submit invalid data", func(mt *mtest.T) {
		router := setupRouter(mt.Client)

		body := bytes.NewBufferString("first_name=&age=-1&location=&email=&salary=0&designation=")
		req, _ := http.NewRequest(http.MethodPost, "/submit", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusForbidden {
			t.Errorf("Expected status Forbidden; got %v", recorder.Code)
		}
	})
}

func TestDeleteData(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mtest.NewClientOptions()))
	defer mt.Close()

	mt.Run("Delete existing document", func(mt *mtest.T) {
		collection := mt.Coll
		docID := primitive.NewObjectID()
		_, err := collection.InsertOne(context.Background(), bson.M{"_id": docID, "first_name": "John"})
		if err != nil {
			t.Fatalf("Could not insert mock document: %v", err)
		}

		router := setupRouter(mt.Client)

		body := bytes.NewBufferString("_id=" + docID.Hex())
		req, _ := http.NewRequest(http.MethodPost, "/delete", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", recorder.Code)
		}
	})
}

func TestUpdateDetails(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mtest.NewClientOptions()))
	defer mt.Close()

	mt.Run("Update existing document", func(mt *mtest.T) {
		collection := mt.Coll
		docID := primitive.NewObjectID()
		_, err := collection.InsertOne(context.Background(), bson.M{"_id": docID, "first_name": "John"})
		if err != nil {
			t.Fatalf("Could not insert mock document: %v", err)
		}

		router := setupRouter(mt.Client)

		body := bytes.NewBufferString("_id=" + docID.Hex() + "&first_name=Jane")
		req, _ := http.NewRequest(http.MethodPost, "/update", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", recorder.Code)
		}
	})
}
