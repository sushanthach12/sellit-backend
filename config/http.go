package config

import (
	"net/http"
	"time"
)

func GetServerConfig(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
}
