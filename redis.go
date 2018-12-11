package utils

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"
)

// REDIS集群配置
var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout uint32 = 1000
var DefaultRedisWriteTimeout uint32 = 1000
var DefaultIdleTimeout uint32 = 10

//var rs *RedisStore
var RedisPool RedisStore
var rsconfigs redisConfig

type RedisStore struct {
	pool    *redis.Pool
	options *configItem
}

// redis配置
type configItem struct {
	Address     string `toml:"address"`
	ListName    string `toml:"listname"`
	Database    int    `toml:"Database,omitempty"`
	MaxIdle     int    `toml:"MaxIdle,omitempty"`     //最大的空闲连接数
	MaxActive   int    `toml:"MaxActive,omitempty"`   //最大的激活连接数
	IdleTimeout int    `toml:"IdleTimeout,omitempty"` //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
	Timeout     int    `toml:"timeout,omitempty"`
}

type redisConfig struct {
	Configs map[string]configItem `toml:"redis"`
}

func RedisLoad(cfgPath string, sub string) *RedisStore {
	if _, err := toml.DecodeFile(cfgPath, &rsconfigs); err != nil {
		panic("Redis Config load Failed")
	}
	subconfig, ishas := rsconfigs.Configs[sub]
	if !ishas {
		panic("Redis SubConfig load Failed")
	}
	if subconfig.Address == "" {
		panic("Redis config need address and listname")
	}
	if subconfig.MaxIdle == 0 {
		subconfig.MaxIdle = 20
	}
	if subconfig.MaxActive == 0 {
		subconfig.MaxActive = 100
	}
	if subconfig.IdleTimeout == 0 {
		subconfig.IdleTimeout = 10
	}
	if subconfig.Timeout == 0 {
		subconfig.Timeout = 1000
	}
	pool := &redis.Pool{
		MaxIdle:     subconfig.MaxIdle,
		MaxActive:   subconfig.MaxActive,
		IdleTimeout: time.Duration(subconfig.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", subconfig.Address)
		},
	}
	rs := &RedisStore{
		pool:    pool,
		options: &subconfig,
	}
	return rs
}

func NewRedisStore(opts *configItem) *RedisStore {
	fmt.Println("NewRedisStore Init")
	pool := &redis.Pool{
		MaxIdle:     opts.MaxIdle,
		MaxActive:   opts.MaxActive,
		IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", opts.Address)
		},
	}
	rs := &RedisStore{
		pool:    pool,
		options: opts,
	}
	return rs
}

func (self *RedisStore) Close() {
	fmt.Println("Redis Closed")
	self.Pool().Close()
}

func (self *RedisStore) Pool() *redis.Pool {
	return self.pool
}

//设置bit 值
func (self *RedisStore) Setbit(key string, offset int, val int) (int, error) {

	return redis.Int(self.Do("SETBIT", key, offset, val))
}

//获取bit 值
func (self *RedisStore) Getbit(key string, offset int) (int, error) {

	return redis.Int(self.Do("GETBIT", key, offset))
}

//获取bit 值
func (self *RedisStore) Incrby(key string, offset int) (int, error) {

	return redis.Int(self.Do("INCRBY", key, offset))
}

func (self *RedisStore) SetList(key string, data []struct{}) (bool, error) {

	return redis.Bool(self.Do("SET", key, data))
}

func (self *RedisStore) Set(key string, data string) (bool, error) {

	return redis.Bool(self.Do("SET", key, data))
}

func (self *RedisStore) Incr(key string) (bool, error) {

	return redis.Bool(self.Do("INCR", key))
}
func (self *RedisStore) Setex(key string, timeout int, data string) (string, error) {

	return redis.String(self.Do("SETEX", key, timeout, data))
}

func (self *RedisStore) IsKeyExist(key string) (bool, error) {

	return redis.Bool(self.Do("EXISTS", key))
}

func (self *RedisStore) Get(key string) (string, error) {

	return redis.String(self.Do("GET", key))
}

func (self *RedisStore) Lpush(key string, data string) (bool, error) {

	return redis.Bool(self.Do("LPUSH", key, data))
}

func (self *RedisStore) Rpop(key string) (string, error) {

	return redis.String(self.Do("RPOP", key))
}

func (self *RedisStore) Llen(key string) (uint64, error) {

	return redis.Uint64(self.Do("LLEN", key))
}

func (self *RedisStore) Hincrby(key string, field string, number int64) (int64, error) {

	return redis.Int64(self.Do("HINCRBY", key, field, number))
}

func (self *RedisStore) Sadd(key string, data string) (bool, error) {

	return redis.Bool(self.Do("SADD", key, data))
}

func (self *RedisStore) Smembers(key string) (string, error) {

	return redis.String(self.Do("SMEMBERS", key))
}

func (self *RedisStore) Keys(key string) ([]string, error) {

	return redis.Strings(self.Do("KEYS", "*id_*"))
}

//设置哈希字段的字符串值
func (self *RedisStore) Hset(key string, field string, data string) (int, error) {

	return redis.Int(self.Do("HSET", key, field, data))
}

//为多个哈希字段分别设置它们的值
func (self *RedisStore) Hmset(args ...interface{}) (string, error) {

	return redis.String(self.Do("HMSET", args))
}

func (self *RedisStore) Hget(key string, field string) (string, error) {

	return redis.String(self.Do("HGET", key, field))
}

//获取存储在指定键的哈希中的所有字段和值
func (self *RedisStore) Hgetall(key string) (interface{}, error) {

	return redis.StringMap(self.Do("HGETALL", key))
}

//获取存储在指定键的哈希中的所有字段和值
func (self *RedisStore) HgetallInt(key string) (map[string]int, error) {

	return redis.IntMap(self.Do("HGETALL", key))
}

//设置哈希字段的字符串值
func (self *RedisStore) Hexists(key string, field string) (int, error) {

	return redis.Int(self.Do("HEXISTS", key, field))
}

//开启事务
func (self *RedisStore) Multi() (string, error) {

	return redis.String(self.Do("MULTI"))
}

//提交事务
func (self *RedisStore) Exec() (string, error) {
	return redis.String(self.Do("EXEC"))
}

//回滚MULTI之后发出的所有命令
func (self *RedisStore) Discard() (string, error) {

	return redis.String(self.Do("DISCARD"))
}

func (self *RedisStore) Setnx(key string, data string) (int, error) {
	return redis.Int(self.Do("SETNX", key, data))
}

func (self *RedisStore) Expire(key string, timeout int) (int, error) {
	return redis.Int(self.Do("EXPIRE", key, timeout))
}

func (self *RedisStore) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := self.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
