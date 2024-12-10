package router

import (
	"fmt"
	"github.com/SmirnovND/gofermart/internal/controllers"
	"github.com/SmirnovND/gofermart/internal/pkg/auth"
	"github.com/SmirnovND/gofermart/internal/pkg/config"
	"github.com/SmirnovND/gofermart/internal/pkg/container"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Handler(diContainer *container.Container) http.Handler {
	var db *sqlx.DB
	var cf *config.Config
	var AuthController *controllers.AuthController
	var OrderController *controllers.OrderController
	err := diContainer.Invoke(func(
		d *sqlx.DB,
		c *config.Config,
		authControl *controllers.AuthController,
		orderControl *controllers.OrderController,
	) {
		db = d
		cf = c
		AuthController = authControl
		OrderController = orderControl
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Post("/api/user/register", AuthController.HandleRegisterJSON)
	r.Post("/api/user/login", AuthController.HandleLoginJSON)

	r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(cf, http.HandlerFunc(OrderController.HandleOrdersUpload)).ServeHTTP(w, r)
	})
	r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(cf, http.HandlerFunc(OrderController.HandleListUserOrders)).ServeHTTP(w, r)
	})

	healthcheckController := controllers.NewHealthcheckController(db)
	r.Get("/ping", healthcheckController.HandlePing)

	// Обработчик для неподходящего метода (405 Method Not Allowed)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Обработчик для несуществующих маршрутов (404 Not Found)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}
