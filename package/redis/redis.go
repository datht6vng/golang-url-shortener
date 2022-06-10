package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"trueid-shorten-link/config"

	"github.com/go-redis/redis/v7"
)

type Config struct {
	Addrs     []string `mapstructure:"addrs"`
	Pwd       string   `mapstructure:"password"`
	DB        int      `mapstructure:"db"`
	IsCluster bool     `mapstructure:"cluster"`
}

type Redis struct {
	cluster     *redis.ClusterClient
	single      *redis.Client
	clusterMode bool
}

func Connect() *Redis {
	return NewRedis(Config{
		Addrs:     config.Config.Redis.Address,
		Pwd:       config.Config.Redis.Password,
		DB:        config.Config.Redis.Database,
		IsCluster: config.Config.Redis.IsCluster,
	})
}
func NewRedis(c Config) *Redis {
	if len(c.Addrs) == 0 {
		return nil
	}
	r := &Redis{}
	if len(c.Addrs) == 1 {
		r.single = redis.NewClient(
			&redis.Options{
				Addr:         c.Addrs[0], // use default Addr
				Password:     c.Pwd,      // no password set
				DB:           c.DB,       // use default DB
				DialTimeout:  3 * time.Second,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
			})
		if err := r.single.Ping().Err(); err != nil {
			log.Println(err.Error())
			fmt.Println(err.Error())
			return nil
		}
		fmt.Println("Connected to Redis!")
		// r.single.Do("CONFIG", "SET", "notify-keyspace-events", "AKE")
		r.clusterMode = false
		if c.IsCluster {

			slots, err := r.single.ClusterSlots().Result()
			if err != nil {
				log.Println(err.Error())
			}

			var addrs []string
			for _, slot := range slots {
				for _, node := range slot.Nodes {
					addrs = append(addrs, node.Addr)
				}
			}

			//log.Println("Redis Adds: %v", addrs)

			if len(addrs) > 1 {
				r.clusterMode = true

				r.cluster = redis.NewClusterClient(
					&redis.ClusterOptions{
						Addrs:        c.Addrs,
						Password:     c.Pwd,
						DialTimeout:  3 * time.Second,
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 5 * time.Second,
					})
				// r.cluster.Do("CONFIG", "SET", "notify-keyspace-events", "AKE")

			}
		}
		return r
	}

	r.cluster = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Password:     c.Pwd,
			DialTimeout:  3 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	if err := r.cluster.Ping().Err(); err != nil {
		log.Println(err.Error())
	}
	// r.cluster.Do("CONFIG", "SET", "notify-keyspace-events", "AKE")
	r.clusterMode = true
	return r
}

func (r *Redis) Close() {
	if r.single != nil {
		r.single.Close()
	}
	if r.cluster != nil {
		r.cluster.Close()
	}
}

func (r *Redis) Set(k, v string, t time.Duration) error {
	if r.clusterMode {
		return r.cluster.Set(k, v, t).Err()
	}
	return r.single.Set(k, v, t).Err()
}
func (r *Redis) SetJSON(key string, value interface{}, t time.Duration) error {
	encodedValue, _ := json.Marshal(value)
	if r.clusterMode {
		return r.cluster.Set(key, encodedValue, t).Err()
	}
	return r.single.Set(key, encodedValue, t).Err()
}
func (r *Redis) GetJSON(key string, result interface{}) error {
	encodedValue := r.Get(key)
	if encodedValue == "" {
		result = nil
		return nil
	}
	err := json.Unmarshal([]byte(encodedValue), result)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) SetExpireAt(k, v string, t time.Time) error {
	if r.clusterMode {
		r.cluster.Set(k, v, 0).Err()
		return r.cluster.ExpireAt(k, t).Err()
	}
	r.single.Set(k, v, 0).Err()
	return r.single.ExpireAt(k, t).Err()
}

func (r *Redis) Get(k string) string {
	if r.clusterMode {
		return r.cluster.Get(k).Val()
	}
	return r.single.Get(k).Val()
}

func (r *Redis) HSet(k, field string, value interface{}) error {
	if r.clusterMode {
		return r.cluster.HSet(k, field, value).Err()
	}
	return r.single.HSet(k, field, value).Err()
}

func (r *Redis) HKeys(key string) ([]string, error) {
	if r.clusterMode {
		cmd := r.cluster.HKeys(key)
		return cmd.Val(), cmd.Err()
	}
	cmd := r.single.HKeys(key)
	return cmd.Val(), cmd.Err()
}

