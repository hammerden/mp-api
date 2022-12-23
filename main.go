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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/letthefireflieslive/mp-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"os"
)

var mealPlansHandler *handlers.MealPlansHandler
var authHandler *handlers.AuthHandler

func init() {
	ctx := context.Background() //TODO: Understand https://go.dev/blog/context
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("mealPlan")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println("Redis Status: ", status)

	mealPlansHandler = handlers.NewMealPlansHandler(ctx,
		collection, redisClient)
	authHandler = &handlers.AuthHandler{}
}

func main() {
	router := gin.Default()
	router.GET("/mp", mealPlansHandler.ListMealPlansHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/mp", mealPlansHandler.NewMealPlanHandler)
		authorized.PUT("/mp/:id", mealPlansHandler.UpdateMealPlanHandler)
	}
	router.Run()
}
