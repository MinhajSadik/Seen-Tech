package Utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckIfObjExistingByObjId(collection *mongo.Collection, objID primitive.ObjectID) error {
	filter := bson.M{"_id": objID}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)
	fmt.Println("Count : ", len(results))

	if len(results) == 0 {
		return errors.New("obj not found")
	}

	return nil
}

func AdaptCurrentTimeByUnit(unit string, period int) time.Time {
	now := time.Now()
	if unit == "Month" {
		now = now.AddDate(0, period, 0)
	} else if unit == "Week" {
		now = now.AddDate(0, 0, period*7)
	} else if unit == "Day" {
		now = now.AddDate(0, 0, period)
	} else if unit == "Year" {
		now = now.AddDate(period, 0, 0)
	}
	return now
}

func AdaptRefernceTimeByUnit(refernceTime time.Time, unit string, period int) time.Time {
	if unit == "Month" {
		refernceTime = refernceTime.AddDate(0, period, 0)
	} else if unit == "Week" {
		refernceTime = refernceTime.AddDate(0, 0, period*7)
	} else if unit == "Day" {
		refernceTime = refernceTime.AddDate(0, 0, period)
	} else if unit == "Year" {
		refernceTime = refernceTime.AddDate(period, 0, 0)
	}
	return refernceTime
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func UploadImage(c *fiber.Ctx) string {
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Failed in saving Image")
		c.Status(500).Send([]byte("Invalid data sent for uploading"))
		return "Error"
	}

	// Save file to root directory
	var filePath = fmt.Sprintf("Resources/Images/img_%d_%d.png", rand.Intn(1024), MakeTimestamp())
	saveing_err := c.SaveFile(file, "./"+filePath)
	if saveing_err != nil {
		c.Status(500).Send([]byte("Failed to save the uploaded image"))
		return "Error"
	} else {
		c.Status(200).Send([]byte("Saved Successfully"))
		return filePath
	}
}

func FindByFilter(collection *mongo.Collection, filter bson.M) (bool, []bson.M) {
	results := []bson.M{}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return true, results
}

func Contains(arr []primitive.ObjectID, elem primitive.ObjectID) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func Unique(inSlice []primitive.ObjectID) []primitive.ObjectID {
	keys := make(map[string]bool)
	var list []primitive.ObjectID
	for _, entry := range inSlice {
		if _, value := keys[entry.Hex()]; !value {
			keys[entry.Hex()] = true
			list = append(list, entry)
		}
	}
	return list
}

func ArrayStringContains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func DecodeArrData(inStructArr, outStructArr interface{}) error {
	in := struct{ Data interface{} }{Data: inStructArr}
	inStructArrData, err := bson.Marshal(in)
	if err != nil {
		return err
	}
	var out struct{ Data bson.Raw }
	if err := bson.Unmarshal(inStructArrData, &out); err != nil {
		return err
	}
	return bson.Unmarshal(out.Data, &outStructArr)
}

func SendTextResponseAsJSON(c *fiber.Ctx, msg string) {
	response, _ := json.Marshal(
		bson.M{"result": msg},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
}

func DateToJulianDay(year int, month time.Month, day int) (julianDay int) {
	if year <= 0 {
		year++
	}
	a := int(14-month) / 12
	y := year + 4800 - a
	m := int(month) + 12*a - 3
	julianDay = int(day) + (153*m+2)/5 + 365*y + y/4
	if year > 1582 || (year == 1582 && (month > time.October || (month == time.October && day >= 15))) {
		return julianDay - y/100 + y/400 - 32045
	} else {
		return julianDay - 32083
	}
}

func GenerateBatchNumber() string {
	var str string
	year := time.Now().Year()
	str = strconv.Itoa(year%10) + "/" + strconv.Itoa(DateToJulianDay(year, time.Now().Month(), time.Now().Day())) + "/"
	switch hr := time.Now().Hour(); {
	case hr >= 7 && hr < 15:
		str += "B"
	case hr >= 15 && hr < 23:
		str += "C"
	default:
		str += "A"
	}
	return str
}
