package main

import (
  "log"
  "time"
  "github.com/pilu/traffic"
  "github.com/garyburd/redigo/redis"
)

const VERSION = "0.1.0"

var (
  router    *traffic.Router
  dbPool    *redis.Pool
  settings  map[string]string
)

func init() {
  initSettings()
  initRedis()
  initRouter()
}

func initSettings() {
  settings = make(map[string]string)
  settings["base_url"] = getConfig("base_url")
  settings["info"]     = getConfig("info")
}

func initRedis() {
  host, password := redisSettings()
  dbPool = &redis.Pool{
    MaxIdle: 10,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) {
      c, err := redis.Dial("tcp", host)
      if err != nil {
        return nil, err
      }

      if len(password) > 0 {
        if _, err := c.Do("AUTH", password); err != nil {
          log.Fatal(err)
        }
      }

      return c, err
    },
    TestOnBorrow: func(c redis.Conn, t time.Time) error {
      _, err := c.Do("PING")
      return err
    },
  }
}

func initRouter() {
  router = traffic.New()
  router.AddBeforeFilter(SetHeaders)
  router.Get("/", IndexHandler)
  router.Post("/", CreateHandler)
  router.Get("/:shorty\\+", InfoHandler)
  router.Get("/:shorty", RedirectHandler)
}

func main() {
  router.Run()
}
