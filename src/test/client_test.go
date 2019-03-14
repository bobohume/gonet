package main_test

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type(
	TopRank struct{
		Id int64	`sql:"primary;name:id"			json:"id"		json:"id"`
		Type int8	`sql:"primary;name:type"		json:"type"		json:"type"`
		Name string `sql:"name:name"				json:"name"		json:"name"`
		Score int `sql:"name:score"					json:"score"	json:"score"`
		Value [2]int `sql:"name:value"				json:"value"	json:"value"`
		LastTime int64 `sql:"datetime;name:last_time"	json:"last_time"	json:"last_time"`
	}
)

var(
	ntimes = 1000000
)

func TestJson(t *testing.T){
	for i := 0; i < ntimes; i++{
		json.Marshal(&TopRank{})
	}
}

func TestBson(t *testing.T){
	for i := 0; i < ntimes; i++{
		bson.Marshal(&TopRank{})
	}
	proto.NewBuffer([]byte{}).EncodeStringBytes()
}

func TestUJson(t *testing.T){
	buff, _ := json.Marshal(&TopRank{})
	for i := 0; i < ntimes; i++{
		json.Unmarshal(buff, &TopRank{})
	}
}

func TestUbson(t *testing.T){
	buff, _ := bson.Marshal(&TopRank{})
	for i := 0; i < ntimes; i++{
		bson.Unmarshal(buff, &TopRank{})
	}
}