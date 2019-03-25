package player

type (
	SimplePlayerData struct{
		AccountId int64 `sql:"name:account_id"`
		PlayerId int64 `sql:"primary;name:player_id"`
		PlayerName string `sql:"name:player_name"`
		Level int `sql:"name:level"`
		Sex	  int `sql:"name:sex"`
		Gold  int `sql:"name:gold"`
		DrawGold int `sql:"name:draw_gold"`
		Vip int `sql:"name:vip"`
		LastLogoutTime int64 `sql:"datetime;name:last_logout_time"`
		LastLoginTime int64	`sql:"datetime;name:last_login_time"`
	}
)

//-----load blob---//
/*rows, err := world.SERVER.GetDB().Query("select `blob` from tbl_player where player_id = ?" , pData.PlayerId)
rs := db.Query(rows)
if rs.Next(){
fmt.Println(rs.Row().Get("blob"))
}

//-----set blob-----//
byte, _ := bson.Marshal(pData)
_, err = world.SERVER.GetDB().Exec("update tbl_player set `blob` = ? where player_id = ?", byte, pData.PlayerId)
fmt.Println(err)*/