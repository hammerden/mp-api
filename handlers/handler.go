package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/letthefireflieslive/mp-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type MealPlansHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewMealPlansHandler(ctx context.Context, collection *mongo.Collection) *MealPlansHandler {
	return &MealPlansHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation GET /mealPlans listMealPlans
// Returns list of mealPlans
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *MealPlansHandler) ListMealPlansHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
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
	c.JSON(http.StatusOK, mealPlans)
}

// swagger:operation POST /mealPlans newMealPlan
// Create a new mealPlan
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

	c.JSON(http.StatusOK, mealPlan)
}

// swagger:operation PUT /mealPlans/{id} mealPlans updateMealPlan
// Update an existing mealPlan
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

	c.JSON(http.StatusOK, gin.H{"message": "MealPlan has been updated"})
}

// swagger:operation DELETE /mealPlans/{id} mealPlans deleteMealPlan
// Delete an existing mealPlan
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

// swagger:operation GET /mealPlans/{id} mealPlans
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

// swagger:operation GET /mealPlans/search mealPlans findMealPlan
// Search mealPlans based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: mealPlan tag
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
/*func SearchMealPlansHandler(c *gin.Context) {
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
}*/
