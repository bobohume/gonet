package db_test

import (
	"fmt"
	"gonet/db"
	"gonet/db/model"
	"gonet/server/message"
	"testing"
)

func TestInsert(t *testing.T) {
	nameMap := make(map[string]bool)
	//nameMap["1"][2] = true
	nVal := nameMap["1"]
	fmt.Println(nVal)
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}, Vip: [8]int{1, 2, 3, 4, 5}}
	t.Log(db.InsertSql(data))
}

func TestUpdate(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(db.UpdateSql(data))
}

func TestLoad(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(db.LoadSql(data, db.WithOutWhere()))
	t.Log(db.LoadSql(data, db.WithWhere(&model.SimplePlayerData{AccountId: 1, PlayerName: "11"}), db.WithLimit(10)))
}

func TestDelete(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(db.DeleteSql(data))
}

func TestSave(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(db.SaveSql(data))
}
