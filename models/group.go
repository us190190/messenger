package models

import (
	"database/sql"
	"fmt"
	"github.com/us190190/messenger/database"
	"log"
)

type Group struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
}

type GroupMember struct {
	ID                 int `json:"id"`
	GroupID            int `json:"group_id"`
	UserID             int `json:"user_id"`
	LastDeliveredMsgID int `json:"last_delivered_msg_id"`
}

func GetGrpMbrsByGroupID(groupID int) ([]GroupMember, error) {
	var grpMembers []GroupMember
	db := database.GetDB()
	qry := fmt.Sprintf("SELECT id, group_id, user_id, last_delivered_msg_id "+
		"FROM group_members WHERE group_id = %d", groupID)
	rows, err := db.Query(qry)
	if err != nil {
		log.Println(fmt.Sprintf("GetGrpMbrsByGroupID failed qry: %s error: %v\n", qry, err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// log  fatal error
		}
	}(rows)

	for rows.Next() {
		var curGrpMember GroupMember
		err := rows.Scan(&curGrpMember.ID, &curGrpMember.GroupID, &curGrpMember.UserID, &curGrpMember.LastDeliveredMsgID)
		if err != nil {
			return nil, err
		}
		grpMembers = append(grpMembers, curGrpMember)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grpMembers, err
}

func UpdateLastMsgDelvrdUsrInGrp(userID int, groupID int, msgID int) (bool, error) {
	db := database.GetDB()
	qry := fmt.Sprintf("UPDATE group_members SET last_delivered_msg_id = %d WHERE user_id = %d AND group_id = %d", msgID, userID, groupID)
	_, err := db.Exec(qry)
	if err != nil {
		log.Println(fmt.Sprintf("UpdateLastMsgDelvrdUsrInGrp failed qry: %s error: %v\n", qry, err))
		return false, err
	}
	return true, nil
}
