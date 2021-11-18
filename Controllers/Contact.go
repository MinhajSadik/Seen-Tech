package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"example.com/seen-eg/CMS/DBManager"
	"example.com/seen-eg/CMS/Models"
	"example.com/seen-eg/CMS/Utils"
)

func isContactExisting(collection *mongo.Collection, name string) bool {

	var filter bson.M = bson.M{
		"name": name,
	}
	var results []bson.M

	b, results := Utils.FindByFilter(collection, filter)
	return (b && len(results) > 0)
}

func ContactNew(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Contact

	// Fill the received data inside an obj
	var self Models.Contact
	self.UUID, _ = SettingsGenerateUUID()
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	// Check if this obj is already existing
	if isContactExisting(collection, self.Name) {
		c.Status(500)
		return errors.New("Contact is already exist")
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}

	response, _ := json.Marshal(res) //Decode
	c.Set("Content-Type", "application/json")
	c.Status(200)
	c.Send(response)
	return nil
}

func ContactGet(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Contact

	var self Models.ContactSearch
	c.BodyParser(&self)

	var results []bson.M

	b, results := Utils.FindByFilter(collection, self.GetContactSearchBSONObj())
	if !b {
		err := errors.New("db error")
		c.Status(500).Send([]byte(err.Error()))
		return err
	}
	// Decode
	response, _ := json.Marshal(
		bson.M{"result": results},
	)
	c.Set("Content-Type", "application/json")
	c.Send(response)
	return nil
}

func ContactGetById(objID primitive.ObjectID) (Models.Contact, error) {
	var self Models.Contact

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.Contact

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func ContactGetAggregated(c *fiber.Ctx) error {

	collection := DBManager.SystemCollections.Contact

	groupbyFilter := []bson.M{
		{
			"$match": bson.M{
				"iscompany": true,
			},
		},
		{

			"$group": bson.M{
				"_id":   nil,
				"lotno": "$lotno",
				"records": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
	}
	groupInfo, err := collection.Aggregate(context.Background(), groupbyFilter)
	if err != nil {
		return err
	}

	var results []bson.M
	if err = groupInfo.All(context.Background(), &results); err != nil {
		return err
	}

	// fill results from DB into internal struct
	contactBytes, _ := json.Marshal(results)
	var contactDocs []Models.ContactAggregated
	json.Unmarshal(contactBytes, &contactDocs)

	for i, v := range contactDocs {
		contactDocs[i].CompanyData, _ = ContactGetById(v.ID)
	}

	// Decode
	response, _ := json.Marshal(
		bson.M{"result": contactDocs},
	)
	c.Set("Content-Type", "application/json")
	c.Send(response)
	return nil
}

func ContactGetPopulated(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Contact

	var self Models.ContactSearch
	c.BodyParser(&self)

	var results []bson.M

	b, results := Utils.FindByFilter(collection, self.GetContactSearchBSONObj())
	if !b {
		err := errors.New("db error or obj is not exist")
		c.Status(500).Send([]byte(err.Error()))
		return err
	}

	// Decode
	contactResults, _ := json.Marshal(results)
	var contactDocs []Models.Contact
	json.Unmarshal(contactResults, &contactDocs) //Encode
	populatedResult := make([]Models.ContactPopulated, len(contactDocs))

	for i, v := range contactDocs {
		populatedResult[i].CloneFrom(v)
		if !v.CompanyRef.IsZero() {
			tmp, _ := ContactGetById(v.CompanyRef)
			populatedResult[i].CompanyRef = tmp
		}
	}

	allpopulated, _ := json.Marshal(populatedResult)
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func ContactSetStatus(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Contact

	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("Invalid request params")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var newValue = true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}

	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	c.Set("Content-Type", "application/json")
	if updateErr != nil {
		c.Status(500)
		return errors.New("An error occurred when modifing storage condition status")
	} else {
		c.Status(200)

	}
	return nil
}

func ContactModify(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Contact
	// Fill the received data inside an obj
	var self Models.Contact
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": self.GetId()}, updateData)
	c.Set("Content-Type", "application/json")
	if updateErr != nil {
		c.Status(500)
		return errors.New("An error occurred when modifing storage condition status")
	} else {
		c.Status(200)
	}
	return nil
}
