package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"example.com/seen-eg/CMS/DBManager"
	"example.com/seen-eg/CMS/Models"
	"example.com/seen-eg/CMS/Utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ChemicalSolutionInventoryRecordGetById(objID primitive.ObjectID) (Models.ChemicalSolutionInventoryRecord, error) {
	var self Models.ChemicalSolutionInventoryRecord
	b, results := Utils.FindByFilter(DBManager.SystemCollections.ChemicalSolutionInventoryRecord, bson.M{"_id": objID})
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := json.Marshal(results[0]) // Decode
	json.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func ChemicalSolutionInventoryRecordPopulatedGetById(objID primitive.ObjectID, ptr *Models.ChemicalSolutionInventoryRecord) (Models.ChemicalSolutionInventoryRecordPopulated, error) {
	var recordDoc Models.ChemicalSolutionInventoryRecord
	if ptr == nil {
		recordDoc, _ = ChemicalSolutionInventoryRecordGetById(objID)
	} else {
		recordDoc = *ptr
	}

	populatedResult := Models.ChemicalSolutionInventoryRecordPopulated{}
	populatedResult.CloneFrom(recordDoc)
	populatedResult.PreparationRef, _ = ChemicalSolutionPreparationgetById(recordDoc.PreparationRef)
	populatedResult.ChemicalSolutionTemplateRef, _ = ChemicalSolutionTemplateGetById(DBManager.SystemCollections.ChemicalSolutionTemplate, recordDoc.ChemicalSolutionTemplateRef)
	populatedResult.UserRef, _ = UserGetById(recordDoc.UserRef)
	populatedResult.DesiredAmountUoM, _ = UnitsOfMeasurementGetById(recordDoc.UserRef)
	populatedResult.ToSiteRef, _ = LabDivisionGetById(recordDoc.ToSiteRef)
	populatedResult.ToLabRef, _ = LabDivisionGetById(recordDoc.ToLabRef)
	populatedResult.ToAreaRef, _ = LabDivisionGetById(recordDoc.ToAreaRef)
	populatedResult.ToLocationRef, _ = LabDivisionGetById(recordDoc.ToLocationRef)
	return populatedResult, nil
}

func ChemicalSolutionInventoryRecordGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord
	var self Models.ChemicalSolutionInventoryRecordSearch
	c.BodyParser(&self)

	_, results := Utils.FindByFilter(collection, self.GetSearchBSONObj())

	recordsResults, _ := json.Marshal(results)
	var recordsDocs []Models.ChemicalSolutionInventoryRecord
	json.Unmarshal(recordsResults, &recordsDocs)

	populatedResult := make([]Models.ChemicalSolutionInventoryRecordPopulated, len(recordsDocs))

	for i, v := range recordsDocs {
		populatedResult[i], _ = ChemicalSolutionInventoryRecordPopulatedGetById(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(
		bson.M{"result": populatedResult},
	)

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(allpopulated)

	return nil
}

func ChemicalSolutionInventoryRecordGetAllPopulatedAggregated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord

	matchStage := bson.D{{"$match", bson.D{{"status", "Active"}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$chemicalsolutiontemplateref"}, {"totalamount", bson.D{{"$sum", "$actualamount"}}}}}}

	results := []bson.M{}
	cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	// fill results from DB into internal struct
	aggreBytes, _ := json.Marshal(results)
	var aggreDocs []Models.ChemicalSolutionInventoryAggregationResult
	json.Unmarshal(aggreBytes, &aggreDocs)

	populatedResult := make([]Models.ChemicalSolutionInventoryAggregatedAmounts, len(aggreDocs))
	for i, v := range aggreDocs {
		populatedResult[i].ChemicalSolutionRef, _ = ChemicalSolutionTemplatePopulatedGetById(v.ID, nil)
		populatedResult[i].TotalAmount = v.TotalAmount
		populatedResult[i].Records, _ = ChemicalSolutionInventoryRecordGetByChemicalSolutionRef(c, v.ID)
	}

	allpopulated, _ := json.Marshal(populatedResult)

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(allpopulated)
	return nil
}

func ChemicalSolutionInventoryRecordGetByChemicalSolutionRef(c *fiber.Ctx, objID primitive.ObjectID) ([]Models.ChemicalSolutionInventoryRecord, error) {
	var self []Models.ChemicalSolutionInventoryRecord
	b, results := Utils.FindByFilter(DBManager.SystemCollections.ChemicalSolutionInventoryRecord, bson.M{"chemicalsolutiontemplateref": objID})
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := json.Marshal(results) // Decode
	json.Unmarshal(bsonBytes, &self)      // Encode

	return self, nil
}

func ChemicalSolutionInventoryRecordGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord

	// Fill the received search obj data
	var self Models.ChemicalSolutionInventoryRecordSearch
	c.BodyParser(&self)

	b, results := Utils.FindByFilter(collection, self.GetSearchBSONObj())
	if !b {
		err := errors.New("db error")
		c.Status(500).Send([]byte(err.Error()))
		return err
	}

	response, _ := json.Marshal(results) // Decode
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

// !! when creating new ChemicalSolInventory from where I get LabDivision ID
func ChemicalSolutionInventoryRecordNew(c *fiber.Ctx, row Models.ChemicalSolutionPreparation, preparationId primitive.ObjectID) error {

	var self Models.ChemicalSolutionInventoryRecord
	uuid, err := SettingsGenerateUUID()
	if err != nil {
		return err
	}
	codeImg := QREncode(uuid)
	self.UUID = uuid
	self.PreparationRef = preparationId
	self.ChemicalSolutionTemplateRef = row.Data[0].ChemicalSolutionRef
	self.PreparationDate = row.PreparationDate
	self.ExpirationDate = row.ExpirationDate
	self.BatchNumber = row.BatchNumber
	self.UserRef = row.UserRef
	self.DesiredAmount = row.DesiredAmount
	self.DesiredAmountUoM = row.DesiredAmountUoM
	self.CodeImage = codeImg
	self.ActualAmount = row.DesiredAmount
	self.Status = "Active"

	_, err = DBManager.SystemCollections.ChemicalSolutionInventoryRecord.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}

	return nil
}

func ChemicalSolutionInventoryRecordSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord

	if c.Params("id") == "" || c.Params("new_state") == "" {
		c.Status(404)
		return errors.New("all params not sent correctilly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	_, err := ChemicalSolutionInventoryRecordGetById(objID)
	if err != nil {
		c.Status(500)
		return errors.New("chemical Solution inventory record is not found")
	}

	updateData := bson.M{
		"$set": bson.M{
			"status": c.Params("new_state"),
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred while performing the set status operation")
	}
	c.Status(200).Send([]byte("Updated successfully"))
	return nil
}

func ChemicalSolutionInventoryRecordMoveTo(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord

	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	SolRecord, err := ChemicalSolutionInventoryRecordGetById(objID)
	if err != nil {
		c.Status(500)
		return errors.New("can not make the operation")
	}

	chemicalSolTemplate, err := ChemicalSolutionTemplateGetById(DBManager.SystemCollections.ChemicalSolutionTemplate, SolRecord.ChemicalSolutionTemplateRef)
	if err != nil {
		c.Status(500)
		return errors.New("invalid chemical reference inside the chemical record")
	}

	var self Models.ChemicalSolutionInventoryRecordOperationOpen
	c.BodyParser(&self)

	location, err := LabDivisionGetById(self.ToLocationRef)
	if err != nil {
		c.Status(500)
		return errors.New("invalid destination location")
	}

	if !(location.StorageConditionRef == chemicalSolTemplate.StorageConditionRef) {
		c.Status(500)
		return errors.New("new storage condtion is not match with the new destination location's storage condition")
	}

	modificationQuery := bson.M{
		"tositeref":     self.ToSiteRef,
		"tolabref":      self.ToLabRef,
		"toarearef":     self.ToAreaRef,
		"tolocationref": self.ToLocationRef,
	}

	updateData := bson.M{
		"$set": modificationQuery,
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred while performing the open operation")
	} else {
		c.Status(200)
	}

	return nil
}

func marshalUnmarshalBsonMArrToChemicalSolutionInventoryRecord(arr []bson.M) []Models.ChemicalSolutionInventoryRecord {
	recordsResults, _ := json.Marshal(arr)
	var recordsDocs []Models.ChemicalSolutionInventoryRecord
	json.Unmarshal(recordsResults, &recordsDocs)
	return recordsDocs
}

func ChemicalSolutionInventoryExpirationNotificationGetAll(c *fiber.Ctx) error {

	if c.Params("type") == "" || !(c.Params("type") == "Medium" || c.Params("type") == "High" || c.Params("type") == "All") {
		c.Status(404)
		return errors.New("invalid request params")
	}
	exp_type := c.Params("type")

	setting_collection := DBManager.SystemCollections.Settings
	var self_setting Models.Settings
	err := setting_collection.FindOne(context.Background(), bson.M{}).Decode(&self_setting)
	if err != nil {
		return err
	}
	medium_expiration_safe_line := Utils.AdaptCurrentTimeByUnit(self_setting.MLENPeriodUnit, self_setting.MLENPeriod)
	high_expiration_safe_line := Utils.AdaptCurrentTimeByUnit(self_setting.HLENPeriodUnit, self_setting.HLENPeriod)

	filter_m := bson.M{"$and": []bson.M{
		{"expirationdate": bson.M{"$lte": medium_expiration_safe_line}},
		{"expirationdate": bson.M{"$gt": high_expiration_safe_line}},
		{"status": bson.M{"$eq": "Active"}}}}

	filter_h := bson.M{"$and": []bson.M{
		{"expirationdate": bson.M{"$lte": high_expiration_safe_line}},
		{"status": bson.M{"$eq": "Active"}}}}

	// Initiate the connection
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord
	var results_m []bson.M
	var results_h []bson.M
	_, results_m = Utils.FindByFilter(collection, filter_m)
	_, results_h = Utils.FindByFilter(collection, filter_h)

	recordsDocs_m := marshalUnmarshalBsonMArrToChemicalSolutionInventoryRecord(results_m)
	recordsDocs_h := marshalUnmarshalBsonMArrToChemicalSolutionInventoryRecord(results_h)

	var b bson.M

	if exp_type == "High" {
		b = bson.M{"result": recordsDocs_h}
	} else if exp_type == "Medium" {
		b = bson.M{"result": recordsDocs_m}
	} else {
		b = bson.M{"result": append(recordsDocs_h, recordsDocs_m...)}
	}
	expired, _ := json.Marshal(b)

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(expired)

	return nil
}

//helping function
func validateResults(obj *Models.StandardizationLogEvent) error {
	if len(obj.Results) != obj.Config.ResultsCount {
		return errors.New("results Count is not compatible")
	}
	var sum, mean, sd float64
	for e, _ := range obj.Results {
		// check tolerance is used
		if obj.Config.ToleranceIsUsed {
			// calculate margin
			margin := obj.Config.Tolerance
			if obj.Config.ToleranceType == "Percentage" {
				margin = (obj.Config.Tolerance / 100) * obj.Config.Target
			}
			// result should be within the tolerated region
			if !((obj.Config.Target-margin <= float64(e)) && (float64(e) <= obj.Config.Target+margin)) {
				return errors.New("provided result exceeds the tolerance")
			}
		}
		sum += float64(e)
	}
	mean = sum / float64(len(obj.Results))
	for i := 0; i < len(obj.Results); i++ {
		sd += math.Pow(obj.Results[i]-mean, 2)
	}

	sd = math.Sqrt(sd / float64(len(obj.Results)))
	rsd := (sd / mean) * 100
	recovery := mean / obj.Config.Target
	// check tolerance is used
	if obj.Config.ToleranceIsUsed {
		// calculate margin
		margin := obj.Config.Tolerance
		if obj.Config.ToleranceType == "Percentage" {
			margin = (obj.Config.Tolerance / 100) * obj.Config.Target
		}
		// recovery should be less than or equal tolerance
		if recovery > margin {
			return errors.New("recovery exceeds the tolerance")
		}
	}
	// RSD Validation
	switch obj.Config.RSDValidationType {
	case "Min":
		if rsd < obj.Config.RSDMin {
			return errors.New("rsd shouldn't be less than Min")
		}
	case "Max":
		if rsd > obj.Config.RSDMax {
			return errors.New("rsd shouldn't be more than Max")
		}
	case "MinMax":
		if !(rsd >= obj.Config.RSDMin && rsd <= obj.Config.RSDMax) {
			return errors.New("rsd should be in the right interval")
		}
	case "TargetTolerance":
		margin := obj.Config.RSDTolerance
		if obj.Config.ToleranceType == "Percentage" {
			margin = (obj.Config.RSDTolerance / 100.0) * obj.Config.Target
		}
		if !((obj.Config.Target-margin <= rsd) && (rsd <= obj.Config.Target+margin)) {
			return errors.New("rsd exceeds the tolerance")
		}
	}
	return nil
}

// URL/:ChemicalSolutionInventoryId
func ValidateDraftLastStdLogEvent(c *fiber.Ctx) error {
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	if !hasDraft(&solution) {
		return errors.New("there is no Draft to Validate")
	}
	size := len(solution.StandardizationLog)
	err = validateResults(&solution.StandardizationLog[size-1])
	if err != nil {
		return err
	}
	return nil
}

// URL/:ChemicalSolutionInventoryId
func ValidateAndSaveDraftLastStdLogEvent(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	// check chemicalSolInventory found
	if err != nil {
		return err
	}
	// Check if the last stdLogEvent is Draft
	if !hasDraft(&solution) {
		return errors.New("there is no Draft to Validate")
	}
	// Validate last StdLogEvent
	stdLog := solution.StandardizationLog
	size := len(stdLog)
	err = validateResults(&stdLog[size-1])
	var updateData bson.M
	// Fail
	if err != nil {
		stdLog[size-1].FinalConclusion = "Fail"
		updateData = bson.M{
			"$set": bson.M{
				"standardizationlog": stdLog,
			},
		}
	} else {
		// Pass and Update Time
		stdLog[size-1].FinalConclusion = "Pass"
		chemicalSolTemp, _ := ChemicalSolutionTemplateGetById(DBManager.SystemCollections.ChemicalSolutionTemplate, solution.ChemicalSolutionTemplateRef)
		stdLog[size-1].Date = primitive.NewDateTimeFromTime(Utils.AdaptRefernceTimeByUnit(solution.PreparationDate.Time(), chemicalSolTemp.ExpirationPeriodUnit, chemicalSolTemp.ExpirationPeriod))
		nextStandardizationDueDate := stdLog[size-1].Date
		updateData = bson.M{
			"$set": bson.M{
				"standardizationlog":         stdLog,
				"nextstandardizationduedate": nextStandardizationDueDate,
			},
		}
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": solId}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing ChemicalSolutionInventory")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

// helping function
func hasDraft(obj *Models.ChemicalSolutionInventoryRecord) bool {
	size := len(obj.StandardizationLog)
	if size == 0 || obj.StandardizationLog[size-1].FinalConclusion != "Draft" {
		return false
	}
	return true
}

// HasDraft URL/:ChemicalSolutionInventoryId
func HasDraft(c *fiber.Ctx) error {
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	if !hasDraft(&solution) {
		return errors.New("no Draft Found")
	}
	c.Status(200).Send([]byte("found draft"))
	return nil
}

// URL/:ChemicalSolutionInventoryId
func GetDraft(c *fiber.Ctx) error {
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	if !hasDraft(&solution) {
		return errors.New("no Draft Found")
	}
	size := len(solution.StandardizationLog)
	result := solution.StandardizationLog[size-1]
	response, _ := json.Marshal(result) // Decode
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

// URL/:ChemicalSolutionInventroyId and Body have a StandardizationLogEvent obj
func SaveDraft(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	size := len(solution.StandardizationLog)
	var self Models.StandardizationLogEvent
	c.BodyParser(&self)
	self.FinalConclusion = "Draft"
	stdLog := solution.StandardizationLog
	// Append New StandardLogicEvent
	if !hasDraft(&solution) {
		stdLog = append(stdLog, self)
	}else{
		// Update StandardLogicEvent
		stdLog[size-1] = self
	}
	updateData := bson.M{
		"$set": bson.M{
			"standardizationlog": stdLog,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": solId}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Chemical Solution Inventory stdLogEvent")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

// Create New StdLogEvent URL/:ChemicalSolutionInventoryId and Body have a New StandardizationLogEvent obj
func CreateNewStdLogEvent(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ChemicalSolutionInventoryRecord
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	var self Models.StandardizationLogEvent
	c.BodyParser(&self)
	stdLog := solution.StandardizationLog
	// Check if last is Draft
	if hasDraft(&solution) {
		return errors.New("can't create new one while draft is found, please validate & Save Draft first")
	}
	//validate the new StandardLogicEvent
	err = validateResults(&self)
	if err != nil {
		// Fail
		self.FinalConclusion = "Fail"
	} else {
		// Pass and Update Time
		self.Date = primitive.NewDateTimeFromTime(time.Now())
		self.FinalConclusion = "Pass"
	}
	// update the StdLogicEvent
	stdLog = append(stdLog, self)
	updateData := bson.M{
		"$set": bson.M{
			"standardizationlog": stdLog,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": solId}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing ChemicalSolutionInventory")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

// GetHistory URL/:ChemicalSolutionInventoryId
func GetHistory(c *fiber.Ctx) error {
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	stdLog := solution.StandardizationLog
	var reversedStdLog []Models.StandardizationLogEvent
	size := len(stdLog)
	for i := size - 1; i > 0; i-- {
		reversedStdLog = append(reversedStdLog, stdLog[i])
	}
	response, _ := json.Marshal(reversedStdLog)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

// GetNextDueDate URL/:ChemicalSolutionInventoryId
func GetNextDueDate(c *fiber.Ctx) error {
	solId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	solution, err := ChemicalSolutionInventoryRecordGetById(solId)
	if err != nil {
		return err
	}
	response, _ := json.Marshal(solution.NextStandardizationDueDate)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
