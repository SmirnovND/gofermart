package main

import (
	"github.com/SmirnovND/gofermart/internal/middleware"
	"github.com/SmirnovND/gofermart/internal/pkg/compressor"
	"github.com/SmirnovND/gofermart/internal/pkg/config"
	"github.com/SmirnovND/gofermart/internal/pkg/container"
	"github.com/SmirnovND/gofermart/internal/pkg/loggeer"
	"github.com/SmirnovND/gofermart/internal/router"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	diContainer := container.NewContainer()

	var cf config.Config
	var db *sqlx.DB
	diContainer.Invoke(func(c config.Config, d *sqlx.DB) {
		cf = c
		db = d
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	return http.ListenAndServe(cf.GetFlagRunAddr(), middleware.ChainMiddleware(
		router.Handler(db),
		loggeer.WithLogging,
		compressor.WithDecompression,
		compressor.WithCompression,
	))
}
