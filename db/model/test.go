package model

import "gonet/server/message"

type (
	SimplePlayerData struct {
		table          string              `sql:"table;name:tbl_player"`
		AccountId      int64               `sql:"primary;name:account_id"`
		PlayerId       int64               `sql:"primary;name:player_id"`
		PlayerName     string              `sql:"name:player_name"`
		Level          int                 `sql:"name:level"`
		Sex            int                 `sql:"name:sex"`
		Gold           int                 `sql:"name:gold"`
		DrawGold       int                 `sql:"name:draw_gold"`
		Vip            [8]int              `sql:"name:vip"`
		LastLogoutTime int64               `sql:"datetime;nameg:last_logout_time"`
		LastLoginTime  int64               `sql:"datetime;name:last_login_time"`
		PLayerBlob     *message.PlayerData `sql:"blob;name:plaeyr_blob"`
		PLayerBlobJson *AA                 `sql:"json;name:plaeyr_blob_json"`
	}

	AA struct {
		A int
		B map[int]string
	}
)
