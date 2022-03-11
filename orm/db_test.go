package orm_test

import (
	"fmt"
	"gonet/orm"
	"gonet/orm/model"
	"gonet/server/message"
	"testing"
)

func TestInsert(t *testing.T) {
	nameMap := make(map[string]bool)
	//nameMap["1"][2] = true
	nVal := nameMap["1"]
	fmt.Println(nVal)
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}, Vip: [8]int{1, 2, 3, 4, 5}}
	t.Log(orm.InsertSql(data))
}

func TestUpdate(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(orm.UpdateSql(data))
}

func TestLoad(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(orm.LoadSql(data, orm.WithOutWhere()))
	t.Log(orm.LoadSql(data, orm.WithWhere(&model.SimplePlayerData{PlayerId: 1, PlayerName: "11"}), orm.WithLimit(10)))
}

func TestDelete(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(orm.DeleteSql(data))
}

func TestSave(t *testing.T) {
	data := &model.SimplePlayerData{PLayerBlob: &message.PlayerData{}, PLayerBlobJson: &model.AA{A: 1, B: map[int]string{1: "test", 2: "test2"}}}
	t.Log(orm.SaveSql(data))
}
