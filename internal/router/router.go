package router

import (
	"fmt"
	"github.com/SmirnovND/gofermart/internal/controllers"
	"github.com/SmirnovND/gofermart/internal/pkg/container"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Handler(diContainer *container.Container) http.Handler {
	var db *sqlx.DB
	var authController *controllers.AuthController
	err := diContainer.Invoke(func(d *sqlx.DB, controller *controllers.AuthController) {
		db = d
		authController = controller
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Post("/api/user/register", authController.HandleRegisterJSON)

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
