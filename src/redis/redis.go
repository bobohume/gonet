package redis

import (
	"base"
	"github.com/gomodule/redigo/redis"
	"github.com/gomodule/redigo/redisx"
	"strconv"
)

type(
	Row struct {
		m_Resut map[string] string
	}

	IRow interface {
		init()
		Get(key string) string
		String(key string) string
		Int(key string) int
		Int64(key string) int64
		Flot32(key string) float32
		Flot64(key string) float64
		Bool(key string) bool
		Time(key string) int64
	}
)

func (this *Row) init() {
	this.m_Resut = make(map[string] string)
}

func (this *Row) Get(key string) string{
	v, exist := this.m_Resut[key]
	if exist{
		return v
	}

	return ""
}

func (this *Row) String(key string) string{
	return this.Get(key)
}

func (this *Row) Int(key string) int{
	n, _ := strconv.Atoi(this.Get(key))
	return n
}

func (this *Row) Int64(key string) int64{
	n, _ := strconv.ParseInt(this.Get(key), 0, 64)
	return n
}

func (this *Row) Flot32(key string) float32{
	n, _ := strconv.ParseFloat(this.Get(key), 32)
	return float32(n)
}

func (this *Row) Flot64(key string) float64{
	n, _ := strconv.ParseFloat(this.Get(key), 64)
	return n
}

func (this *Row) Bool(key string) bool{
	n, _ := strconv.ParseBool(this.Get(key))
	return n
}

func (this *Row) Time(key string) int64{
	return base.GetDBTime(this.Get(key)).Unix()
}

func Exec(m *redisx.ConnMux, cmd, key string, args []interface{}){
	c1 := m.Get()
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	vargs = append(vargs, args...)
	c1.Send(cmd, key, vargs)
	c1.Flush()
	c1.Close()
}

func Query(m *redisx.ConnMux, cmd string, args string) *Row{
	row := &Row{}
	row.init()
	c1 := m.Get()
	c1.Send(cmd, args)
	c1.Flush()
	row.m_Resut, _ = redis.StringMap(c1.Receive())
	c1.Close()
	return row
}

//	a := AccountDB{}
//	a.AccountId = 2
//	a.AccountName = "test222"
//	a.LoginIp = "192.168.0.10"
//	fmt.Println(redis2.QueryRow(c, "HMSET", "player", redis2.RedisStr(&a))...)
//	cc := redis2.QueryRow(c, "HGETALL", "player")