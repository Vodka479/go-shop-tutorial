package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func Newserver(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	//Middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	//Modules
	v1 := s.app.Group("v1")

	modules := InitModule(v1, s, middlewares)

	modules.MonitorModule()
	modules.UsersModule()

	s.app.Use(middlewares.RouterCheck())

	// Graceful Shutdown ปิด func ต่างๆ ก่อนจะปิดตัว app ลง
	c := make(chan os.Signal, 1) // ประกาศตัวรับสัญญาณ channel = c
	signal.Notify(c, os.Interrupt)
	go func() { // go concurrency เจอแล้วทำงานเลย // แบบ serial จะจากบนลงมาล่าง
		_ = <-c
		log.Println("server is shutting down ...")
		_ = s.app.Shutdown()
	}() // anonymous func ต้องใส้ ()

	// Listen to host:port
	log.Printf("server is starting on %v", s.cfg.App().URL())
	s.app.Listen(s.cfg.App().URL())
}
