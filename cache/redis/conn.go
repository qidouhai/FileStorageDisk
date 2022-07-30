package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "123456"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 50,  // 连接池里最多可用的连接数
		MaxActive: 30,  // 同时能够使用的连接数
		IdleTimeout: 300 * time.Second,  // 超时回收
		Dial: func() (redis.Conn, error) { // 用来创建和配置连接的方法
			// 1. 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			// 2. 访问认证
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error { // 定时检查连接是否可用
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},  
	}
}

// 初始化Redis连接池
func init() {
	pool = newRedisPool()
}

// 对外接口： 返回pool
func RedisPool() *redis.Pool {
	return pool
}