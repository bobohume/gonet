package model

type(
	Item struct {
		table string	`sql:"table;name:tbl_item"`
		Id int64 `sql:"primary;name:id"`				//物品唯一Id
		PlayerId int64 `sql:"primary;name:player_id"`	//玩家Id
		ItemId int	`sql:"name:item_id"`				//模板Id
		Quantity int `sql:name:quantity`				//数量(对于装备，只能为1)
	}

	Equip struct {
		table string	`sql:"table;name:tbl_equip"`
		Id int64 `sql:"primary;name:id"`				//物品唯一Id
		PlayerId int64 `sql:"primary;name:player_id"`	//玩家Id
		ItemId int `sql:"name:item_id"`					//模板Id
		Level int	`sql:"name:level"`					//等级
		StrengthenLv int `sql:"name:strengthen_lv"`		//强化等级
	}
)
