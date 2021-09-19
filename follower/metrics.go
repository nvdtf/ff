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
		"contract", "failed",
	})
)

func RegisterTagMetrics(contract string, failed bool) {
	tagsObserved.WithLabelValues(contract, strconv.FormatBool(failed)).Inc()
}
