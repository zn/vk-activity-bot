package dbworker

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Значение события == event_id в таблице events
type Event int

const (
	Subscribe Event = iota + 1 // event_id = 1
	Like                       // event_id = 2 и т.д.
	Repost
	Comment
)

type DbWorker struct {
	Client  *mongo.Client
	Context context.Context
}

// Database collection users_activities
type userActivity struct {
	UserId  int       `bson:"user_id"`
	GroupId string    `bson:"group_id"`
	EventId Event     `bson:"event_id"`
	Date    time.Time `bson:"date"`
}

// Возвращает все записи коллекции, удовлетворяющие фильтру
func getAll(collection *mongo.Collection, ctx context.Context, filter bson.M, readTo interface{}) error {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, readTo); err != nil {
		return err
	}
	return nil
}
