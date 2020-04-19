package db_worker

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

// Значение события == event_id в таблице events
type Event int
const (
	Subscribe Event = iota + 1 // event_id = 1
	Like // event_id = 2 и т.д.
	Repost
	Comment
)

// Обновляет список подписчиков в базе данных
func UpdateSubscribers(updatedSubs []int, groupId string){
	db, err := sql.Open("mysql", "root:@tcp(localhost:3308)/vk_activity")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM users_activities WHERE event_id=?", Subscribe)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var savedSubs []int
	for rows.Next(){
		var userId int
		rows.Scan(&userId)
		savedSubs = append(savedSubs, userId)
	}

	// Удаляем из базы отписавшихся людей
	sqlStr, values := generateSqlDeleteUnsubs(savedSubs, updatedSubs)
	if values != nil{
		stmt,err := db.Prepare(sqlStr)
		if err != nil{
			log.Fatal(err)
		}
		_, err = stmt.Exec(values...)
		if err != nil{
			log.Fatal(err)
		}
	}

	// Добавляем в базу новых подписчиков
	sqlStr, values = generateSqlInsertSubs(savedSubs,updatedSubs, groupId)
	if values != nil{
		stmt, err := db.Prepare(sqlStr)
		if err != nil{
			log.Fatal(err)
		}
		_, err = stmt.Exec(values...)
		if err != nil{
			log.Fatal(err)
		}
	}
}

// Генерирует SQL-запрос на удаление отписавшихся
// Возвращает запрос в виде строки и список параметров
func generateSqlDeleteUnsubs(savedSubs, updatedSubs []int) (sqlStr string, values []interface{}){
	sqlStr = "DELETE FROM users_activities WHERE user_id IN (%s)"
	var placeholders string
	for _, id := range savedSubs {
		if !contains(updatedSubs, id){
			values = append(values, id)
			placeholders += "?,"
		}
	}
	placeholders = strings.TrimSuffix(placeholders, ",")
	sqlStr = fmt.Sprintf(sqlStr,placeholders)
	return
}

// Генерирует SQL-запрос на добавление новых подписчиков
// Возвращает запрос в виде строки и список параметров
func generateSqlInsertSubs(savedSubs, updatedSubs []int, groupId string) (sqlStr string, values []interface{}){
	sqlStr = "INSERT INTO users_activities(user_id,group_id,event_id) VALUES "
	var inserts []string
	const placeholders = "(?,?,?)"
	for _, id := range updatedSubs{
		if !contains(savedSubs, id){
			inserts = append(inserts, placeholders)
			values = append(values, id, groupId, Subscribe)
		}
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	return
}

// Проверяет, содержится ли элемент value в массиве arr
func contains(arr []int, value int) bool{
	for _, item := range arr {
		if item == value{
			return true
		}
	}
	return false
}