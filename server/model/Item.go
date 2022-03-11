package model

type(
	Item struct {
		Id int64 		//物品唯一Id
		PlayerId int64	//玩家Id
		ItemId int		//模板Id
		Quantity int	//数量(对于装备，只能为1)
	}

	Equip struct {
		Id int64 			//物品唯一Id
		PlayerId int64 		//玩家Id
		ItemId int 			//模板Id
		Level int			//等级
		StrengthenLv int 	//强化等级
	}
)
