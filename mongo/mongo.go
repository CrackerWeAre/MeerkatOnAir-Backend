package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	m "github.com/ssoyyoung.p/GoDirectory/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database Name
const dbName = "meerkatonair"

// Collection name
const colNameLive = "live_list"
const colNameCrawl = "crawl_target"
const colNameUser = "user_info"

// get Collection
func getCollection(client *mongo.Client, colName string) *mongo.Collection {
	return client.Database(dbName).Collection(colName)
}

// get MongoDB Authorization info7
func getAuth() m.Auth {
	data, err := os.Open("mongo/mongodb_auth.json")
	checkErr(err)

	var auth m.Auth
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &auth)

	return auth
}

// checkErr function
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// connect to MongoDB
func connectDB() (client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	// Timeout 설정을 위한 Context생성
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	Authrization := getAuth()

	// Auth에러 처리를 위한 client option 구성
	clientOptions := options.Client().ApplyURI("mongodb://" + Authrization.Hostname + Authrization.Port).SetAuth(options.Credential{
		Username: Authrization.Username,
		Password: Authrization.Password,
	})

	// MongoDB 연결
	client, err := mongo.Connect(ctx, clientOptions)
	checkErr(err)

	// MongoDB 연결 검증
	checkErr(client.Ping(ctx, readpref.Primary()))

	return client, ctx, cancel
}

// define bson.M type data
var datas []bson.M

// jsonData marshal to string func
func jsonMarshalString(datas []bson.M) string {
	jsonBytes, err := json.Marshal(datas)
	checkErr(err)
	jsonString := string(jsonBytes)

	return jsonString
}

// string ID convert to OjectID
func convertID(id string) primitive.ObjectID {
	docID, err := primitive.ObjectIDFromHex(id)
	checkErr(err)

	return docID
}

// LiveList func
func LiveTrueList() string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"liveAttdc", -1}})

	res, err := getCollection(client, colNameLive).Find(ctx, bson.M{"onLive": true}, findOptions)
	checkErr(err)

	if err = res.All(ctx, &datas); err != nil {
		fmt.Println(err)
	}

	return jsonMarshalString(datas)
}

// AllList func
func LiveAllList() string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	res, err := getCollection(client, colNameCrawl).Find(ctx, bson.M{})
	checkErr(err)

	if err = res.All(ctx, &datas); err != nil {
		fmt.Println(err)
	}

	return jsonMarshalString(datas)
}

// SearchDBbyID func
func SearchDBbyID(id string) string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	res, _ := getCollection(client, colNameCrawl).Find(ctx, bson.M{"_id": convertID(id)})
	if err := res.All(ctx, &datas); err != nil {
		fmt.Println(err)
	}

	return jsonMarshalString(datas)
}

// DeleteDBbyID func
func DeleteDBbyID(id string) string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	_, err := getCollection(client, colNameCrawl).DeleteOne(ctx, bson.M{"_id": convertID(id)})
	checkErr(err)

	return "Delete!"
}

//UpdateDBbyID func
func UpdateDBbyID(id, platform, channel, channelID string) string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	filter := bson.M{"_id": convertID(id)}
	update := bson.D{
		{"$set", bson.D{
			{"platform", platform},
			{"channel", channel},
			{"channelID", channelID},
		},
		},
	}
	_, err := getCollection(client, colNameCrawl).UpdateOne(ctx, filter, update)
	checkErr(err)

	return "Update!"
}

// CreateDB func
func CreateDB(platform, channel, channelID string) string {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	newData := m.CrawlTarget{
		Platform:  platform,
		Channel:   channel,
		ChannelID: channelID,
	}
	_, err := getCollection(client, colNameCrawl).InsertOne(ctx, newData)
	checkErr(err)

	return "create!"
}

//CheckUser func
func CheckUser(googleID, name, email string) bool {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	res, err := getCollection(client, colNameUser).CountDocuments(ctx, bson.M{"googleId": googleID, "name": name})
	checkErr(err)
	if res == 0 {
		createUser(googleID, name, email)
	}
	return true
}

// createUser func
func createUser(googleID, name, email string) {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	_, err := getCollection(client, colNameUser).InsertOne(ctx, bson.M{
		"googleId": googleID,
		"name":     name,
		"email":    email,
	})
	checkErr(err)
}

// UpdateUser
func UpdateUser(googleID, token string) {
	client, ctx, cancel := connectDB()
	defer client.Disconnect(ctx)
	defer cancel()

	filter := bson.M{"googleId": googleID}
	update := bson.D{
		{"$set", bson.D{
			{"token", token},
		},
		},
	}
	res, err := getCollection(client, colNameUser).UpdateOne(ctx, filter, update)
	fmt.Println(res)
	checkErr(err)
}
