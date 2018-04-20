package player

import "base"

type (
	SimplePlayerData struct{
		AccountId int
		PlayerId int `primary`
		PlayerName string
		Level int
		Sex	  int
		Gold  int
		DrawGold int
		Vip int
		LastLogoutTime int64 `datetime`
		LastLoginTime int64	`datetime`
	}
)

func (this *SimplePlayerData) ReadData(b *base.BitStream){
	base.ReadData(this, b)
}

func (this *SimplePlayerData) WriteData(b *base.BitStream){
	base.WriteData(this, b)
}