package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"time"
)

var mealPlans []MealPlan

func init() {
	mealPlans = make([]MealPlan, 0)
}

type MealPlan struct {
	ID                 string   `json:"id"`
	Customer           string   `json:"customer"`      //todo Name object
	Diet               string   `json:"diet"`          //todo Enum
	ContactNumber      string   `json:"contactNumber"` //todo double check type for phone numbers
	Allergies          []string `json:"allergies"`
	AvoidedIngredients []string `json:"avoidedIngredients"`
	//deliveryMonday     time.Time `json:"deliveryMonday" binding:"required,bookabledate" time_format:"2006-01-02"` //todo Delivery object and other days
	CreatedAt time.Time `json:"createdAt"`
}

//func bookableDate(
//	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
//	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
//) bool {
//	if date, ok := field.Interface().(time.Time); ok {
//		today := time.Now()
//		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
//			return false
//		}
//	}
//	return true
//}

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

func main() {
	router := gin.Default()
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("bookabledate", bookableDate)
	//}
	router.POST("/mealplans", NewMealPlanHandler)
	router.Run()
}
