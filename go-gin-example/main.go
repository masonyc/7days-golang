package main

import (
	"fmt"
	"net/http"

	"github.com/masonyc/7days-golang/go-gin-example/pkg/settings"
	"github.com/masonyc/7days-golang/go-gin-example/routers"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", settings.HTTPPort),
		Handler:        router,
		ReadTimeout:    settings.ReadTimeout,
		WriteTimeout:   settings.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
