package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	totalRequests  prometheus.Gauge
	genURLRequests *prometheus.GaugeVec
	getURLRequests *prometheus.GaugeVec
}

func (this *Metrics) Init() *Metrics {
	this.totalRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_requests",
			Help: "Number of requests call to server.",
		},
	)
	this.genURLRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gen_url_requests",
			Help: "Number of gen short url requests.",
		},
		[]string{"user"},
	)
	this.getURLRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "get_url_requests",
			Help: "Number of get shor url requests.",
		},
		[]string{"url", "user"},
	)
	prometheus.MustRegister(this.totalRequests)
	prometheus.MustRegister(this.genURLRequests)
	prometheus.MustRegister(this.getURLRequests)
	return this
}

func (this *Metrics) IncreaseTotalRequests() {
	this.totalRequests.Inc()
}
func (this *Metrics) ResetTotalRequests() {
	this.totalRequests.Set(0)
}
func (this *Metrics) IncreaseGenURLRequests(key string) {
	this.genURLRequests.WithLabelValues(key).Inc()
}
func (this *Metrics) ResetGenURLRequests(key string) {
	this.genURLRequests.WithLabelValues(key).Set(0)
}
func (this *Metrics) AddGenURLRequests(key string, x float64) {
	this.genURLRequests.WithLabelValues(key).Add(x)
}

func (this *Metrics) IncreaseGetURLRequests(url string, user string) {
	this.getURLRequests.With(prometheus.Labels{
		"url":  url,
		"user": user,
	}).Inc()
}
func (this *Metrics) ResetGetURLRequests(url string, user string) {
	this.getURLRequests.With(prometheus.Labels{
		"url":  url,
		"user": user,
	}).Set(0)
}
