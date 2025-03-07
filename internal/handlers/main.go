package handlers

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/justaskz/infra-app/internal/memoryload"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func StatusHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	currentMemLimit := debug.SetMemoryLimit(-1)

	c.JSON(http.StatusOK, gin.H{
		"port":                                os.Getenv("PORT"),
		"GOMEMLIMIT":                          os.Getenv("GOMEMLIMIT"),
		"max_procs":                           runtime.GOMAXPROCS(0),
		"current_mem_limit":                   float64(currentMemLimit),
		"current_mem_limit_gb":                float64(currentMemLimit) / (1024 * 1024 * 1024),
		"heap_system_mem_mb":                  float64(m.HeapSys) / 1024 / 1024,
		"heap_allocated_mem_mb":               float64(m.HeapAlloc) / 1024 / 1024,
		"heap_in_use_mb":                      float64(m.HeapInuse) / 1024 / 1024,
		"heap_released_mem_mb":                float64(m.HeapReleased) / 1024 / 1024,
		"total_mem_in_use_mb":                 float64(m.Sys) / 1024 / 1024,
		"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"),
		"OTEL_METRICS_EXPORTER":               os.Getenv("OTEL_METRICS_EXPORTER"),
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

func EchoHandler(c *gin.Context) {
	message := c.DefaultQuery("message", "no message found")

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": message,
	})
}

var meter = otel.Meter("go-infra-app")
var counter, _ = meter.Int64Counter(
	"metric.attributes.infra_app_counter",
	metric.WithUnit("1"),
	metric.WithDescription("Infra app counter demo"),
)

func CounterHandler(c *gin.Context) {
	counter.Add(c.Request.Context(), 1)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Counter incremented",
	})
}
