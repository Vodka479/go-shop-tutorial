package monitorHandlers // รับ api request หรือ protocal ต่างๆ ในที่นี้ คือ Http

import (
	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/Vodka479/go-shop-tutorial/modules/entities"
	"github.com/Vodka479/go-shop-tutorial/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMontiorHandler interface { // Handler ให้รับ monitor แค่ Context ของ fiber เท่านั้น
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMontiorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
