package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	totalRequests  prometheus.Gauge
	genUrlRequests *prometheus.GaugeVec
	getUrlRequests *prometheus.GaugeVec
}

func (this *Metrics) Init() *Metrics {
	this.totalRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_requests",
			Help: "Number of requests call to server.",
		},
	)
	this.genUrlRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gen_url_requests",
			Help: "Number of gen short url requests.",
		},
		[]string{"user"},
	)
	this.getUrlRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "get_url_requests",
			Help: "Number of get shor url requests.",
		},
		[]string{"url", "user"},
	)
	prometheus.MustRegister(this.totalRequests)
	prometheus.MustRegister(this.genUrlRequests)
	prometheus.MustRegister(this.getUrlRequests)
	return this
}
func (this *Metrics) IncreaseTotalRequests() {
	this.totalRequests.Inc()
}
func (this *Metrics) ResetTotalRequests() {
	this.totalRequests.Set(0)
}
func (this *Metrics) IncreaseGenUrlRequests(key string) {
	this.genUrlRequests.WithLabelValues(key).Inc()
}
func (this *Metrics) ResetGenUrlRequests(key string) {
	this.genUrlRequests.WithLabelValues(key).Set(0)
}
func (this *Metrics) AddGenUrlRequests(key string, x float64) {
	this.genUrlRequests.WithLabelValues(key).Add(x)
}

func (this *Metrics) IncreaseGetUrlRequests(url string, user string) {
	this.getUrlRequests.With(prometheus.Labels{
		"url":  url,
		"user": user,
	}).Inc()
}
func (this *Metrics) ResetGetUrlRequests(url string, user string) {
	this.getUrlRequests.With(prometheus.Labels{
		"url":  url,
		"user": user,
	}).Set(0)
}
