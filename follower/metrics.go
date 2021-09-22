package follower

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	tagsObserved = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "import_tag_total",
		Help: "The total number of observed known imports",
	}, []string{
		"contract", "address", "failed",
	})

	eventsObserved = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_total",
		Help: "The total number of observed Cadence events",
	}, []string{
		"type",
	})
)

func RegisterImportMetrics(imp CadenceImport, failed bool) {
	tagsObserved.WithLabelValues(imp.Contract, imp.Address, strconv.FormatBool(failed)).Inc()
}
func RegisterEventMetrics(eventType string) {
	eventsObserved.WithLabelValues(eventType).Inc()
}
