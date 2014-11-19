package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/pilu/traffic"
)

var validUrlRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+://([^?/#\.\s]+)\.([^?/#\.\s]+)`)

func validUrl(url string) bool {
	return validUrlRegexp.MatchString(url)
}

func getConfig(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("You need to set %s the env variable", key)
	}

	return value
}

func redisSettings() (string, string) {
	redisUrl := getConfig("REDIS_URL")
	redisInfo, _ := url.Parse(redisUrl)
	host := redisInfo.Host

	var password string

	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}

	return host, password
}

func urlHash(url string) string {
	h := sha1.New()
	io.WriteString(h, url)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func SetHeaders(w traffic.ResponseWriter, r *traffic.Request) {
	w.Header().Set("Cortado-Version", VERSION)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func Error(message string, w traffic.ResponseWriter, status int) {
	w.WriteHeader(status)
	body := map[string]string{
		"error": message,
	}

	w.WriteJSON(body)
}

func HandleError(err error, w traffic.ResponseWriter) {
	if e, ok := err.(ShortyNotFound); ok {
		Error(e.Error(), w, http.StatusNotFound)
	} else if e, ok := err.(InvalidUrl); ok {
		Error(e.Error(), w, http.StatusBadRequest)
	} else {
		log.Println(err.Error())
		Error("internal server error", w, http.StatusInternalServerError)
	}
}
