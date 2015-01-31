package api

import (
  "github.com/gin-gonic/gin"
  "github.com/magneticio/vamp-loadbalancer/haproxy"
  "github.com/magneticio/vamp-loadbalancer/metrics"
  gologger "github.com/op/go-logging"
  "time"
)

// override the standard Gin-Gonic middleware to add the CORS headers
func CORSMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {

    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
  }
}


func HaproxyMiddleware(haConfig *haproxy.Config, haRuntime *haproxy.Runtime) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("haConfig", haConfig)
        c.Set("haRuntime", haRuntime)

    }
}

func LoggerMiddleware(log *gologger.Logger) gin.HandlerFunc {
  return func(c *gin.Context) {
  // Start timer
    start := time.Now()

    // Process request
    c.Next()

    // Stop timer
    end := time.Now()
    latency := end.Sub(start)

    method := c.Request.Method
    statusCode := c.Writer.Status()
    format := "%-5s %-50s %3d %12v "

    switch {
      case statusCode >= 200 && statusCode <= 399:
        log.Notice(format,method, c.Request.URL.Path, statusCode, latency)
      case statusCode >= 400 && statusCode <= 499:
        log.Warning(format,method, c.Request.URL.Path, statusCode, latency)
      default:
        log.Error(format,method, c.Request.URL.Path, statusCode, latency)
    }
  }
}

func SSEMiddleware(SSEBroker *metrics.SSEBroker) gin.HandlerFunc {
  return func(c *gin.Context) {

        c.Set("sseBroker", SSEBroker)
        c.Writer.Header().Set("Content-Type", "text/event-stream")
        c.Writer.Header().Set("Cache-Control", "no-cache")
        c.Writer.Header().Set("Connection", "keep-alive")

  }
}
