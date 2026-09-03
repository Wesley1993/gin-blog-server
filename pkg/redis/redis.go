package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端封装
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// DefaultClient 全局默认 Redis 客户端实例
var DefaultClient *RedisClient

// InitRedis 初始化全局 Redis 客户端
func InitRedis(addr, password string, db, poolSize, minIdleConns int) {
	DefaultClient = NewRedisClient(RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})
}

// GetClient 获取全局默认 Redis 客户端
func GetClient() *RedisClient {
	return DefaultClient
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(config RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

// Set 设置键值对
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// SetEX 设置键值对并设置过期时间
func (r *RedisClient) SetEX(key string, value interface{}, seconds int64) error {
	return r.client.Set(r.ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// Del 删除键
func (r *RedisClient) Del(keys ...string) (int64, error) {
	return r.client.Del(r.ctx, keys...).Result()
}

// Incr 自增
func (r *RedisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Decr 自减
func (r *RedisClient) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// HSet Hash设置
func (r *RedisClient) HSet(key string, values ...interface{}) (int64, error) {
	return r.client.HSet(r.ctx, key, values...).Result()
}

// HGet Hash获取
func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// LPush 列表左推
func (r *RedisClient) LPush(key string, values ...interface{}) (int64, error) {
	return r.client.LPush(r.ctx, key, values...).Result()
}

// RPop 列表右弹
func (r *RedisClient) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

// LRange 列表范围获取
func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// Ping 检测连接可用性，超时控制由外部 ctx 决定
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(r.ctx, key, expiration).Result()
}

// SetWithTTL 设置键值对并指定过期时间（语义别名）
func (r *RedisClient) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

// SAdd 向集合添加成员
func (r *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(r.ctx, key, members...).Result()
}

// SMembers 获取集合所有成员
func (r *RedisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// ScanDel 使用 SCAN 迭代匹配前缀模式的键并批量删除（避免 KEYS 阻塞）
func (r *RedisClient) ScanDel(pattern string) (int64, error) {
	var cursor uint64
	var deleted int64
	for {
		keys, next, err := r.client.Scan(r.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return deleted, err
		}
		if len(keys) > 0 {
			n, err := r.client.Del(r.ctx, keys...).Result()
			if err != nil {
				return deleted, err
			}
			deleted += n
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return deleted, nil
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client 获取底层Redis客户端
func (r *RedisClient) Client() *redis.Client {
	return r.client
}
