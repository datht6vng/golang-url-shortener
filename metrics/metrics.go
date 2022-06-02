package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	totalRequests  prometheus.Counter
	genUrlRequests *prometheus.CounterVec
	getUrlRequests *prometheus.CounterVec
}

func (this *Metrics) Init() *Metrics {
	this.totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Number of requests call to server.",
		},
	)
	this.genUrlRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gen_url_requests",
			Help: "Number of gen short url requests.",
		},
		[]string{"user"},
	)
	this.getUrlRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_url_requests",
			Help: "Number of get shor url requests.",
		},
		[]string{"url"},
	)
	prometheus.MustRegister(this.totalRequests)
	prometheus.MustRegister(this.genUrlRequests)
	prometheus.MustRegister(this.getUrlRequests)
	return this
}
func (this *Metrics) IncreaseTotalRequests() {
	this.totalRequests.Inc()
}
func (this *Metrics) AddTotalRequests(x float64) {
	this.totalRequests.Add(x)
}
func (this *Metrics) IncreaseGenUrlRequests(key string) {
	this.genUrlRequests.WithLabelValues(key).Inc()
}
func (this *Metrics) AddGenUrlRequests(key string, x float64) {
	this.genUrlRequests.WithLabelValues(key).Add(x)
}
func (this *Metrics) IncreaseGetUrlRequests(key string) {
	this.getUrlRequests.WithLabelValues(key).Inc()
}
func (this *Metrics) AddGetUrlRequests(key string, x float64) {
	this.getUrlRequests.WithLabelValues(key).Add(x)
}
