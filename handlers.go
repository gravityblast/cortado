package main

import (
  "fmt"
  "net/http"
  "github.com/pilu/traffic"
)

func IndexHandler(w traffic.ResponseWriter, r *traffic.Request) {
  body := map[string]string {
    "version":  VERSION,
    "base_url": settings["base_url"],
    "info":     settings["info"],
  }

  w.WriteJSON(body)
}

func CreateHandler(w traffic.ResponseWriter, r *traffic.Request) {
  r.ParseForm()
  url := r.PostForm.Get("url")

  shorty, found, err := Shorten(url)

  if err != nil {
    HandleError(err, w)
    return
  }

  shortUrl := fmt.Sprintf("%s/%s", settings["base_url"], shorty)

  body := map[string]string {
    "long_url":   url,
    "short_url":  shortUrl,
    "shorty":     shorty,
  }

  w.Header().Set("Location", shortUrl)
  if found {
    w.WriteHeader(http.StatusFound)
  } else {
    w.WriteHeader(http.StatusCreated)
  }
  w.WriteJSON(body)
}

func RedirectHandler(w traffic.ResponseWriter, r *traffic.Request) {
  shorty := r.Param("shorty")
  url, err := FindByShorty(shorty)
  if err != nil {
    HandleError(err, w)
    return
  }

  IncrementClicks(shorty)
  http.Redirect(w, r.Request, url, http.StatusMovedPermanently)
}

func InfoHandler(w traffic.ResponseWriter, r *traffic.Request) {
  shorty := r.Param("shorty")
  url, err := FindByShorty(shorty)
  if err != nil {
    HandleError(err, w)
    return
  }

  clicks := fmt.Sprintf("%d", Clicks(shorty))

  body := map[string]string {
    "long_url":   url,
    "short_url":  fmt.Sprintf("%s/%s", settings["base_url"], shorty),
    "shorty":     shorty,
    "clicks":     clicks,
  }

  w.WriteJSON(body)
}
