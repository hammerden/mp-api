// Meal Plan API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
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
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var mealPlans []MealPlan
var ctx context.Context
var err error
var client *mongo.Client

func init() {
	mealPlans = make([]MealPlan, 0)
	file, _ := os.ReadFile("meal_plans.json")
	_ = json.Unmarshal([]byte(file), &mealPlans)
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	var listOfMealPlans []interface{}
	for _, mealPlan := range mealPlans {
		listOfMealPlans = append(listOfMealPlans, mealPlan)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("hd")
	insertManyResult, err := collection.InsertMany(ctx, listOfMealPlans)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted mealPlans: ", len(insertManyResult.InsertedIDs))

}

type MealPlan struct {
	ID                 string    `json:"id"`
	Customer           string    `json:"customer"`      //todo Name object
	Diet               string    `json:"diet"`          //todo Enum
	ContactNumber      string    `json:"contactNumber"` //todo double check type for phone numbers
	Allergies          []string  `json:"allergies"`
	AvoidedIngredients []string  `json:"avoidedIngredients"`
	DeliveryMonday     time.Time `json:"deliveryMonday"'`  //todo Delivery object and other days, validator
	DeliveryTuesday    time.Time `json:"deliveryTuesday"'` //todo Delivery object and other days, validator
	Tags               []string  `json:"tags"`
	CreatedAt          time.Time `json:"createdAt"`
}

func main() {
	router := gin.Default()
	router.POST("/mealplans", NewMealPlanHandler)
	router.GET("/mealplans", ListMealPlanHandler)
	router.PUT("/mealplans/:id", UpdateMealPlanHandler)
	router.DELETE("/mealplans/:id", DeleteMealPlanHandler)
	router.GET("/mealplans/search", SearchMealPlanHandler)
	router.Run()
}

func SearchMealPlanHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfMealPlans := make([]MealPlan, 0)
	for i := 0; i < len(mealPlans); i++ {
		found := false
		for _, t := range mealPlans[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfMealPlans = append(listOfMealPlans, mealPlans[i])
		}

	}
	c.JSON(http.StatusOK, listOfMealPlans)
}

func DeleteMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(mealPlans); i++ {
		if mealPlans[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Meal not found",
		})
		return
	}
	mealPlans = append(mealPlans[:index], mealPlans[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Meal Plan has been deleted",
	})
}

// swagger:operation PUT /mealplans/{id} mealplans updateMealPlan
// Update an existing meal plan
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the mealplan
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
//	'404':
//	    description: Invalid mealplan ID
func UpdateMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	var mealPlan MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	mealPlan.ID = id
	index := -1
	for i := 0; i < len(mealPlans); i++ {
		if mealPlans[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipes not found",
		})
		return
	}
	mealPlans[index] = mealPlan
	c.JSON(http.StatusOK, mealPlan)
}

func NewMealPlanHandler(c *gin.Context) {
	var mealPlan MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	mealPlan.ID = xid.New().String()
	mealPlan.CreatedAt = time.Now()
	mealPlans = append(mealPlans, mealPlan)
	c.JSON(http.StatusOK, mealPlan)
}

// swagger:operation GET /mealplans mealplans listMealPlans
// Returns list of mealplans
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func ListMealPlanHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mealPlans)
}
