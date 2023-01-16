package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/xiabing/goat.git/prom"
)

var (
	request = prom.NewPromVec("cobra_http").
		Counter("request_total", "Total number of request made.", []string{"code", "method", "path"}).
		Histogram("response_duration_seconds", "Bucketed histogram of api response time duration",
			[]string{"code", "method", "path"},
			[]float64{0.01, 0.02, 0.05, 0.1, 0.5, 1, 2, 5, 10, 15, 20},
		)
)

func Prometheus() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		handlsTime := time.Since(start).Seconds()

		path := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		request.Inc(status, method, path)
		request.HandleTimeWithSeconds(handlsTime, status, method, path)
	}
}
