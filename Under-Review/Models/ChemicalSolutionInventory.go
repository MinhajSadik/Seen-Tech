package Models

import (
	"fmt"
	"reflect"
	"strings"

	"example.com/seen-eg/CMS/Utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StandardizationLogEvent struct {
	Config          Standardization    `json:"config,omitempty"`  // {embedded} Standardization
	Results         []float64          `json:"results,omitempty"` //Len MUST = Standardization.ResultsCount <- inputs
	Mean            float64            `json:"mean,omitempty"`
	StandardDiv     float64            `json:"standarddiv,omitempty"`
	RSD             float64            `json:"rsd,omitempty"`
	Recovery        float64            `json:"recovery,omitempty"`                   //actual(avg results) / expected
	Date            primitive.DateTime `json:"date,omitempty" bson:"date,omitempty"` //Update Date case finalConc: Pass Only
	FinalConclusion string             `json:"finalconclusion,omitempty"`            // "Draft","Pass", "Fail"
}

func (obj StandardizationLogEvent) GetModifcationBSONObj() bson.M {
	self := bson.M{
		"config":          obj.Config,
		"results":         obj.Results,
		"mean":            obj.Mean,
		"standarddiv":     obj.StandardDiv,
		"rsd":             obj.RSD,
		"recovery":        obj.Recovery,
		"date":            obj.Date,
		"finalconclusion": obj.FinalConclusion,
	}
	return self
}

func (obj StandardizationLogEvent) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.FinalConclusion, validation.In("Draft", "Pass", "Fail")),
	)
}

type Standardization struct {
	Name                  string  `json:"name,omitempty"`
	Description           string  `json:"description,omitempty"`
	Target                float64 `json:"target,omitempty"`
	TargetNumericType     string  `json:"targetnumerictype,omitempty"` //"abs" / "percentage"
	TargetNumberOfDigits  int     `json:"targetnumberofdigits,omitempty"`
	ToleranceIsUsed       bool    `json:"toleranceisused,omitempty"`
	Tolerance             float64 `json:"tolerance,omitempty"`
	ToleranceType         string  `json:"tolerancetype,omitempty"` //"abs" / "percentage"
	ExpirationPeriodInDay int     `json:"expirationperiodinday,omitempty"`
	ResultsCount          int     `json:"resultscount,omitempty"`
	RSDValidationType     string  `json:"rsdvalidationtype,omitempty"` //("Min", "Max", "MinMax", "TargetTolerance)
	RSDMin                float64 `json:"rsdmin,omitempty"`
	RSDMax                float64 `json:"rsdmax,omitempty"`
	RSDTarget             float64 `json:"rsdtarget,omitempty"`
	RSDTolerance          float64 `json:"rsdtolerance,omitempty"`
	Status                bool    `json:"status,omitempty"`
}

func (obj Standardization) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Target, validation.Required),
		validation.Field(&obj.TargetNumericType, validation.Required, validation.In("Abs", "Percentage")),
		validation.Field(&obj.ToleranceType, validation.Required, validation.In("Abs", "Percentage")),
		validation.Field(&obj.RSDValidationType, validation.Required, validation.In("Min", "Max", "MinMax", "TargetTolerance")),
	)
}

