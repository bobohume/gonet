package rd

import (
	"base"
	"db"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"reflect"
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

///Get 获取一个值
func Get(key string, database int) ([]byte, error){
	c := POOL.Get()
	defer c.Close()

	c.Do("SELECT", database)
	data, err := redis.Bytes(c.Do("GET", key))
	return data, err
}

//Set 设置一个值
func Set(key string, val interface{}, timeout int, database int) (err error) {
	c := POOL.Get()
	defer c.Close()

	c.Do("SELECT", database)
	_, err = c.Do("SETEX", key, timeout, val)
	if err != nil {
		return err
	}

	return nil
}

//IsExist 判断key是否存在
func Exist(key string, database int) bool {
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
func Delete(key string, database int) error {
	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	if _, err := c.Do("DEL", key); err != nil {
		return err
	}

	return nil
}

//Expire 超时
func Expire(key string, time, database int) error {
	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	if _, err := c.Do("EXPIRE", key, time); err != nil {
		return err
	}

	return nil
}

func isJson(obj interface{}) bool{
	oTpye := reflect.TypeOf(obj)
	if oTpye.Kind() == reflect.Ptr {
		oTpye = reflect.TypeOf(obj).Elem()
	}

	if oTpye.Kind() == reflect.Struct{
		return true
	}
	return false
}

//redis set
func Exec(time int, database int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode()
			//base.GLOG.Println("redis exec", err)
		}
	}()
	c := POOL.Get()
	defer c.Close()
	bJson := false
	if len(args) > 0{
		bJson = isJson(args[0])
	}

	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	if bJson{
		for _, v := range args{
			data, err := json.Marshal(v)
			if err == nil{
				vargs = append(vargs, data)
			}
		}
	} else{
		vargs = append(vargs, args...)
	}
	c.Send(cmd, vargs...)
	if time != -1{
		c.Send("EXPIRE", key, time)
	}
	c.Flush()
	return nil
}

func ExecKV(time int, database int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode()
			//base.GLOG.Println("redis execkv", err)
		}
	}()

	c := POOL.Get()
	defer c.Close()
	bJson := false
	if len(args) > 1{
		bJson = isJson(args[1])
	}

	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	if bJson{
		for i, v := range args{
			if i % 2 != 0{
				data, err := json.Marshal(v)
				if err == nil{
					vargs = append(vargs, data)
				}else{
					vargs = append(vargs, "")
				}
			}else{
				vargs = append(vargs, v)
			}
		}
	} else{
		vargs = append(vargs, args...)
	}
	c.Send(cmd, vargs...)
	if time != -1{
		c.Send("EXPIRE", key, time)
	}
	c.Flush()
	return nil
}

//redis get
func Query(database int, cmd, key string, obj interface{}, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode()
			//base.GLOG.Println("redis query", err)
		}
	}()

	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	vargs = append(vargs, args...)
	c.Send(cmd, vargs...)
	c.Flush()
	reply, err := c.Receive()
	if err != nil{
		return err
	}

	bJson := isJson(obj)

	switch reply.(type) {
	case int64:
		*(obj.(*int64)), _ = redis.Int64(reply, err)
	case float64:
		*(obj.(*float64)), _ = redis.Float64(reply, err)
	case bool:
		*(obj.(*bool)), _ = redis.Bool(reply, err)
	case string:
		*(obj.(*string)), _ = redis.String(reply, err)
	case int:
		*(obj.(*int)), _ = redis.Int(reply, err)
	case []interface{}:
		var aData [][]byte
		aData, err = redis.ByteSlices(reply, err)
		if err == nil {
			r := reflect.Indirect(reflect.ValueOf(obj))
			isPtr := false
			if kind := r.Kind(); kind == reflect.Slice {
				rType := r.Type().Elem()
				if rType.Kind() == reflect.Ptr {
					isPtr = true
					rType = rType.Elem()
				}

				for _, v := range aData{
					elem := reflect.New(rType).Elem()
					if bJson{
						if json.Unmarshal(v, elem.Addr().Interface()) == nil{
							if isPtr{
								r.Set(reflect.Append(r, elem.Addr()))
							}else{
								r.Set(reflect.Append(r, elem))
							}
						}
					}else{
						elem.SetString(string(v))
						if isPtr{
							r.Set(reflect.Append(r, elem.Addr()))
						}else{
							r.Set(reflect.Append(r, elem))
						}
					}
				}
			}
		}
	case interface{}:
		if bJson{
			var data []byte
			data, err = redis.Bytes(reply, err)
			if err == nil {
				return json.Unmarshal(data, obj)
			}
		}else{
			var val string
			val, err = redis.String(reply, err)
			if err == nil{
				*obj.(*string) = val
			}
		}
	}
	return err
}

//redis get map
func QueryKV(database int, cmd, key string, obj interface{}, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode()
			//base.GLOG.Println("redis querykv", err)
		}
	}()

	c := POOL.Get()
	defer c.Close()
	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	vargs = append(vargs, args...)
	c.Send(cmd, vargs...)
	c.Flush()
	reply, err := c.Receive()
	if err != nil{
		return err
	}

	bJson := false
	var aData [][]byte
	aData, err = redis.ByteSlices(reply, err)
	if err == nil {
		r := reflect.Indirect(reflect.ValueOf(obj))
		isPtr := false
		if kind := r.Kind(); kind == reflect.Map{
			rType := r.Type().Elem()
			rkType := r.Type().Key()
			if rType.Kind() == reflect.Ptr {
				isPtr = true
				rType = rType.Elem()
			}

			if rType.Kind() == reflect.Struct{
				bJson = true
			}

			tKey := reflect.New(rkType)
			var vKey reflect.Value
			for i, v := range aData{
				if i % 2==0{
					row := db.NewRow()
					row.Set("key", string(v))
					switch tKey.Elem().Kind() {
					case reflect.Int:
						vKey = reflect.ValueOf(row.Int("key"))
					case reflect.Int32:
						vKey = reflect.ValueOf(int32(row.Int("key")))
					case reflect.Int64:
						vKey = reflect.ValueOf(row.Int64("key"))
					case reflect.String:
						vKey = reflect.ValueOf(row.String("key"))
					case reflect.Float32:
						vKey = reflect.ValueOf(row.Float32("key"))
					case reflect.Float64:
						vKey = reflect.ValueOf(row.Float64("key"))
					case reflect.Bool:
						vKey = reflect.ValueOf(row.Bool("key"))
					}
				}else{
					elem := reflect.New(rType).Elem()
					if bJson{
						if json.Unmarshal(v, elem.Addr().Interface()) == nil{
							if isPtr{
								r.SetMapIndex(vKey, elem.Addr())
							}else{
								r.SetMapIndex(vKey, elem)
							}
						}
					}else{
						elem.SetString(string(v))
						if isPtr{
							r.SetMapIndex(vKey, elem.Addr())
						}else{
							r.SetMapIndex(vKey, elem)
						}
					}
				}

			}
		}
	}
	return err
}