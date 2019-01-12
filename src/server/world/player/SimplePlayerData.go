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

/*func (this *SimplePlayerData) ReadData(b *base.BitStream){
	base.ReadData(this, b)
}*/

/*func (this *SimplePlayerData) WriteData(b *base.BitStream){
	base.WriteData(this, b)
}*/