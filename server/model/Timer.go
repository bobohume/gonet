package model

type(
	Timer struct {
		table 	 string	`sql:"table;name:tbl_timerset"`
		Id int `sql:"primary;name:id"`						//定时器Id
		PlayerId int64 `sql:"primary;name:player_id"`		//玩家Id
		Flag int64	`sql:"name:flag"`						//定时器数据
		ExpireTime int64 `sql:datetime;name:expire_time`	//定时器过期时间
	}
)
