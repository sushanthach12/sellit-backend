package config

import (
	"net/http"
	"time"
)

func GetServerConfig(handler http.Handler, cfg ConfigVars) *http.Server {
	port := ":" + cfg.Port

	return &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
}
