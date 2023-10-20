package collect

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type appVersionCollector struct {
	upversionMetric   *prometheus.Desc
	downversionMetric *prometheus.Desc
}

func NewAppVersionCollector() *appVersionCollector {
	return &appVersionCollector{
		prometheus.NewDesc("app_up_version", "The up app version and nums", []string{"ns", "app", "version"}, nil),
		prometheus.NewDesc("app_down_version", "The down app version and nums", []string{"ns", "app", "version"}, nil),
	}
}

func (collect *appVersionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.upversionMetric
	ch <- collect.downversionMetric
}

func (collect *appVersionCollector) Collect(ch chan<- prometheus.Metric) {
	appversion, _ := AppFilter("deploy_type=dep")
	for _, v := range appversion {
		fmt.Printf("%v\n", v)
		if v.downCount > 0 {
			ch <- prometheus.MustNewConstMetric(collect.downversionMetric, prometheus.GaugeValue, float64(v.downCount), v.ns, v.app, v.downImage)
		}
		ch <- prometheus.MustNewConstMetric(collect.upversionMetric, prometheus.GaugeValue, float64(v.upCount), v.ns, v.app, v.upImage)
	}
}
