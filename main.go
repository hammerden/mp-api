// Meal Plan API
//
// This is a sample mealPlans API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/letthefireflieslive/mp-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"os"
)

var mealPlansHandler *handlers.MealPlansHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("mealPlan")
	mealPlansHandler = handlers.NewMealPlansHandler(ctx,
		collection)
}

func main() {
	router := gin.Default()
	router.POST("/mealpslans", mealPlansHandler.NewMealPlanHandler)
	router.GET("/mealplans", mealPlansHandler.ListMealPlansHandler)
	router.PUT("/mealplans/:id", mealPlansHandler.UpdateMealPlanHandler)
	router.Run()
}
