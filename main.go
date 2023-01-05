// Meal Plan API
//
// This is a sample mealPlans API. You can find out more about the API at https://github.com/hammerden/mp-api
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
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
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
var redisUri, mongoDatabase, mongoURI string

func init() {
	setEnvConfig()
	ctx := context.Background() //TODO: Understand https://go.dev/blog/context
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(mongoURI))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(mongoDatabase).Collection("mealPlan")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println("Redis Status: ", status)

	mealPlansHandler = handlers.NewMealPlansHandler(ctx,
		collection, redisClient)

	collectionUsers := client.Database(mongoDatabase).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func setEnvConfig() {
	redisUri = os.Getenv("REDIS_URI")
	if redisUri == "" {
		log.Fatal("REDIS_URI env var not found!")
	}

	mongoDatabase = os.Getenv("MONGO_DATABASE")
	if os.Getenv("MONGO_DATABASE") == "" {
		log.Fatal("MONGO_DATABASE env var not found!")
	}

	mongoURI = os.Getenv("MONGO_URI")
	if os.Getenv("MONGO_URI") == "" {
		log.Fatal("MONGO_URI env var not found!")
	}
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp",
		redisUri, "", []byte("secret"))
	router.Use(sessions.Sessions("meal_plans_api", store))

	router.GET("/mp", mealPlansHandler.ListMealPlansHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/mp", mealPlansHandler.NewMealPlanHandler)
		authorized.PUT("/mp/:id", mealPlansHandler.UpdateMealPlanHandler)
	}
	router.RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")
}
