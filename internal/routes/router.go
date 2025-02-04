package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justaskz/infra-app/internal/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Init() *gin.Engine {
	router := gin.Default()

	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.Use(router)

	router.Use(otelgin.Middleware("go-infra-app"))

	router.GET("/", handlers.StatusHandler)
	router.GET("/health", handlers.HealthHandler)
	router.GET("/memoryload", handlers.MemoryLoadHandler)
	router.GET("/echo", handlers.EchoHandler)

	router.GET("/metrics2", gin.WrapH(promhttp.Handler()))

	return router
}
