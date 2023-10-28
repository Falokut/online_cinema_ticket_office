package metrics

import (
	"strconv"

	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncCacheHits(status int, method string)
	IncCacheMiss(status int, method string)
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	CacheHits *prometheus.CounterVec
	CacheMiss *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(address string, name string, logger logging.Logger) (Metrics, error) {
	var metr PrometheusMetrics
	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}

	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}

	metr.CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_cache_hits",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.CacheHits); err != nil {
		return nil, err
	}

	metr.CacheMiss = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_cache_miss",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.CacheMiss); err != nil {
		return nil, err
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return nil, err
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}
	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		if err := router.Start(address); err != nil {
			logger.Fatal(err)
		}
	}()
	return &metr, nil
}

func (metr *PrometheusMetrics) IncHits(status int, method, path string) {
	metr.HitsTotal.Inc()
	metr.CacheHits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metr *PrometheusMetrics) IncCacheHits(status int, method string) {
	metr.HitsTotal.Inc()
	metr.CacheHits.WithLabelValues(strconv.Itoa(status), method).Inc()
}

func (metr *PrometheusMetrics) IncCacheMiss(status int, method string) {
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(strconv.Itoa(status), method).Inc()
}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metr.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}