func (r *Redis) HSetNX(k, field string, value interface{}) error {
	if r.clusterMode {
		return r.cluster.HSetNX(k, field, value).Err()
	}
	return r.single.HSetNX(k, field, value).Err()
}

func (r *Redis) HIncrBy(k, field string, value int64) int64 {
	if r.clusterMode {
		return r.cluster.HIncrBy(k, field, value).Val()
	}
	return r.single.HIncrBy(k, field, value).Val()
}

func (r *Redis) HGet(k, field string) string {
	if r.clusterMode {
		return r.cluster.HGet(k, field).Val()
	}
	return r.single.HGet(k, field).Val()
}

func (r *Redis) HGetAll(k string) map[string]string {
	if r.clusterMode {
		return r.cluster.HGetAll(k).Val()
	}
	return r.single.HGetAll(k).Val()
}

func (r *Redis) HDel(k, field string) error {
	if r.clusterMode {
		return r.cluster.HDel(k, field).Err()
	}
	return r.single.HDel(k, field).Err()
}

func (r *Redis) Expire(k string, t time.Duration) error {
	if r.clusterMode {
		return r.cluster.Expire(k, t).Err()
	}

	return r.single.Expire(k, t).Err()
}

func (r *Redis) HSetTTL(k, field string, value interface{}, t time.Duration) error {
	if r.clusterMode {
		if err := r.cluster.HSet(k, field, value).Err(); err != nil {
			return err
		}
		return r.cluster.Expire(k, t).Err()
	}
	if err := r.single.HSet(k, field, value).Err(); err != nil {
		return err
	}
	return r.single.Expire(k, t).Err()
}

func (r *Redis) Keys(k string) []string {
	if r.clusterMode {
		return r.cluster.Keys(k).Val()
	}
	return r.single.Keys(k).Val()
}

func (r *Redis) Del(k string) error {
	if r.clusterMode {
		return r.cluster.Del(k).Err()
	}
	return r.single.Del(k).Err()
}

func (r *Redis) RPush(k string, value interface{}) error {
	if r.clusterMode {
		return r.cluster.RPush(k, value).Err()
	}
	return r.single.RPush(k, value).Err()
}

func (r *Redis) LPush(k string, value interface{}) error {
	// logger.Infof("LPUSH: %s --> %v", k, value)
	if r.clusterMode {
		return r.cluster.LPush(k, value).Err()
	}
	return r.single.LPush(k, value).Err()
}

func (r *Redis) LRange(k string, start, stop int64) []string {
	if r.clusterMode {
		return r.cluster.LRange(k, start, stop).Val()
	}
	return r.single.LRange(k, start, stop).Val()
}

func (r *Redis) LRem(k string, value interface{}) error {
	// logger.Infof("LREM: %s --> %v", k, value)
	if r.clusterMode {
		return r.cluster.LRem(k, 0, value).Err()
	}
	return r.single.LRem(k, 0, value).Err()
}

func (r *Redis) BRPop(duration time.Duration, key string) []string {
	if r.clusterMode {
		return r.cluster.BRPop(duration, key).Val()
	}
	return r.single.BRPop(duration, key).Val()
}

func (r *Redis) RPop(key string) string {
	if r.clusterMode {
		return r.cluster.RPop(key).Val()
	}
	return r.single.RPop(key).Val()
}

func (r *Redis) LPop(key string) string {
	if r.clusterMode {
		return r.cluster.LPop(key).Val()
	}
	return r.single.LPop(key).Val()
}

func (r *Redis) LLen(key string) int64 {
	if r.clusterMode {
		return r.cluster.LLen(key).Val()
	}
	return r.single.LLen(key).Val()
}

func (r *Redis) SAdd(k string, value interface{}) int64 {
	if r.clusterMode {
		return r.cluster.SAdd(k, value).Val()
	}
	return r.single.SAdd(k, value).Val()
}

func (r *Redis) SIsMember(k string, value interface{}) bool {
	if r.clusterMode {
		return r.cluster.SIsMember(k, value).Val()
	}
	return r.single.SIsMember(k, value).Val()
}

func (r *Redis) SMembers(k string) []string {
	if r.clusterMode {
		return r.cluster.SMembers(k).Val()
	}
	return r.single.SMembers(k).Val()
}

func (r *Redis) SRem(key string, value interface{}) int64 {
	if r.clusterMode {
		return r.cluster.SRem(key, value).Val()
	}
	return r.single.SRem(key, value).Val()
}

