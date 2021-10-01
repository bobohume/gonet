package model


type(
	SimplePlayerData struct{
		table string	`sql:"table;name:tbl_player"`
		AccountId int64 `sql:"name:account_id"`
		PlayerId int64 `sql:"primary;name:player_id"`
		PlayerName string `sql:"name:player_name"`
		Level int `sql:"name:level"`
		Sex	  int `sql:"name:sex"`
		Gold  int `sql:"froce;name:gold"`
		DrawGold int `sql:"name:draw_gold"`
		Vip int `sql:"name:vip"`
		LastLogoutTime int64 `sql:"datetime;name:last_logout_time"`
		LastLoginTime int64	`sql:"datetime;name:last_login_time"`
	}

	PlayerKvData struct {
		table 	  string	`sql:"table;name:tbl_player_kv"`
		PlayerId int64     `sql:"primary;name:player_id"`
		Key      int 	   `sql:"primary;name:key"`
		Value    int64     `sql:"name:value"`
	}
)
