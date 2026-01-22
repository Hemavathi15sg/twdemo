package telemetry

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request duration histogram
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_milliseconds",
			Help:    "HTTP request duration in milliseconds",
			Buckets: []float64{10, 25, 50, 75, 100, 150, 200, 300, 500, 1000},
		},
		[]string{"endpoint", "method", "status_code"},
	)

	// HTTP request counter
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoint", "method", "status_code"},
	)

	// Grade letter distribution - Figma design tokens
	GradeLetterDistribution = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grade_letter_distribution",
			Help: "Distribution of grades by letter (A+, A, B, etc.) - Figma tokens",
		},
		[]string{"letter_grade", "grade_color"},
	)

	// Grade status distribution - Figma alert states
	GradesCalculatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grades_calculated_total",
			Help: "Total grades calculated by status (success, info, warning, critical) - Figma states",
		},
		[]string{"grade_status", "letter_grade"},
	)

	// Grade curve usage tracking
	GradesWithCurveTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grades_with_curve_total",
			Help: "Total grades calculated with curve applied",
		},
		[]string{"curve_applied"},
	)

	// Cache metrics
	RedisCacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_hits",
			Help: "Total number of Redis cache hits",
		},
	)

	RedisCacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_misses",
			Help: "Total number of Redis cache misses",
		},
	)

	// Performance SLA metrics - TEC-31
	GradeCalculationUnder200ms = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "grade_calculation_under_200ms_total",
			Help: "Total grade calculations completed under 200ms (TEC-31 SLA)",
		},
	)

	GradeCalculationOver200ms = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "grade_calculation_over_200ms_total",
			Help: "Total grade calculations that exceeded 200ms (TEC-31 SLA violation)",
		},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(endpoint, method, statusCode string, duration time.Duration) {
	durationMs := float64(duration.Milliseconds())

	HttpRequestDuration.WithLabelValues(endpoint, method, statusCode).Observe(durationMs)
	HttpRequestsTotal.WithLabelValues(endpoint, method, statusCode).Inc()

	// Track TEC-31 SLA compliance for grade calculation endpoint
	if endpoint == "/api/grades/calculate" {
		if durationMs < 200 {
			GradeCalculationUnder200ms.Inc()
		} else {
			GradeCalculationOver200ms.Inc()
		}
	}
}

// RecordGradeCalculation records grade calculation metrics with Figma tokens
func RecordGradeCalculation(letterGrade, gradeColor, gradeStatus string, curveApplied bool) {
	// Update grade distribution
	GradeLetterDistribution.WithLabelValues(letterGrade, gradeColor).Set(1)

	// Update grade status distribution (Figma alert states)
	GradesCalculatedTotal.WithLabelValues(gradeStatus, letterGrade).Inc()

	// Update curve usage
	curveStr := "false"
	if curveApplied {
		curveStr = "true"
	}
	GradesWithCurveTotal.WithLabelValues(curveStr).Inc()
}

// RecordCacheHit records a cache hit
func RecordCacheHit() {
	RedisCacheHits.Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss() {
	RedisCacheMisses.Inc()
}