func (r *Redis) SPop(key string) string {
	if r.clusterMode {
		return r.cluster.SPop(key).Val()
	}
	return r.single.SPop(key).Val()
}

func (r *Redis) SCard(key string) int64 {
	if r.clusterMode {
		return r.cluster.SCard(key).Val()
	}
	return r.single.SCard(key).Val()
}

func (r *Redis) ZAdd(key string, member *redis.Z) int64 {
	if r.clusterMode {
		return r.cluster.ZAdd(key, member).Val()
	}
	return r.single.ZAdd(key, member).Val()
}

func (r *Redis) ZRem(key string, member interface{}) int64 {
	if r.clusterMode {
		return r.cluster.ZRem(key, member).Val()
	}
	return r.single.ZRem(key, member).Val()
}

func (r *Redis) BZPopMin(duration time.Duration, key string) *redis.ZWithKey {
	if r.clusterMode {
		return r.cluster.BZPopMin(duration, key).Val()
	}
	return r.single.BZPopMin(duration, key).Val()
}

func (r *Redis) ZIncrBy(key string, increment float64, member string) float64 {
	if r.clusterMode {
		return r.cluster.ZIncrBy(key, increment, member).Val()
	}
	return r.single.ZIncrBy(key, increment, member).Val()
}

func (r *Redis) ZInterStore(destination string, store *redis.ZStore) int64 {
	if r.clusterMode {
		return r.cluster.ZInterStore(destination, store).Val()
	}
	return r.single.ZInterStore(destination, store).Val()
}

func (r *Redis) ZRange(key string, start, stop int64) []string {
	if r.clusterMode {
		return r.cluster.ZRange(key, start, stop).Val()
	}
	return r.single.ZRange(key, start, stop).Val()
}

func (r *Redis) ZRangeWithScores(key string, start, stop int64) []redis.Z {
	if r.clusterMode {
		return r.cluster.ZRangeWithScores(key, start, stop).Val()
	}
	return r.single.ZRangeWithScores(key, start, stop).Val()
}

func (r *Redis) Exists(keys string) int64 {
	if r.clusterMode {
		return r.cluster.Exists(keys).Val()
	}
	return r.single.Exists(keys).Val()
}

func (r *Redis) Incr(key string) int64 {
	if r.clusterMode {
		return r.cluster.Incr(key).Val()
	}
	return r.single.Incr(key).Val()
}

func (r *Redis) Decr(key string) int64 {
	if r.clusterMode {
		return r.cluster.Decr(key).Val()
	}
	return r.single.Decr(key).Val()
}

func (r *Redis) TxPipeline() redis.Pipeliner {
	if r.clusterMode {
		return r.cluster.TxPipeline()
	}
	return r.single.TxPipeline()
}

func (r *Redis) Watch(fn func(*redis.Tx) error, keys ...string) error {
	if r.clusterMode {
		return r.cluster.Watch(fn, keys...)
	}
	return r.single.Watch(fn, keys...)
}

func (r *Redis) SetNX(key string, value string, duration time.Duration) bool {
	if r.clusterMode {
		return r.cluster.SetNX(key, value, duration).Val()
	}
	return r.single.SetNX(key, value, duration).Val()
}

func (r *Redis) TTL(key string) time.Duration {
	if r.clusterMode {
		return r.cluster.TTL(key).Val()
	}
	return r.single.TTL(key).Val()
}

func (r *Redis) ReleaseLock(lockname, identifier string) bool {
	lockname = "lock:" + lockname
	var flag = true
	for flag {
		err := r.Watch(func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()
			if tx.Get(lockname).Val() == identifier {
				pipe.Del(lockname)
				if _, err := pipe.Exec(); err != nil {
					return err
				}
				flag = true
				return nil
			}

			tx.Unwatch()
			flag = false
			return nil
		}, lockname)

		if err != nil {
			// logger.Errorf("watch failed in ReleaseLock, err is: ", err)
			return false
		}
	}
	return true
}

func (r *Redis) Publish(channel string, data interface{}) error {
	//return fmt.Errorf("return err for test direct flow push")
	if r.cluster != nil {
		return r.cluster.Publish(channel, data).Err()
	}
	return r.single.Publish(channel, data).Err()
}

func (r *Redis) Subscribe(channel string) *redis.PubSub {
	if r.cluster != nil {
		return r.cluster.Subscribe(channel)
	}
	return r.single.Subscribe(channel)
}
