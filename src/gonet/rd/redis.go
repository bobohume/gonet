package rd

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"gonet/base"
	"gonet/db"
	"reflect"
	"runtime"
	"time"
)

var(
	Redis_DB_TopRank = 0
)

var POOL *redis.Pool

type BYTES_TYPE int
const(
	BYTES_NONE BYTES_TYPE = iota
	BYTES_PB BYTES_TYPE = iota
	BYTES_JSON BYTES_TYPE = iota
)

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
	_, err = c.Do("SETEX", key, timeout, val)
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

func BytesType(obj interface{}) BYTES_TYPE{
	oTpye := reflect.TypeOf(obj)
	for oTpye.Kind() == reflect.Ptr {
		oTpye = oTpye.Elem()
	}

	if oTpye.Kind() == reflect.Struct{
		if oTpye.NumField() != 0{
			sf := oTpye.Field(0)
			if len(sf.Tag.Get("protobuf")) > 0{
				return BYTES_PB
			}else if len(sf.Tag.Get("josn")) > 0{
				return BYTES_JSON
			}
		}
		return BYTES_NONE
	}
	return BYTES_NONE
}

//redis set
func Exec(database int, timeout int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
			//base.GLOG.Println("redis exec", err)
		}
	}()
	c := POOL.Get()
	defer c.Close()
	nType := BYTES_NONE
	if len(args) > 0{
		nType = BytesType(args[0])
	}

	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	if nType == BYTES_JSON{
		for _, v := range args{
			data, err := json.Marshal(v)
			if err == nil{
				vargs = append(vargs, data)
			}
		}
	}else if nType == BYTES_PB{
		for _, v := range args{
			val := reflect.ValueOf(v)
			for val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			data, err := proto.Marshal(val.Addr().Interface().(proto.Message))
			if err == nil{
				vargs = append(vargs, data)
			}
		}
	} else{
		vargs = append(vargs, args...)
	}
	c.Send(cmd, vargs...)
	if timeout != -1{
		c.Send("EXPIRE", key, timeout)
	}
	c.Flush()
	return nil
}

func ExecKV(database int, timeout int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
			//base.GLOG.Println("redis execkv", err)
		}
	}()

	c := POOL.Get()
	defer c.Close()
	nType := BYTES_NONE
	if len(args) > 1{
		nType = BytesType(args[1])
	}

	c.Do("SELECT", database)
	vargs := make([]interface{}, 0)
	vargs = append(vargs, key)
	if nType == BYTES_JSON{
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
	} else if nType == BYTES_PB{
		for i, v := range args{
			if i % 2 != 0{
				val := reflect.ValueOf(v)
				for val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				data, err := proto.Marshal(val.Addr().Interface().(proto.Message))
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
	if timeout != -1{
		c.Send("EXPIRE", key, timeout)
	}
	c.Flush()
	return nil
}

//redis get
func Query(obj interface{}, database int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
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

	nType := BytesType(obj)

	switch reply.(type) {
	case int64:
		ptr, ok := (obj.(*int64))
		if ok{
			*ptr, _ = redis.Int64(reply, err)
		}
	case float64:
		ptr, ok := (obj.(*float64))
		if ok{
			*ptr, _ = redis.Float64(reply, err)
		}
	case bool:
		ptr, ok := (obj.(*bool))
		if ok{
			*ptr, _ = redis.Bool(reply, err)
		}
	case string:
		ptr, ok := (obj.(*string))
		if ok{
			*ptr, _ = redis.String(reply, err)
		}
	case int:
		ptr, ok := (obj.(*int))
		if ok{
			*ptr, _ = redis.Int(reply, err)
		}
	case []interface{}:
		var aData [][]byte
		aData, err = redis.ByteSlices(reply, err)
		if err == nil {
			r := reflect.Indirect(reflect.ValueOf(obj))
			isPtr := false
			if kind := r.Kind(); kind == reflect.Slice {
				rType := r.Type().Elem()
				for rType.Kind() == reflect.Ptr {
					isPtr = true
					rType = rType.Elem()
				}

				for _, v := range aData{
					elem := reflect.New(rType).Elem()
					if nType == BYTES_JSON{
						if json.Unmarshal(v, elem.Addr().Interface()) == nil{
							if isPtr{
								r.Set(reflect.Append(r, elem.Addr()))
							}else{
								r.Set(reflect.Append(r, elem))
							}
						}
					} else if nType == BYTES_PB{
						if proto.Unmarshal(v, elem.Addr().Interface().(proto.Message)) == nil{
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
		if nType == BYTES_JSON{
			var data []byte
			data, err = redis.Bytes(reply, err)
			if err == nil {
				return json.Unmarshal(data, obj)
			}
		} else if nType == BYTES_PB{
			var data []byte
			data, err = redis.Bytes(reply, err)
			if err == nil {
				val := reflect.ValueOf(obj)
				for val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				return proto.Unmarshal(data, val.Addr().Interface().(proto.Message))
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
func QueryKV(obj interface{}, database int, cmd, key string, args ...interface{}) error{
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
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

	nType := BYTES_NONE
	var aData [][]byte
	aData, err = redis.ByteSlices(reply, err)
	if err == nil {
		r := reflect.ValueOf(obj)
		for r.Kind() == reflect.Ptr {
			r = r.Elem()
		}
		isPtr := false
		if kind := r.Kind(); kind == reflect.Map{
			rType := r.Type().Elem()
			rkType := r.Type().Key()
			for rType.Kind() == reflect.Ptr {
				isPtr = true
				rType = rType.Elem()
			}

			if rType.Kind() == reflect.Struct{
				if rType.NumField() != 0{
					sf := rType.Field(0)
					if len(sf.Tag.Get("protobuf")) > 0{
						nType = BYTES_PB
					}else if len(sf.Tag.Get("json")) > 0{
						nType = BYTES_JSON
					}
				}
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
					if nType == BYTES_JSON{
						if json.Unmarshal(v, elem.Addr().Interface()) == nil{
							if isPtr{
								r.SetMapIndex(vKey, elem.Addr())
							}else{
								r.SetMapIndex(vKey, elem)
							}
						}
					}else if nType == BYTES_PB{
						if proto.Unmarshal(v, elem.Addr().Interface().(proto.Message)) == nil{
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