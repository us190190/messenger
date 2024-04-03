package models

import (
	"database/sql"
	"fmt"
	"github.com/us190190/messenger/database"
	"log"
	"time"
)

type Message struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	ReceiverID  int       `json:"receiver_id"`
	GroupID     int       `json:"group_id"`
	Message     string    `json:"message"`
	IsDelivered bool      `json:"is_delivered"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetUndeliveredPvtMsgsByUserID(receiverID int) ([]Message, error) {
	var undeliveredPvtMsgs []Message
	db := database.GetDB()
	qry := fmt.Sprintf("SELECT id, sender_id, receiver_id, group_id, created_at, updated_at "+
		"FROM messages WHERE receiver_id = %d AND is_delivered = 0 AND group_id = 0 "+
		"ORDER BY created_at", receiverID)
	rows, err := db.Query(qry)
	if err != nil {
		log.Println(fmt.Sprintf("GetUndeliveredPvtMsgsByUserID failed qry: %s error: %v\n", qry, err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// log  fatal error
		}
	}(rows)

	for rows.Next() {
		var curMsg Message
		err := rows.Scan(&curMsg.ID, &curMsg.SenderID, &curMsg.ReceiverID, &curMsg.GroupID, &curMsg.CreatedAt, &curMsg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		undeliveredPvtMsgs = append(undeliveredPvtMsgs, curMsg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return undeliveredPvtMsgs, err
}

func MarkPvtMsgDelivered(ID int) (bool, error) {
	// Insert message into database along with the status whether it was delivered or not
	db := database.GetDB()
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	qry := fmt.Sprintf("UPDATE messages SET is_delivered = 1, updated_at = '%s' where ID = %d", timeNow, ID)
	_, err := db.Exec(qry)
	if err != nil {
		log.Println(fmt.Sprintf("MarkPvtMsgDelivered failed qry: %s error: %v\n", qry, err))
		return false, err
	}
	return true, nil
}

func InsertNewMessage(msg Message) (int, error) {
	db := database.GetDB()
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	isDelivered := 0
	if msg.IsDelivered {
		isDelivered = 1
	} else {
		isDelivered = 0
	}
	result, err := db.Exec(fmt.Sprintf("INSERT INTO messages (sender_id, receiver_id, group_id, message, is_delivered, created_at, updated_at) VALUES (%d, %d, %d, '%s', %d, '%s', '%s')", msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Message, isDelivered, timeNow, timeNow))
	if err != nil {
		return 0, err
	}

	// Get the last insert ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}
