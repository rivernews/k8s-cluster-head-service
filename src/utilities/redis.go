package utilities

import (
	"log"
	"net/url"

	"github.com/gomodule/redigo/redis"
)

// RedisPool - a connection pool for redigo client
// https://github.com/gomodule/redigo
var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		// refer to
		// https://godoc.org/github.com/gomodule/redigo/redis#Pool
		address, password := parseRedisURL()
		c, err := redis.Dial("tcp", address, redis.DialDatabase(GetRedisDB()), redis.DialPassword(password))
		if err != nil {
			return nil, err
		}
		// TODO: remove this
		// if _, err := c.Do("AUTH", password); err != nil {
		// 	c.Close()
		// 	return nil, err
		// }
		// if _, err := c.Do("SELECT", 0); err != nil {
		// 	c.Close()
		// 	return nil, err
		// }
		return c, nil
	},
}

// parse credentials from redis URL
// to accomodate flexible use in redis client
// https://github.com/go-redis/redis/issues/129#issuecomment-118889469
func parseRedisURL() (string, string) {
	redisURLString, _ := GetRedisURL()
	redisURL, redisURLParseError := url.Parse(redisURLString)
	if redisURLParseError != nil {
		log.Fatalln("Cannot parse redis URL")
	}

	redisAddress := redisURL.Host
	redisPassword, _ := redisURL.User.Password()

	return redisAddress, redisPassword
}
