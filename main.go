package main

import (
	"context"
	"log"

	"github.com/zn/vk-activity-bot/dbworker"
	"github.com/zn/vk-activity-bot/vk"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GROUP_ID string = "englishyo" //"mudakoff"
const API_VERSION string = "5.103"
const API_TOKEN string = ""

func main() {
	vkHandler := vk.ApiHandler{
		GroupId:           GROUP_ID,
		Version:           API_VERSION,
		ApiToken:          API_TOKEN,
		RequestsPerSecond: 3,
	}
	allMembers, err := vkHandler.GetMembersList()
	if err != nil {
		log.Fatal(err)
	}
	db := InitMongoClient()
	defer db.Client.Disconnect(db.Context)
	err = db.UpdateSubscribers(allMembers, GROUP_ID)
	if err != nil {
		log.Fatal(err)
	}
}

func InitMongoClient() dbworker.DbWorker {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return dbworker.DbWorker{Client: client, Context: ctx}
}
