package dbworker

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Обновляет список подписчиков в базе данных
func (db DbWorker) UpdateSubscribers(updatedSubsList []int, groupId string) error {
	if len(updatedSubsList) == 0 {
		return errors.New("subscribers list is empty")
	}
	collection := db.Client.Database("vk_activity").Collection("users_activities")

	// удаляем отписавшихся
	deleteResult, err := collection.DeleteMany(db.Context, bson.M{
		"group_id": groupId,
		"user_id": bson.D{
			{"$nin", updatedSubsList},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("Deleted:", deleteResult.DeletedCount)

	// добавляем новых подписчиков
	var oldSubs []userActivity
	if err = getAll(collection, db.Context, bson.M{"group_id": groupId}, &oldSubs); err != nil {
		return err
	}
	newSubsMap := getNewSubscribers(updatedSubsList, oldSubs)
	if len(newSubsMap) == 0 {
		fmt.Println("New subscribers: 0")
		return nil
	}
	itemsToInsert := createDocuments(newSubsMap, groupId)
	insertResult, err := collection.InsertMany(db.Context, itemsToInsert)
	if err != nil {
		return err
	}
	fmt.Println("New subscribers:", len(insertResult.InsertedIDs))
	return nil
}

// Возвращает map, где ключами являются IDs новых подписчиков
func getNewSubscribers(updatedSubs []int, oldSubs []userActivity) map[int]*struct{} {
	// Переносим из updatedUsers в map, где ключами являются IDs всех подписчиков
	newSubsMap := map[int]*struct{}{}
	empty := struct{}{}
	for _, userId := range updatedSubs {
		newSubsMap[userId] = &empty
	}

	// Из мапы удаляем все те IDs, которые уже есть в бд.
	for _, item := range oldSubs {
		delete(newSubsMap, item.UserId)
	}
	return newSubsMap
}

// Возвращает список документов для добавления в базу данных
func createDocuments(sourceMap map[int]*struct{}, groupId string) []interface{} {
	var itemsToInsert []interface{} // []userActivity
	timeNow := time.Now()
	for userId, _ := range sourceMap {
		itemsToInsert = append(itemsToInsert, userActivity{
			UserId:  userId,
			GroupId: groupId,
			EventId: Subscribe,
			Date:    timeNow,
		})
	}
	return itemsToInsert
}
