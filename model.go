package main

import (
  "fmt"
  "github.com/pilu/go-base62"
  "github.com/garyburd/redigo/redis"
)

type ShortyNotFound struct {
  shorty string
}

func (s ShortyNotFound) Error() string {
  return fmt.Sprintf("shorty not found: `%s`", s.shorty)
}

func Shorten(url string) (string, bool, error) {
  hash := urlHash(url)
  shorty, err := FindByHash(hash)

  if err != nil {
    return "", false, err
  }

  if shorty != "" {
    return shorty, true, nil
  }

  id, err := redis.Int(db.Do("INCR", "next_id"))
  if err != nil {
    return "", false, err
  }

  shorty    =  base62.Encode(int(id))
  shortyKey := fmt.Sprintf("shorties:%s", shorty)
  urlKey    := fmt.Sprintf("urls:%s",     hash)

  db.Do("MULTI")
  db.Do("SET", shortyKey, url)
  db.Do("SET", urlKey,    shorty)
  _, err = db.Do("EXEC")

  if err != nil {
    return "", false, err
  }

  return shorty, false, nil
}

func IncrementClicks(shorty string) {
  db.Do("ZINCRBY", "clicks", 1, shorty)
}

func Clicks(shorty string) int {
  count, err := redis.Int(db.Do("ZSCORE", "clicks", shorty))

  if err != nil {
    return 0
  }

  return count
}

func FindByShorty(shorty string) (string, error) {
  key := fmt.Sprintf("shorties:%s", shorty)
  url, err := db.Do("GET", key)

  if err != nil {
    return "", err
  }

  if url == nil {
    return "", ShortyNotFound{ shorty }
  }

  return fmt.Sprintf("%s", url), nil
}

func FindByHash(hash string) (string, error) {
  key := fmt.Sprintf("urls:%s", hash)
  url, err := db.Do("GET", key)

  if err != nil {
    return "", err
  }

  if url == nil {
    return "", nil
  }

  return fmt.Sprintf("%s", url), nil
}

