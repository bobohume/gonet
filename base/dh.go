package base

import (
	"math/big"
	"math/rand"
)

type (
	Dh struct {
		q  big.Int //素数q
		a  big.Int //q的原根a
		x  big.Int //私钥
		Y1 big.Int //自己公钥
		Y2 big.Int //对方公钥
	}

	IDh interface {
		Init()
		generatePrik()          //生成私钥
		generatePubk()          //生成公钥
		ExchangePubk(key int64) //交换公钥
		PubKey() int64          //公钥
		ShareKey() int64        //生成共享密钥
	}
)

func (d *Dh) Init() {
	d.q = *big.NewInt(97)
	d.a = *big.NewInt(5)
	d.generatePrik()
	d.generatePubk()
}

func (d *Dh) generatePrik() {
	r := big.NewInt(int64(rand.Int()))
	d.x.Mod(r, &d.q)
	d.x.Add(&d.x, big.NewInt(1))
}

func (d *Dh) generatePubk() {
	d.Y1.Exp(&d.a, &d.x, &d.q)
}

func (d *Dh) ExchangePubk(key int64) {
	d.Y2 = *big.NewInt(key)
}

func (d *Dh) PubKey() int64 {
	return d.Y1.Int64()
}

func (d *Dh) ShareKey() int64 {
	return big.NewInt(0).Exp(&d.Y2, &d.x, &d.q).Int64()
}
