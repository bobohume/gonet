package model

type(
	SocialItem struct {
		table 	 string	`sql:"table;name:tbl_social"`
		PlayerId int64	`sql:"primary;name:player_id"`
		TargetId int64	`sql:"primary;name:target_id"`
		Type	int8	`sql:"primary;name:type"`
		FriendValue	int `sql:"name:friend_value"`
	}
)