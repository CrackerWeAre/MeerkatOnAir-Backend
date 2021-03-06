package utils

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JSONMarshalString : jsonData marshal to string func
func JSONMarshalString(datas []bson.M) string {
	jsonBytes, err := json.Marshal(datas)
	CheckErr(err)
	jsonString := string(jsonBytes)

	return jsonString
}

// ConvertID : string ID convert to OjectID
func ConvertID(id string) primitive.ObjectID {
	docID, err := primitive.ObjectIDFromHex(id)
	CheckErr(err)

	return docID
}

// CheckErr function
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
