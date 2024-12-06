package container

import (
	"github.com/SmirnovND/gofermart/internal/controllers"
	"github.com/SmirnovND/gofermart/internal/pkg/config"
	"github.com/SmirnovND/gofermart/internal/pkg/db"
	"github.com/SmirnovND/gofermart/internal/repo"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/SmirnovND/gofermart/internal/usecase"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
)

// Container - структура контейнера, обертывающая dig-контейнер
type Container struct {
	container *dig.Container
}

func NewContainer() *Container {
	c := &Container{container: dig.New()}
	c.provideDependencies()
	c.provideRepo()
	c.provideService()
	c.provideUsecase()
	c.provideController()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
	// Регистрируем конфигурацию
	c.container.Provide(config.NewConfigCommand)
	c.container.Provide(db.NewDB)
}

func (c *Container) provideUsecase() {
	c.container.Provide(usecase.NewAuthUseCase)
	c.container.Provide(usecase.NewOrderUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewUserRepo)
}

func (c *Container) provideService() {
	c.container.Provide(service.NewAuthService)
	c.container.Provide(service.NewUserService)
	c.container.Provide(service.NewOrderService)
}

func (c *Container) provideController() {
	c.container.Provide(controllers.NewAuthController)
	c.container.Provide(controllers.NewOrderController)
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
