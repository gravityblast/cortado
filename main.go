package main

import (
  "log"
  "github.com/pilu/traffic"
  "github.com/garyburd/redigo/redis"
)

const VERSION = "0.1.0"

var (
  router    *traffic.Router
  db        redis.Conn
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
  var err error
  host, password := redisSettings()
  db, err = redis.Dial("tcp", host)

  if err != nil {
    log.Fatal(err)
  }

  if len(password) > 0 {
    if _, err := db.Do("AUTH", password); err != nil {
      log.Fatal(err)
    }
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
