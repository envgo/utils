package utils

// // 设置redis
// utils.RedisPool = *utils.NewRedisStore(
// 	&utils.RedisStoreOptions{
// 		Network:        "tcp",
// 		Address:        conf.Config.Redis.Address,
// 		ConnectTimeout: time.Duration(50) * time.Millisecond,
// 		ReadTimeout:    time.Duration(50) * time.Millisecond,
// 		WriteTimeout:   time.Duration(50) * time.Millisecond,
// 		IdleTimeout:    time.Duration(50) * time.Millisecond,
// 		Database:       conf.Config.Redis.Database,
// 		MaxIdle:        conf.Config.Redis.MaxIdle,
// 		MaxActive:      conf.Config.Redis.MaxActive,
// 	})
// // 使用redis
// utils.RedisPool.Get("keepredis")

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

// REDIS集群配置
var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout uint32 = 1000
var DefaultRedisWriteTimeout uint32 = 1000

var rs *RedisStore
var RedisPool RedisStore

type RedisStoreOptions struct {
	Network              string
	Address              string
	ConnectTimeout       time.Duration
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	Database             int           // Redis database to use for session keys
	KeyPrefix            string        // If set, keys will be KeyPrefix:SessionID (semicolon added)
	BrowserSessServerTTL time.Duration // Defaults to 2 days
	MaxIdle              int
	MaxActive            int
}

type RedisStore struct {
	pool    *redis.Pool
	options *RedisStoreOptions
}

func NewRedisStore(opts *RedisStoreOptions) *RedisStore {
	fmt.Println("NewRedisStore Init")
	pool := &redis.Pool{
		MaxIdle:     opts.MaxIdle,
		MaxActive:   opts.MaxActive,
		IdleTimeout: opts.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", opts.Address)
		},
	}
	rs := &RedisStore{
		pool:    pool,
		options: opts,
	}
	// 初始化
	return rs
}

func (self *RedisStore) Pool() *redis.Pool {
	return self.pool
}

//设置bit 值
func (self *RedisStore) Setbit(key string, offset int, val int) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("SETBIT", key, offset, val))
}

//获取bit 值
func (self *RedisStore) Getbit(key string, offset int) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GETBIT", key, offset))
}

//获取bit 值
func (self *RedisStore) Incrby(key string, offset int) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("INCRBY", key, offset))
}

func (self *RedisStore) SetList(key string, data []struct{}) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("SET", key, data))
}

func (self *RedisStore) Set(key string, data string) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("SET", key, data))
}

func (self *RedisStore) Incr(key string) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("INCR", key))
}
func (self *RedisStore) Setex(key string, timeout int, data string) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("SETEX", key, timeout, data))
}

func (self *RedisStore) IsKeyExist(key string) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func (self *RedisStore) Get(key string) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (self *RedisStore) Lpush(key string, data string) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("LPUSH", key, data))
}

func (self *RedisStore) Rpop(key string) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("RPOP", key))
}

func (self *RedisStore) Llen(key string) (uint64, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Uint64(conn.Do("LLEN", key))
}

func (self *RedisStore) Hincrby(key string, field string, number int64) (int64, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("HINCRBY", key, field, number))
}

func (self *RedisStore) Sadd(key string, data string) (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("SADD", key, data))
}

func (self *RedisStore) Smembers(key string) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("SMEMBERS", key))
}

func (self *RedisStore) Keys(key string) ([]string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", "*id_*"))
}

//设置哈希字段的字符串值
func (self *RedisStore) Hset(key string, field string, data string) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HSET", key, field, data))
}

//为多个哈希字段分别设置它们的值
func (self *RedisStore) Hmset(args ...interface{}) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("HMSET", args))
}

func (self *RedisStore) Hget(key string, field string) (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("HGET", key, field))
}

//获取存储在指定键的哈希中的所有字段和值
func (self *RedisStore) Hgetall(key string) (interface{}, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

//获取存储在指定键的哈希中的所有字段和值
func (self *RedisStore) HgetallInt(key string) (map[string]int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.IntMap(conn.Do("HGETALL", key))
}

//设置哈希字段的字符串值
func (self *RedisStore) Hexists(key string, field string) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HEXISTS", key, field))
}

//开启事务
func (self *RedisStore) Multi() (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("MULTI"))
}

//提交事务
func (self *RedisStore) Exec() (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	a, err := redis.String(conn.Do("EXEC"))
	log.Print("Exec", a, err)
	return redis.String(conn.Do("EXEC"))
}

//回滚MULTI之后发出的所有命令
func (self *RedisStore) Discard() (string, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("DISCARD"))
}

func (self *RedisStore) Setnx(key string, data string) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("SETNX", key, data))
}

func (self *RedisStore) Expire(key string, timeout int) (int, error) {
	conn := self.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("EXPIRE", key, timeout))
}
