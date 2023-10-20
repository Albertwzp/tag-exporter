package main

import (
	"tag-exporter/collect"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*func recordMetrics() {
	go func() {
		for {
			apps, _ :=
			for _, app := range apps {
			}
			versionUp.Inc()
			versionDown.Add(11)
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	versionUp = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "app_version_on",
		Help:        "The up app version",
		ConstLabels: map[string]string{},
	})
	versionDown = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "app_version_down",
		Help:        "The down app version",
		ConstLabels: map[string]string{},
	})
)*/

func main() {
	prometheus.MustRegister(collect.NewAppVersionCollector())
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
