package model

type (
	ModelData struct {
		Dirty bool `sql:"-"` //脏标记
	}

	SimplePlayerData struct {
		ModelData
		table          string `sql:"table;name:tbl_player"`
		PlayerId       int64  `sql:"primary;name:player_id"`
		AccountId      int64  `sql:"name:account_id"`
		PlayerName     string `sql:"name:player_name"`
		Level          int    `sql:"name:level"`
		Sex            int    `sql:"name:sex"`
		Gold           int    `sql:"froce;name:gold"`
		DrawGold       int    `sql:"name:draw_gold"`
		Vip            int    `sql:"name:vip"`
		LastLogoutTime int64  `sql:"datetime;name:last_logout_time"`
		LastLoginTime  int64  `sql:"datetime;name:last_login_time"`
	}

	PlayerKvItem struct {
		Key   int
		Value int64
	}

	PlayerKvData struct {
		ModelData
		table    string                  `sql:"table;name:tbl_player_kv"`
		PlayerId int64                   `sql:"primary;name:player_id"`
		DataMap  map[int64]*PlayerKvItem `sql:"blob;name:data_map"`
	}

	Bag struct {
		DataMap map[int64]*Item `sql:"blob;name:data_map"`
	}

	ItemData struct {
		ModelData
		table    string `sql:"table;name:tbl_item"`
		PlayerId int64  `sql:"primary;name:player_id"`
		Bag      Bag    `sql:""`
	}

	EquipData struct {
		ModelData
		table    string           `sql:"table;name:tbl_equip"`
		PlayerId int64            `sql:"primary;name:player_id"`
		DataMap  map[int64]*Equip `sql:"blob;name:data_map"`
	}

	MailData struct {
		ModelData
		table    string              `sql:"table;name:tbl_mail"`
		PlayerId int64               `sql:"primary;name:player_id"`
		DataMap  map[int64]*MailItem `sql:"blob;name:data_map"`
	}

	SocialData struct {
		ModelData
		table    string                `sql:"table;name:tbl_social"`
		PlayerId int64                 `sql:"primary;name:player_id"`
		DataMap  map[int64]*SocialItem `sql:"blob;name:data_map"`
	}
)

type (
	//人物结构
	PlayerData struct {
		SimplePlayerData
		PlayerKvData
		ItemData
		EquipData
		MailData
		SocialData
	}
)

// 玩家初始化数据
func (this *PlayerData) Init(PlayerId int64) {
	this.PlayerKvData = PlayerKvData{PlayerId: PlayerId, DataMap: map[int64]*PlayerKvItem{}}
	this.ItemData = ItemData{PlayerId: PlayerId, Bag: Bag{DataMap: map[int64]*Item{}}}
	this.EquipData = EquipData{PlayerId: PlayerId, DataMap: map[int64]*Equip{}}
	this.MailData = MailData{PlayerId: PlayerId, DataMap: map[int64]*MailItem{}}
	this.SocialData = SocialData{PlayerId: PlayerId, DataMap: map[int64]*SocialItem{}}
}
