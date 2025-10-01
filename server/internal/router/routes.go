package router

import (
	"my_project/server/internal/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitializeRoutes(h *server.Hertz, taskService *service.TaskService) {
	//h.POST("/probes/icmp", taskService.HandleICMPTask()) // ICMP 下发
	//h.POST("/probes/tcp", taskService.HandleTCPTask())   // TCP 下发

	//h.DELETE("/probes/icmp/:ip", taskService.HandlePauseICMP())
	//h.DELETE("/probes/tcp/:ip/:port", taskService.HandlePauseTCP())
}
