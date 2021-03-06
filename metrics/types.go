package metrics

import (
	"fmt"
	"sort"
	"time"
)

// Scope is a namespace wrapper around a stats reporter, ensuring that
// all emitted values have a given prefix or set of tags.
type Scope interface {
	// Counter returns the Counter object corresponding to the name.
	Counter(name string) Counter

	// Gauge returns the Gauge object corresponding to the name.
	Gauge(name string) Gauge

	IntegerGauge(name string) IntegerGauge

	// Timer returns the Timer object corresponding to the name.
	Timer(name string) Timer

	// Histogram returns the Histogram object corresponding to the name.
	// To use default value and duration buckets configured for the scope
	// simply pass metrics.DefaultBuckets or nil.
	// You can use metrics.ValueBuckets{x, y, ...} for value buckets.
	// You can use metrics.DurationBuckets{x, y, ...} for duration buckets.
	// You can use metrics.MustMakeLinearValueBuckets(start, width, count) for linear values.
	// You can use metrics.MustMakeLinearDurationBuckets(start, width, count) for linear durations.
	// You can use metrics.MustMakeExponentialValueBuckets(start, factor, count) for exponential values.
	// You can use metrics.MustMakeExponentialDurationBuckets(start, factor, count) for exponential durations.
	Histogram(name string, buckets Buckets) Histogram

	// Tagged returns a new child scope with the given tags and current tags.
	Tagged(tags map[string]string) Scope

	// SubScope returns a new child scope appending a further name prefix.
	SubScope(name string) Scope

	// Capabilities returns a description of metrics reporting capabilities.
	Capabilities() Capabilities
}

// Counter is the interface for emitting counter type metrics.
type Counter interface {
	// Inc increments the counter by a delta.
	Inc(delta int64)
}

// Gauge is the interface for emitting gauge metrics.
type Gauge interface {
	// Update sets the gauges absolute value.
	Update(value float64)
}

type IntegerGauge interface {
	// Update sets the gauges absolute value.
	Update(value int64)

	Inc(value int64)

	Dec(value int64)
}

// Timer is the interface for emitting timer metrics.
type Timer interface {
	// Record a specific duration directly.
	Record(value time.Duration)

	// Start gives you back a specific point in time to report via Stop.
	Start() Stopwatch
}

// Histogram is the interface for emitting histogram metrics
type Histogram interface {
	// RecordValue records a specific value directly.
	// Will use the configured value buckets for the histogram.
	RecordValue(value float64)

	// RecordDuration records a specific duration directly.
	// Will use the configured duration buckets for the histogram.
	RecordDuration(value time.Duration)

	// Start gives you a specific point in time to then record a duration.
	// Will use the configured duration buckets for the histogram.
	Start() Stopwatch
}

// Stopwatch is a helper for simpler tracking of elapsed time, use the
// Stop() method to report time elapsed since its created back to the
// timer or histogram.
type Stopwatch struct {
	start    time.Time
	recorder StopwatchRecorder
}

// NewStopwatch creates a new immutable stopwatch for recording the start
// time to a stopwatch reporter.
func NewStopwatch(start time.Time, r StopwatchRecorder) Stopwatch {
	return Stopwatch{start: start, recorder: r}
}

// Stop reports time elapsed since the stopwatch start to the recorder.
func (sw Stopwatch) Stop() {
	sw.recorder.RecordStopwatch(sw.start)
}

// StopwatchRecorder is a recorder that is called when a stopwatch is
// stopped with Stop().
type StopwatchRecorder interface {
	RecordStopwatch(stopwatchStart time.Time)
}

// Buckets is an interface that can represent a set of buckets
// either as float64s or as durations.
type Buckets interface {
	fmt.Stringer
	sort.Interface

	// AsValues returns a representation of the buckets as float64s
	AsValues() []float64

	// AsDurations returns a representation of the buckets as time.Durations
	AsDurations() []time.Duration
}

// BucketPair describes the lower and upper bounds
// for a derived bucket from a buckets set.
type BucketPair interface {
	LowerBoundValue() float64
	UpperBoundValue() float64
	LowerBoundDuration() time.Duration
	UpperBoundDuration() time.Duration
}

// Capabilities is a description of metrics reporting capabilities.
type Capabilities interface {
	// Reporting returns whether the reporter has the ability to actively report.
	Reporting() bool

	// Tagging returns whether the reporter has the capability for tagged metrics.
	Tagging() bool
}
