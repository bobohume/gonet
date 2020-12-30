package base

import (
	"math/big"
	"math/rand"
)

type (
	Dh struct{
		q big.Int//素数q
		a big.Int//q的原根a
		x big.Int//私钥
		Y1 big.Int//自己公钥
		Y2 big.Int//对方公钥
	}

	IDh interface {
		Init()
		generatePrik()//生成私钥
		generatePubk()//生成公钥
		ExchangePubk(key int64)//交换公钥
		PubKey() int64//公钥
		ShareKey()int64//生成共享密钥
	}
)

func (this *Dh) Init(){
	this.q = *big.NewInt(97)
	this.a = *big.NewInt(5)
	this.generatePrik()
	this.generatePubk()
}

func (this *Dh) generatePrik(){
	r := big.NewInt(int64(rand.Int()))
	this.x.Mod(r, &this.q)
}

func (this *Dh) generatePubk(){
	this.Y1.Exp(&this.a, &this.x, &this.q)
}

func (this *Dh) ExchangePubk(key int64){
	this.Y2 = *big.NewInt(key)
}

func (this *Dh) PubKey() int64{
	return this.Y1.Int64()
}

func (this *Dh) ShareKey() int64{
	return big.NewInt(0).Exp(&this.Y2, &this.x, &this.q).Int64()
}