package cm

import (
	"gonet/base"
	"gonet/base/vector"
	"math"
	"strings"
)

type (
	//随即组
	RandUnit struct {
		Key    int32 //Key
		Val    int32 //Val
		LowVal int32 //下限
		UpVal  int32 //上限
	}

	//随机组
	RandGroup struct {
		Units  []*RandUnit
		MaxVal int32
	}
)

// 纯随机
func (r *RandGroup) Rand() *RandUnit {
	nNeed := len(r.Units)
	if nNeed > 0 {
		nIndex := base.RandI(0, nNeed-1)
		return r.Units[nIndex]
	}
	return nil
}

// 不重复随机属性
func (r *RandGroup) RandEx(Id []int32, need int) []*RandUnit {
	randomBuff := []*RandUnit{}
	buffVec := vector.Vector[*RandUnit]{}
	for _, v := range r.Units {
		//招到重复的
		bFind := false
		for _, v1 := range Id {
			if v.Key == v1 {
				bFind = true
				break
			}
		}
		if !bFind {
			buffVec.PushBack(v)
		}
	}
	nNeed := math.Min(float64(need), float64(buffVec.Len()))
	for ; nNeed > 0; nNeed-- {
		nIndex := base.RandI(0, buffVec.Len()-1)
		randomBuff = append(randomBuff, buffVec.Get(nIndex))
		buffVec.Remove(nIndex)
	}
	return randomBuff
}

// -------产生随机组-------//
func NewRandGroup(str string, bVal bool) *RandGroup {
	randGroup := &RandGroup{}
	stream := GetParamStream(str)
	nRandVal := int32(0)
	for stream.ReadFlag() {
		randUnit := &RandUnit{}
		randUnit.Key = int32(stream.ReadInt(32))
		if bVal {
			randUnit.Val = int32(stream.ReadInt(32))
		}
		randUnit.LowVal = nRandVal
		nRandVal += int32(stream.ReadInt(32))
		randUnit.UpVal = nRandVal
		randGroup.Units = append(randGroup.Units, randUnit)
	}
	randGroup.MaxVal = nRandVal
	return randGroup
}

// 解析数组  格式 value1;value2;
func GetArrayIntStream(str string) base.IBitStream {
	msg := make([]byte, 256)
	bitstrem := base.NewBitStream(msg, 256)
	rows := strings.Split(str, ";")
	for _, v := range rows {
		if v == "0" || v == "-1" || v == "" || v == " " {
			continue
		}
		bitstrem.WriteFlag(true)
		bitstrem.WriteInt(base.Int(v), 32)
	}
	bitstrem.SetPosition(0)
	return bitstrem
}

// 解析pair数组  格式 key1:value1;key2:value2;
func GetParamStream(str string) base.IBitStream {
	msg := make([]byte, 256)
	bitstrem := base.NewBitStream(msg, 256)
	rows := strings.Split(str, ";")
	for _, v := range rows {
		if v == "0" || v == "-1" || v == "" || v == " " {
			continue
		}
		bitstrem.WriteFlag(true)
		cols := strings.Split(v, ":")
		for _, col := range cols {
			bitstrem.WriteInt(base.Int(col), 32)
		}
	}
	bitstrem.SetPosition(0)
	return bitstrem
}
