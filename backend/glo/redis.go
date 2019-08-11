package glo

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisCache Cache
)

type Cache struct {
	Pool              *redis.Pool
	DefaultExpiration time.Duration
}

// Serialization 序列化
func Serialization(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	case reflect.Map:
	}
	k, err := json.Marshal(value)
	return k, err
}

// Deserialization 反序列化
func Deserialization(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt
		return
	}
	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				fmt.Printf("Deserialization: failed to parse int '%s': %s", string(byt), err)
			} else {
				p.SetInt(i)
			}
			return

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				fmt.Printf("Deserialization: failed to parse uint '%s': %s", string(byt), err)
			} else {
				p.SetUint(i)
			}
			return
		}
	}
	err = json.Unmarshal(byt, &ptr)
	return
}

// NewRedisPool func redis pool
func NewRedisPool() {
	maxIdle := Config.GopsAPI.Redis.MaxIdle
	expriedSec := time.Duration(Config.GopsAPI.Redis.DefaultExpried) * time.Second
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(Config.GopsAPI.Redis.IdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Config.GopsAPI.Redis.Addr)
			if err != nil {
				log.Panicln("打开redis数据库失败:", err)
				return c, fmt.Errorf("redis connection error: %s", err)
			}
			if Config.GopsAPI.Redis.Password != `` {
				if _, authErr := c.Do("AUTH", Config.GopsAPI.Redis.Password); authErr != nil {
					log.Panicln("验证redis密码数据库失败:", authErr)
					return c, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
	RedisCache = Cache{Pool: pool, DefaultExpiration: expriedSec}
}

// 判断所在的 key 是否存在
func (cache *Cache) Exist(name string) (bool, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", name))
	return v, err
}

// StringSet string 类型 添加, v 可以是任意类型
func (cache *Cache) StringSet(name string, v interface{}) error {
	conn := cache.Pool.Get()
	//s, _ := Serialization(v) // 序列化
	defer conn.Close()
	_, err := conn.Do("SET", name, v)
	return err
}

// StringGet 获取字符串类型的值
func (cache *Cache) StringGet(name string) (v interface{}, err error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	temp, _ := redis.Bytes(conn.Do("Get", name))
	err = Deserialization(temp, &v) // 反序列化
	return v, err
}

// //////////////////  hash ///////////
// 删除指定的 hash 键
func (cache *Cache) Hdel(name, key string) (bool, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var err error
	v, err := redis.Bool(conn.Do("HDEL", name, key))
	return v, err
}

// 设置单个值, value 还可以是一个 map slice 等
func (cache *Cache) HSet(name string, key string, value interface{}) (err error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	v, _ := Serialization(value)
	_, err = conn.Do("HSET", name, key, v)
	return
}

// 获取单个hash 中的值
func (cache *Cache) HGet(name, field string) (v interface{}, ok bool, err error) {
	ok = false
	conn := cache.Pool.Get()
	defer conn.Close()
	temp, _ := redis.Bytes(conn.Do("HGET", name, field))
	err = Deserialization(temp, &v) // 反序列化
	if err == nil {
		ok = true
	}
	return
}

// HMSet 设置单个值, value 还可以是一个 map slice 等
func (cache *Cache) HMSet(name string, kv map[string]interface{}) (err error) {
	// res := make(map[string]interface{}, len(kv))
	// for k, v := range kv {
	// 	d, _ := Serialization(v)
	// 	res[k] = d
	// }
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("HMSET", redis.Args{}.Add(name).AddFlat(kv)...)
	return
}

// HIncrby hash key 自增
func (cache *Cache) HIncrby(name, field string) (ok bool, err error) {
	ok = false
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("HINCRBY", name, field, 1)
	if err == nil {
		ok = true
	}
	return
}

// HIncrby hash key 自增
func (cache *Cache) HGetAll(name string) (data map[string]string, err error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	tmp, err := redis.StringMap(conn.Do("HGETALL", name))
	if err != nil {
		data = map[string]string{}
	} else {
		data = tmp
	}
	return
}

// HGetString 获取单个hash中的值string
func (cache *Cache) HGetString(name, field string) (v interface{}, ok bool, err error) {
	ok = false
	conn := cache.Pool.Get()
	defer conn.Close()
	tmp, err := redis.Bytes(conn.Do("HGET", name, field))
	if err == nil {
		ok = true
		v = string(tmp)
	}
	return
}

// 查看hash 中指定是否存在
func (cache *Cache) HExists(name, field string) (bool, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var err error
	v, err := redis.Bool(conn.Do("HEXISTS", name, field))
	return v, err
}

// Del 删除Key
func (cache *Cache) Del(name string) (err error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", name)
	return
}

// RdsDisConnect DisConnect Redis Pool
func RdsDisConnect() {
	RedisCache.Pool.Close()
}
