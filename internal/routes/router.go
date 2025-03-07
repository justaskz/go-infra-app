package routes

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/justaskz/infra-app/internal/handlers"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
// 	var shutdownFuncs []func(context.Context) error

// 	// shutdown calls cleanup functions registered via shutdownFuncs.
// 	// The errors from the calls are joined.
// 	// Each registered cleanup will be invoked once.
// 	shutdown = func(ctx context.Context) error {
// 		var err error
// 		for _, fn := range shutdownFuncs {
// 			err = errors.Join(err, fn(ctx))
// 		}
// 		shutdownFuncs = nil
// 		return err
// 	}

// 	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
// 	handleErr := func(inErr error) {
// 		err = errors.Join(inErr, shutdown(ctx))
// 	}

// 	// Set up propagator.
// 	prop := newPropagator()
// 	otel.SetTextMapPropagator(prop)

// 	// Set up trace provider.
// 	// tracerProvider, err := newTracerProvider()
// 	// if err != nil {
// 	// 	handleErr(err)
// 	// 	return
// 	// }
// 	// shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
// 	// otel.SetTracerProvider(tracerProvider)

// 	meterProvider, err := newMeterProvider(ctx)
// 	if err != nil {
// 		handleErr(err)
// 		return
// 	}
// 	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
// 	otel.SetMeterProvider(meterProvider)

// 	return
// }

// func newPropagator() propagation.TextMapPropagator {
// 	return propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	)
// }

// func newTracerProvider() (*trace.TracerProvider, error) {
// 	traceExporter, err := stdouttrace.New(
// 		stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, err
// 	}

// 	tracerProvider := trace.NewTracerProvider(
// 		trace.WithBatcher(traceExporter,
// 			trace.WithBatchTimeout(time.Second)),
// 	)
// 	return tracerProvider, nil
// }

// func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
// 	metricExporter, err := otlpmetrichttp.New(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meterProvider := metric.NewMeterProvider(
// 		metric.WithReader(metric.NewPeriodicReader(metricExporter,
// 			metric.WithInterval(3*time.Second))),
// 	)
// 	return meterProvider, nil
// }

func initMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	interval := metric.WithInterval(3 * time.Second)
	reader := metric.NewPeriodicReader(exporter, interval)
	readerOptions := metric.WithReader(reader)

	serviceName := semconv.ServiceNameKey.String("go-infra-app")
	resource := resource.NewWithAttributes(semconv.SchemaURL, serviceName)
	attributesOptions := metric.WithResource(resource)
	meterProvider := metric.NewMeterProvider(readerOptions, attributesOptions)

	return meterProvider, nil
}

func Init() *gin.Engine {
	router := gin.Default()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	meterProvider, err := initMeterProvider(ctx)
	if err != nil {
		log.Fatalf("failed to initialize meter provider: %v", err)
	}
	defer func() { _ = meterProvider.Shutdown(ctx) }()
	otel.SetMeterProvider(meterProvider)

	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.Use(router)

	router.Use(otelgin.Middleware("go-infra-app"))
	router.GET("/", handlers.StatusHandler)
	router.GET("/health", handlers.HealthHandler)
	router.GET("/memoryload", handlers.MemoryLoadHandler)
	router.GET("/echo", handlers.EchoHandler)
	router.GET("/counter", handlers.CounterHandler)

	// router.GET("/metrics2", gin.WrapH(promhttp.Handler()))

	return router
}
