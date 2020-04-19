package db_worker

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type Event int
const (
	Subscribe Event = iota + 1
	Like
	Repost
	Comment
)

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

	// Удаляем отписавшихся
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

	// Добавляем новых
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

func contains(arr []int, value int) bool{
	for _, item := range arr {
		if item == value{
			return true
		}
	}
	return false
}