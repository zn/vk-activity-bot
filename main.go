package main

import (
	"context"
	db_worker "github.com/zn/vk-activity-bot/db-worker"
	vk "github.com/zn/vk-activity-bot/vk-api-handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const GROUP_ID string = "siberianfederal"
const API_VERSION string = "5.103"
const API_TOKEN string = ""

func main(){
	vkHandler := vk.VkApiHandler{
		GroupId: GROUP_ID,
		Version:  API_VERSION,
		ApiToken: API_TOKEN,
		RequestsPerSecond: 3,
	}
	allMembers := vkHandler.GetMembersList()
	db := InitMongoClient()
	defer db.Client.Disconnect(db.Context)
	db.UpdateSubscribers(allMembers, GROUP_ID)
}

func InitMongoClient() db_worker.DbWorker{
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return db_worker.DbWorker{Client: client, Context: ctx}
}

