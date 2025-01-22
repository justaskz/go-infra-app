package main

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		currentMemLimit := debug.SetMemoryLimit(-1)

		c.JSON(http.StatusOK, gin.H{
			"port":                  os.Getenv("PORT"),
			"GOMEMLIMIT":            os.Getenv("GOMEMLIMIT"),
			"max_procs":             runtime.GOMAXPROCS(0),
			"current_mem_limit":     float64(currentMemLimit),
			"current_mem_limit_gb":  float64(currentMemLimit) / (1024 * 1024 * 1024),
			"heap_system_mem_mb":    float64(m.HeapSys) / 1024 / 1024,
			"heap_allocated_mem_mb": float64(m.HeapAlloc) / 1024 / 1024,
			"heap_in_use_mb":        float64(m.HeapInuse) / 1024 / 1024,
			"heap_released_mem_mb":  float64(m.HeapReleased) / 1024 / 1024,
			"total_mem_in_use_mb":   float64(m.Sys) / 1024 / 1024,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/metrics2", gin.WrapH(promhttp.Handler()))

	r.Run()
}
