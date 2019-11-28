package rd

import (
	"github.com/gomodule/redigo/redis"
	"runtime"
	"time"
)

var(
	Redis_DB_TopRank = 0
)

var POOL *redis.Pool

//@title 启动redis, redispo.Pool（连接池）
func OpenRedisPool(ip, pwd string) error {
	cpuNum := runtime.NumCPU()
	POOL = &redis.Pool{
		MaxIdle:     cpuNum,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ip)
			if err != nil {
				return nil, err
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

///Do a func can do no defer close
func Do(database int, pFunc func(c redis.Conn)){
	c := POOL.Get()
	defer c.Close()

	c.Do("SELECT", database)
	pFunc(c)
}

///Get 获取一个值
func Get(database int, key string) ([]byte, error){
	c := POOL.Get()
	defer c.Close()

	c.Do("SELECT", database)
	data, err := redis.Bytes(c.Do("GET", key))
	return data, err
}

//Set 设置一个值
func Set(database int, timeout int, key string, val interface{}) (err error) {
	c := POOL.Get()
	defer c.Close()

	c.Do("SELECT", database)
	_, err = c.Do("SET", key, val, "EX", timeout, )
	if err != nil {
		return err
	}

	return nil
}

//IsExist 判断key是否存在
func Exist(database int, key string) bool {
	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	a, _ := c.Do("EXISTS", key)
	i := a.(int64)
	if i > 0 {
		return true
	}
	return false
}

//Delete 删除
func Delete(database int, key string) error {
	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	if _, err := c.Do("DEL", key); err != nil {
		return err
	}

	return nil
}

//Expire 超时
func Expire(database, timeout int, key string) error {
	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	if _, err := c.Do("EXPIRE", key, timeout); err != nil {
		return err
	}

	return nil
}
