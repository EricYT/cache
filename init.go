package cache

import (
	"strings"
	"time"
  "os"
  "path"

  "fmt"
)

func init() {
  fmt.Println("---------------------> cache init 1")

  cwd, _ := os.Getwd()
  fmt.Println("---------------------> cache init cwd:", cwd)
  fmt.Println("---------------------> cache init cwd:", path.Join(cwd, "cache.conf"))

  // Load config
  var err error
  Config, err = LoadConfig(path.Join(cwd, "cache.conf"))
  if err != nil {
    fmt.Println("------> cache init error:", err)
    panic("Init config error")
  }

}

func init() {
  fmt.Println("---------------------> cache init 2")

	// Set the default expiration time.
	defaultExpiration := time.Hour // The default for the default is one hour.
	if expireStr, found := Config.String("cache.expires"); found {
		var err error
		if defaultExpiration, err = time.ParseDuration(expireStr); err != nil {
			panic("Could not parse default cache expiration duration " + expireStr + ": " + err.Error())
		}
	}

	// make sure you aren't trying to use both memcached and redis
	if Config.BoolDefault("cache.memcached", false) && Config.BoolDefault("cache.redis", false) {
		panic("You've configured both memcached and redis, please only include configuration for one cache!")
	}

	// Use memcached?
	if Config.BoolDefault("cache.memcached", false) {
		hosts := strings.Split(Config.StringDefault("cache.hosts", ""), ",")
		if len(hosts) == 0 {
			panic("Memcache enabled but no memcached hosts specified!")
		}

		Instance = NewMemcachedCache(hosts, defaultExpiration)
		return
	}

fmt.Println("---------------------> cache init 2.5")
	// Use Redis (share same config as memcached)?
	if Config.BoolDefault("cache.redis", false) {
fmt.Println("---------------------> cache init 3")
		hosts := strings.Split(Config.StringDefault("cache.hosts", ""), ",")
		if len(hosts) == 0 {
			panic("Redis enabled but no Redis hosts specified!")
		}
		if len(hosts) > 1 {
			panic("Redis currently only supports one host!")
		}
		password := Config.StringDefault("cache.redis.password", "")
		Instance = NewRedisCache(hosts[0], password, defaultExpiration)
		return
	}

	// By default, use the in-memory cache.
  Instance = NewInMemoryCache(defaultExpiration)
}
