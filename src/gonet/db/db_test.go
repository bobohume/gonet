package db_test

import (
	"fmt"
	"gonet/db"
	"gonet/message"
	"testing"
)

type (

	AA struct {
		A int
		B map[int] string
	}

	SimplePlayerData struct{
		AccountId int64 `sql:"primary;name:account_id"`
		PlayerId int64 `sql:"primary;name:player_id"`
		PlayerName string `sql:"name:player_name"`
		Level int `sql:"name:level"`
		Sex	  int `sql:"name:sex"`
		Gold  int `sql:"name:gold"`
		DrawGold int `sql:"name:draw_gold"`
		Vip [8]int `sql:"name:vip"`
		LastLogoutTime int64 `sql:"datetime;nameg:last_logout_time"`
		LastLoginTime int64	`sql:"datetime;name:last_login_time"`
		PLayerBlob *message.PlayerData	`sql:"blob;name:plaeyr_blob"`
		PLayerBlobJson *AA	`sql:"json;name:plaeyr_blob_json"`
	}
)


func TestInsert(t *testing.T)  {
	nameMap := make(map[string]  bool)
	//nameMap["1"][2] = true
	nVal := nameMap["1"]
	fmt.Println(nVal)
	data := &SimplePlayerData{PLayerBlob:&message.PlayerData{}, PLayerBlobJson:&AA{A:1, B: map[int]string{1:"test", 2:"test2"}}, Vip:[8]int{1,2,3,4,5,6,7,8}}
	t.Log(db.InsertSql(data, "tbl_player"))
	t.Log(db.InsertSqlEx(data, "tbl_player", "account_id", "gold", "plaeyr_blob", "plaeyr_blob_json", "vip0", "vip7"))
}


func TestUpdate(t *testing.T)  {
	data := &SimplePlayerData{PLayerBlob:&message.PlayerData{}, PLayerBlobJson:&AA{A:1, B: map[int]string{1:"test", 2:"test2"}}}
	t.Log(db.UpdateSql(data, "tbl_player"))
	t.Log(db.UpdateSqlEx(data, "tbl_player", "gold", "plaeyr_blob", "plaeyr_blob_json", "vip7"))
}

func TestLoad(t *testing.T)  {
	data := &SimplePlayerData{PLayerBlob:&message.PlayerData{}, PLayerBlobJson:&AA{A:1, B: map[int]string{1:"test", 2:"test2"}}}
	t.Log(db.LoadSql(data, "tbl_player", "where player_id = 0"))
	t.Log(db.LoadSqlEx(data, "tbl_player", "where player_id = 0", "account_id", "gold", "plaeyr_blob", "plaeyr_blob_json", "vip5"))
}

func TestDelete(t *testing.T)  {
	data := &SimplePlayerData{PLayerBlob:&message.PlayerData{}, PLayerBlobJson:&AA{A:1, B: map[int]string{1:"test", 2:"test2"}}}
	t.Log(db.DeleteSql(data, "tbl_player"))
	t.Log(db.DeleteSqlEx(data, "tbl_player", "player_id", "vip0"))
}