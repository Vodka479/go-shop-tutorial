package servers

import (
	"github.com/Vodka479/go-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github.com/Vodka479/go-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github.com/Vodka479/go-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github.com/Vodka479/go-shop-tutorial/modules/monitor/monitorHandlers"
	"github.com/Vodka479/go-shop-tutorial/modules/users/usersHandlers"
	"github.com/Vodka479/go-shop-tutorial/modules/users/usersRepositories"
	"github.com/Vodka479/go-shop-tutorial/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get(",", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
}
