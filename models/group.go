package models

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
