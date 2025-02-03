package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justaskz/infra-app/internal/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func Init() *gin.Engine {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/", handlers.StatusHandler)
	r.GET("/health", handlers.HealthHandler)
	r.GET("/memoryload", handlers.MemoryLoadHandler)
	r.GET("/echo", handlers.EchoHandler)

	r.GET("/metrics2", gin.WrapH(promhttp.Handler()))

	return r
}
