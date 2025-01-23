package handlers

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/justaskz/infra-app/internal/memoryload"
)

func StatusHandler(c *gin.Context) {
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
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func MemoryLoadHandler(c *gin.Context) {
	mc := memoryload.NewMemoryConsumer().ConsumeMemory(1)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"size":   mc.Size(),
	})
}