type ChemicalSolutionInventoryRecord struct {
	ID                          primitive.ObjectID        `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID                        string                    `json:"uuid,omitempty"`
	PreparationRef              primitive.ObjectID        `json:"preparationref,omitempty" bson:"preparationref,omitempty"`
	ChemicalSolutionTemplateRef primitive.ObjectID        `json:"chemicalsolutiontemplateref,omitempty" bson:"chemicalsolutiontemplateref,omitempty"`
	PreparationDate             primitive.DateTime        `json:"preparationdate,omitempty" bson:"preparationdate,omitempty"`
	ExpirationDate              primitive.DateTime        `json:"expirationdate,omitempty" bson:"expirationdate,omitempty"`
	BatchNumber                 string                    `json:"batchnumber,omitempty"`
	UserRef                     primitive.ObjectID        `json:"userref,omitempty" bson:"userref,omitempty"`
	DesiredAmount               float64                   `json:"desiredamount,omitempty"`
	DesiredAmountUoM            primitive.ObjectID        `json:"desiredamountuom,omitempty" bson:"desiredamountuom,omitempty"`
	CodeImage                   string                    `json:"codeimage,omitempty"`
	ActualAmount                float64                   `json:"actualamount,omitempty"`
	Status                      string                    `json:"status,omitempty"`
	ToSiteRef                   primitive.ObjectID        `json:"tositeref,omitempty" bson:"tositeref,omitempty"`
	ToLabRef                    primitive.ObjectID        `json:"tolabref,omitempty" bson:"tolabref,omitempty"`
	ToAreaRef                   primitive.ObjectID        `json:"toarearef,omitempty" bson:"toarearef,omitempty"`
	ToLocationRef               primitive.ObjectID        `json:"tolocationref,omitempty" bson:"tolocationref,omitempty"`
	StandardizationConfig       Standardization           `json:"standardizationConfig,omitempty"` //{embedded}
	StandardizationLog          []StandardizationLogEvent `json:"standardizationlog,omitempty"`    //{embedded}
	NextStandardizationDueDate  primitive.DateTime        `json:"nextstandardizationduedate,omitempty" bson:"nextstandardizationduedate,omitempty"`
}

type ChemicalSolutionInventoryRecordPopulated struct {
	ID                          primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID                        string                      `json:"uuid,omitempty"`
	PreparationRef              ChemicalSolutionPreparation `json:"preparationref,omitempty" bson:"preparationref,omitempty"`
	ChemicalSolutionTemplateRef ChemicalSolutionTemplate    `json:"chemicalsolutiontemplateref,omitempty" bson:"chemicalsolutiontemplateref,omitempty"`
	PreparationDate             primitive.DateTime          `json:"preparationdate,omitempty" bson:"preparationdate,omitempty"`
	ExpirationDate              primitive.DateTime          `json:"expirationdate,omitempty" bson:"expirationdate,omitempty"`
	BatchNumber                 string                      `json:"batchnumber,omitempty"`
	UserRef                     User                        `json:"userref,omitempty" bson:"userref,omitempty"`
	DesiredAmount               float64                     `json:"desiredamount,omitempty"`
	DesiredAmountUoM            UnitOfMeasurement           `json:"desiredamountuom,omitempty" bson:"desiredamountuom,omitempty"`
	CodeImage                   string                      `json:"codeimage,omitempty"`
	ActualAmount                float64                     `json:"actualamount,omitempty"`
	Status                      string                      `json:"status,omitempty"`
	ToSiteRef                   LabDivision                 `json:"tositeref,omitempty" bson:"tositeref,omitempty"`
	ToLabRef                    LabDivision                 `json:"tolabref,omitempty" bson:"tolabref,omitempty"`
	ToAreaRef                   LabDivision                 `json:"toarearef,omitempty" bson:"toarearef,omitempty"`
	ToLocationRef               LabDivision                 `json:"tolocationref,omitempty" bson:"tolocationref,omitempty"`
	StandardizationConfig       Standardization             `json:"standardizationConfig,omitempty"` //{embedded}
	StandardizationLog          []StandardizationLogEvent   `json:"standardizationlog,omitempty"`    //{embedded}
	NextStandardizationDueDate  primitive.DateTime          `json:"nextstandardizationduedate,omitempty" bson:"nextstandardizationduedate,omitempty"`
}

type ChemicalSolutionInventoryAggregatedAmounts struct {
	ChemicalSolutionRef ChemicalSolutionTemplatePopulated `json:"chemicalsolutionref,omitempty"`
	TotalAmount         float64                           `json:"totalamount,omitempty"`
	Records             []ChemicalSolutionInventoryRecord `json:"records,omitempty"`
}

type ChemicalSolutionInventoryAggregationResult struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TotalAmount float64            `json:"totalamount,omitempty" bson:"totalamount,omitempty"`
}

func (obj *ChemicalSolutionInventoryRecordPopulated) CloneFrom(other ChemicalSolutionInventoryRecord) {
	obj.ID = other.ID
	obj.UUID = other.UUID
	obj.PreparationRef = ChemicalSolutionPreparation{}
	obj.ChemicalSolutionTemplateRef = ChemicalSolutionTemplate{}
	obj.PreparationDate = other.PreparationDate
	obj.ExpirationDate = other.ExpirationDate
	obj.BatchNumber = other.BatchNumber
	obj.UserRef = User{}
	obj.DesiredAmount = other.DesiredAmount
	obj.DesiredAmountUoM = UnitOfMeasurement{}
	obj.CodeImage = other.CodeImage
	obj.ActualAmount = other.ActualAmount
	obj.Status = other.Status
	obj.ToSiteRef = LabDivision{}
	obj.ToLabRef = LabDivision{}
	obj.ToAreaRef = LabDivision{}
	obj.ToLocationRef = LabDivision{}
	obj.StandardizationConfig = other.StandardizationConfig
	obj.StandardizationLog = other.StandardizationLog
	obj.NextStandardizationDueDate = other.NextStandardizationDueDate
}

func (obj ChemicalSolutionInventoryRecord) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Status, validation.In("Active", "Inactive", "Scraped", "Consumed")),
	)
}

func (obj ChemicalSolutionInventoryRecord) GetIdString() string {
	return obj.ID.String()
}

func (obj ChemicalSolutionInventoryRecord) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj ChemicalSolutionInventoryRecord) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()

	invalidFieldNames := []string{"ID", "UUID", "CodeImage"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}

	return self
}

type ChemicalSolutionInventoryRecordSearch struct {
	ReferencesKeys    []string             `json:"referenceskeys"`
	ReferencesValues  []primitive.ObjectID `json:"referencesvalues" bson:"referencesvalues,omitempty"`
	ReferencesAreUsed bool                 `json:"referencesareused"`

	StringsKeys    []string `json:"stringskeys"`
	StringsValues  []string `json:"stringsvalues"`
	StringsAreUsed bool     `json:"stringsareused"`

	BooleansKeys    []string `json:"booleanskeys"`
	BooleansValues  []bool   `json:"booleansvalues"`
	BooleansAreUsed bool     `json:"booleansareused"`
}

func (obj ChemicalSolutionInventoryRecordSearch) GetSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.ReferencesAreUsed {
		for i, v := range obj.ReferencesKeys {
			self[v] = obj.ReferencesValues[i]
		}
	}

	if obj.StringsAreUsed {
		for i, v := range obj.StringsKeys {
			regexPattern := fmt.Sprintf(".*%s.*", obj.StringsValues[i])
			self[v] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
		}
	}

	if obj.BooleansAreUsed {
		for i, v := range obj.BooleansKeys {
			self[v] = obj.BooleansValues[i]
		}
	}

	return self
}

type ChemicalSolutionInventoryRecordOperationOpen struct {
	IsNewExpirationDate bool               `json:"isnewexpirationdate,omitempty"`
	NewExpirationDate   primitive.DateTime `json:"newexpirationdate,omitempty" bson:"newexpirationdate,omitempty"`
	StorageConditionRef primitive.ObjectID `json:"storageconditionref,omitempty"`
	ToSiteRef           primitive.ObjectID `json:"tositeref,omitempty" bson:"tositeref,omitempty"`
	ToLabRef            primitive.ObjectID `json:"tolabref,omitempty" bson:"tolabref,omitempty"`
	ToAreaRef           primitive.ObjectID `json:"toarearef,omitempty" bson:"toarearef,omitempty"`
	ToLocationRef       primitive.ObjectID `json:"tolocationref,omitempty" bson:"tolocationref,omitempty"`
	ApprovalCycle       ApprovalCycle      `json:"approvalCycle,omitempty" bson:"tolocationref,omitempty"`
}
