package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID           string             `json:"uuid,omitempty"`
	Type           string             `json:"type,omitempty" binding:"required"`
	IsCompany      bool               `json:"iscompany,omitempty"`
	CompanyRef     primitive.ObjectID `json:"companyref,omitempty"`
	Status         bool               `json:"status,omitempty"`
	Name           string             `json:"name,omitempty" binding:"required"`
	Mobile1        string             `json:"mobile1,omitempty"`
	Mobile2        string             `json:"mobile2,omitempty"`
	LandlineNumber string             `json:"landlinenumber,omitempty"`
	WhatsApp       string             `json:"whatsapp,omitempty"`
	Website        string             `json:"website,omitempty"`
	Email          string             `json:"email,omitempty"`
	Address        string             `json:"address,omitempty"`
	Notes          string             `json:"notes,omitempty"`
}

func (obj Contact) Validate() error {
	return validation.ValidateStruct(&obj,
		// Type Can't be empty, and it has two variables only
		validation.Field(&obj.Type, validation.Required, validation.In("Vendor", "Supplier")),
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&obj.Name, validation.Required, validation.Length(2, 50)),
	)

}

func (obj Contact) GetIdString() string {
	return obj.ID.String()
}

func (obj Contact) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Contact) GetModifcationBSONObj() bson.M {
	self := bson.M{
		"_id":            obj.ID,
		"type":           obj.Type,
		"iscompany":      obj.IsCompany,
		"companyref":     obj.CompanyRef,
		"status":         obj.Status,
		"name":           obj.Name,
		"mobile1":        obj.Mobile1,
		"mobile2":        obj.Mobile2,
		"landlinenumber": obj.LandlineNumber,
		"whatsApp":       obj.WhatsApp,
		"website":        obj.Website,
		"email":          obj.Email,
		"address":        obj.Address,
		"notes":          obj.Notes,
	}

	return self
}

type ContactSearch struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	IDIsUsed        bool               `json:"idisused"`
	Name            string             `json:"name,omitempty"`
	NameIsUsed      bool               `json:"nameisused,omitempty"`
	Status          bool               `json:"status,omitempty"`
	StatusIsUsed    bool               `json:"statusisused,omitempty"`
	Type            string             `json:"type,omitempty"`
	TypeIsUsed      bool               `json:"typeisused,omitempty"`
	IsCompany       bool               `json:"iscompany,omitempty"`
	IsCompanyIsUsed bool               `json:"iscompanyisused,omitempty"`
	TextData        string             `json:"textdata,omitempty"`
	IsTextDataUsed  bool               `json:"istextdataused,omitempty"`
}

func (obj ContactSearch) GetContactSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.IsCompanyIsUsed {
		self["iscompany"] = obj.IsCompany
	}

	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}

	if obj.IsTextDataUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.TextData)

		self["$or"] = []bson.M{
			{"name": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"mobile1": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"mobile2": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"landlinenumber": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"whatsapp": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"website": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"email": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"address": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
			{"notes": bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}},
		}
	}

	return self
}

type ContactPopulated struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID           string             `json:"uuid,omitempty"`
	Type           string             `json:"type,omitempty" binding:"required"`
	IsCompany      bool               `json:"iscompany,omitempty"`
	CompanyRef     Contact            `json:"companyref,omitempty"`
	Status         bool               `json:"status,omitempty"`
	Name           string             `json:"name,omitempty" binding:"required"`
	Mobile1        string             `json:"mobile1,omitempty"`
	Mobile2        string             `json:"mobile2,omitempty"`
	LandlineNumber string             `json:"landlinenumber,omitempty"`
	WhatsApp       string             `json:"whatsapp,omitempty"`
	Website        string             `json:"website,omitempty"`
	Email          string             `json:"email,omitempty"`
	Address        string             `json:"address,omitempty"`
	Notes          string             `json:"notes,omitempty"`
}

func (obj *ContactPopulated) CloneFrom(other Contact) {
	obj.ID = other.ID
	obj.UUID = other.UUID
	obj.Type = other.Type
	obj.IsCompany = other.IsCompany
	obj.CompanyRef = Contact{}
	obj.Status = other.Status
	obj.Name = other.Name
	obj.Mobile1 = other.Mobile1
	obj.Mobile2 = other.Mobile2
	obj.LandlineNumber = other.LandlineNumber
	obj.WhatsApp = other.WhatsApp
	obj.Website = other.Website
	obj.Email = other.Email
	obj.Address = other.Address
	obj.Notes = other.Notes
}

type ContactAggregated struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	CompanyData Contact            `json:"companydata"`
	Records     []Contact          `json:"records,omitempty"`
}
