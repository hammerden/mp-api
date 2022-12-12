package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/letthefireflieslive/mp-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type MealPlansHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewMealPlansHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *MealPlansHandler {
	return &MealPlansHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation GET /mp mealPlans list-all
// Returns list of mealPlans
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *MealPlansHandler) ListMealPlansHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(handler.ctx, "mealPlans").Result()

	if err == redis.Nil {
		log.Print("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx,
			bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)
		mealPlans := make([]models.MealPlan, 0)
		for cur.Next(handler.ctx) {
			var mealPlan models.MealPlan
			cur.Decode(&mealPlan)
			mealPlans = append(mealPlans, mealPlan)
		}

		data, _ := json.Marshal(mealPlans)
		handler.redisClient.Set(handler.ctx, "mealPlans", string(data), 0)
		c.JSON(http.StatusOK, mealPlans)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		mealPlans := make([]models.MealPlan, 0)
		json.Unmarshal([]byte(val), &mealPlans)
		c.JSON(http.StatusOK, mealPlans)
	}
}

// swagger:operation POST /mp mealPlans create
// New mealPlan
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//
//	'400':
//	    description: Invalid input
func (handler *MealPlansHandler) NewMealPlanHandler(c *gin.Context) {
	var mealPlan models.MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mealPlan.ID = primitive.NewObjectID()
	mealPlan.CreatedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, mealPlan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new mealPlan"})
		return
	}

	log.Println("Remove data from Redis")
	handler.redisClient.Del(handler.ctx, "mealPlans")

	c.JSON(http.StatusOK, mealPlan)
}

// swagger:operation PUT /mp/{id} mealPlans update
// Existing mealPlan
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the mealPlan
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
//	    description: Invalid mealPlan ID
func (handler *MealPlansHandler) UpdateMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	var mealPlan models.MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"customer", mealPlan.Customer},
		{"instructions", mealPlan.Diet},
		{"ingredients", mealPlan.ContactNumber},
		{"allergies", mealPlan.Allergies},
		{"avoidedIngredients", mealPlan.AvoidedIngredients},
		{"deliveryMonday", mealPlan.DeliveryMonday},
		{"deliveryTuesday", mealPlan.DeliveryTuesday},
		{"tags", mealPlan.Tags},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.redisClient.Del(handler.ctx, "mealPlans")
	c.JSON(http.StatusOK, gin.H{"message": "MealPlan has been updated"})
}

// swagger:operation DELETE /mp/{id} mealPlans delete
// Remove an existing mealPlan
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the mealPlan
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid mealPlan ID
func (handler *MealPlansHandler) DeleteMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MealPlan has been deleted"})
}

// swagger:operation GET /mp/{id} mealPlans describe
// Get one mealPlan
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: mealPlan ID
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *MealPlansHandler) GetOneMealPlanHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var mealPlan models.MealPlan
	err := cur.Decode(&mealPlan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mealPlan)
}
