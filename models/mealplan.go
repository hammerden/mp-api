package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:parameters mealPlans MealPlan
type MealPlan struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	Customer           string             `json:"customer" bson:"customer"`           //todo Name object
	Diet               string             `json:"diet" bson:"diet"`                   //todo Enum
	ContactNumber      string             `json:"contactNumber" bson:"contactNumber"` //todo double check type for phone numbers
	Allergies          []string           `json:"allergies" bson:"allergies"`
	AvoidedIngredients []string           `json:"avoidedIngredients" bson:"avoidedIngredients"`
	DeliveryMonday     time.Time          `json:"deliveryMonday" bson:"DeliveryMonday" `  //todo Delivery object and other days, validator
	DeliveryTuesday    time.Time          `json:"deliveryTuesday" bson:"deliveryTuesday"` //todo Delivery object and other days, validator
	Tags               []string           `json:"tags" bson:"tags"`
	CreatedAt          time.Time          `json:"createdAt" bson:"createdAt"`
}
