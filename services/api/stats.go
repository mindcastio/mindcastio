package main

import (
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

func stats_endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	result, _ := backend.SimpleApiStats()
	backend.Response(w, result)

	// metrics
	metrics.Count("api.total.count", 1)
	metrics.Count("api.stats.count", 1)
	metrics.Histogram("api.stats.duration", (float64)(util.ElapsedTimeSince(start)))
}
