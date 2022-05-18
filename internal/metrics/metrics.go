package metrics

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func InstallRoute(r *gin.Engine) {
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())
}

const (
	ns = "mirage"
)

var (
	OptimizersRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: "optimizer",
		Name:      "running",
		Help:      "Number of images being currently optimized",
	})
	RequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of requested images",
	})
	RequestCachedCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: "http",
		Name:      "requests_cached",
		Help:      "Total number of requested images found in cache",
	})
)
