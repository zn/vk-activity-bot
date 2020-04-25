package db_worker

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// Значение события == event_id в таблице events
type Event int
const (
	Subscribe Event = iota + 1 // event_id = 1
	Like // event_id = 2 и т.д.
	Repost
	Comment
)

type DbWorker struct{
	Client *mongo.Client
	Context context.Context
}

type userActivity struct{
	Id primitive.ObjectID `bson:"_id"`
	UserId int `bson:"user_id"`
	GroupId string `bson:"group_id"`
	EventId Event `bson:"event_id"`
	Date time.Time `bson:"date"`
}

// Обновляет список подписчиков в базе данных
func (db DbWorker) UpdateSubscribers(updatedSubs []int, groupId string) {
	if len(updatedSubs) == 0 {return}
	collection := db.Client.Database("vk_activity").Collection("users_activities")

	var oldSubscribers []userActivity
	cursor, err := collection.Find(db.Context, bson.M{})
	if err != nil{ log.Fatal(err) }
	defer cursor.Close(db.Context)
	if err = cursor.All(db.Context, &oldSubscribers); err != nil{
		log.Fatal(err)
	}

	// удаляем отписавшихся
	var usersToDelete []primitive.ObjectID
	for _, subscriber := range oldSubscribers{
		if !inArray(updatedSubs, subscriber.UserId){
			usersToDelete = append(usersToDelete, subscriber.Id)
		}
	}
	if len(usersToDelete) != 0 {
		deleteResult, err := collection.DeleteMany(db.Context, bson.M{
			"_id": bson.D{{"$in", usersToDelete}},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Unsubscribed:", deleteResult.DeletedCount)
	}

	// добавляем новых подписчиков
	var usersToInsert []interface{}
	for _, id := range updatedSubs{
		if !inDatabase(oldSubscribers, id){
			temp := userActivity{
				Id:      primitive.NewObjectID(),
				UserId:  id,
				GroupId: groupId,
				EventId: Subscribe,
				Date:    time.Now(),
			}
			usersToInsert = append(usersToInsert, temp)
		}
	}
	if len(usersToInsert) != 0{
		insertResult, err := collection.InsertMany(db.Context, usersToInsert)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println("New subscribers:", len(insertResult.InsertedIDs))
	}
}

// Проверяет, содержится ли элемент value в массиве arr
func inArray(arr []int, value int) bool{
	for _, item := range arr {
		if item == value{
			return true
		}
	}
	return false
}

func inDatabase(arr []userActivity, userId int) bool{
	for _, activity := range arr{
		if userId == activity.UserId{
			return true
		}
	}
	return false
}