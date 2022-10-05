package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"os"
	"time"
)

var mealPlans []MealPlan

func init() {
	mealPlans = make([]MealPlan, 0)
	file, _ := os.ReadFile("meal_plans.json")
	_ = json.Unmarshal([]byte(file), &mealPlans)

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
	CreatedAt          time.Time `json:"createdAt"`
}

func main() {
	router := gin.Default()
	router.POST("/mealplans", NewMealPlanHandler)
	router.GET("/mealplans", ListMealPlanHandler)
	router.PUT("/mealplans/:id", UpdateMealPlanHandler)
	router.Run()
}

func UpdateMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	var mealPlan MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(mealPlans); i++ {
		if mealPlans[i].ID == id {
			index = i
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

func ListMealPlanHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mealPlans)
}
