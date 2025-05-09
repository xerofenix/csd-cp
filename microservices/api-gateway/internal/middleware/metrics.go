package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// responseWriter wraps fiber.Ctx to implement http.ResponseWriter
type responseWriter struct {
	ctx *fiber.Ctx
}

func (rw *responseWriter) Header() http.Header {
	header := make(http.Header)
	rw.ctx.Response().Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ctx.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ctx.Status(statusCode)
}

// SetupPrometheus initializes Prometheus metrics and registers the /metrics endpoint
func SetupPrometheus(app *fiber.App) {
	// Metrics middleware to track requests and duration
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		status := c.Response().StatusCode()
		httpRequestsTotal.WithLabelValues(
			c.Method(),
			c.Path(),
			strconv.Itoa(status),
		).Inc()
		httpRequestDuration.WithLabelValues(
			c.Method(),
			c.Path(),
			strconv.Itoa(status),
		).Observe(duration)

		return err
	})

	// Metrics endpoint for Prometheus scraping
	app.Get("/metrics", func(c *fiber.Ctx) error {
		// Create response writer
		rw := &responseWriter{ctx: c}

		// Convert fasthttp.Request to http.Request
		var httpReq http.Request
		if err := fasthttpadaptor.ConvertRequest(c.Context(), &httpReq, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "internal_error",
				"message": "Failed to convert request",
			})
		}

		// Add context to http.Request
		ctx := context.Background()
		httpReq = *httpReq.WithContext(ctx)

		// Serve Prometheus metrics
		handler := promhttp.Handler()
		handler.ServeHTTP(rw, &httpReq)
		return nil
	})
}
