package week06

import (
	"math"
	"sync"
	"time"
)

// RollingCounter contains some common aggregation function.
// Each RollingCounter can compute summary statistics of window.
type RollingCounter interface {
	// Min finds the min value within the window.
	Min() float64
	// Max finds the max value within the window.
	Max() float64
	// Avg computes average value within the window.
	Avg() float64
	// Sum computes sum value within the window.
	Sum() float64
	// Increment increments the number in current window.
	Increment(i float64)
	// UpdateMax updates the maximum value in the current windows.
	UpdateMax(n float64)
}

// RollingCounter tracks a window over a bounded number of time buckets.
type rollingCounter struct {
	buckets  map[int64]*window
	timeSpan int64
	mutex    *sync.RWMutex
}

type window struct {
	Value float64
}

// RollingCounterOpts contains the arguments for creating RollingCounter.
type RollingCounterOpts struct {
	TimeSpan int64
}

// NewRollingCounter initializes a rollingCounter struct.
func NewRollingCounter(opts RollingCounterOpts) RollingCounter {
	r := &rollingCounter{
		buckets:  make(map[int64]*window),
		timeSpan: opts.TimeSpan,
		mutex:    &sync.RWMutex{},
	}
	return r
}

func (r *rollingCounter) getCurrentBucket() *window {
	now := time.Now().Unix()
	var bucket *window
	var ok bool

	if bucket, ok = r.buckets[now]; !ok {
		bucket = &window{}
		r.buckets[now] = bucket
	}

	return bucket
}

func (r *rollingCounter) removeOldBuckets() {
	now := time.Now().Unix() - r.timeSpan

	for timestamp := range r.buckets {
		// TODO: configurable rolling window
		if timestamp <= now {
			delete(r.buckets, timestamp)
		}
	}
}

// Increment increments the number in current timeBucket.
func (r *rollingCounter) Increment(i float64) {
	if i == 0 {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

// UpdateMax updates the maximum value in the current bucket.
func (r *rollingCounter) UpdateMax(n float64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	b := r.getCurrentBucket()
	if n > b.Value {
		b.Value = n
	}
	r.removeOldBuckets()
}

// Sum sums the values over the buckets in the last r.timeSpan seconds.
func (r *rollingCounter) Sum() float64 {
	now := time.Now().Unix()
	sum := float64(0)

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for timestamp, bucket := range r.buckets {
		// TODO: configurable rolling window
		if timestamp >= now-r.timeSpan {
			sum += bucket.Value
		}
	}

	return sum
}

// Min returns the minimum value seen in the last r.timeSpan seconds.
func (r *rollingCounter) Min() float64 {
	var min = math.MaxFloat64
	now := time.Now().Unix()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for timestamp, bucket := range r.buckets {
		// TODO: configurable rolling window
		if timestamp >= now-r.timeSpan {
			if bucket.Value < min {
				min = bucket.Value
			}
		}
	}

	return min
}

// Max returns the maximum value seen in the last r.timeSpan seconds.
func (r *rollingCounter) Max() float64 {
	var max float64
	now := time.Now().Unix()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for timestamp, bucket := range r.buckets {
		// TODO: configurable rolling window
		if timestamp >= now-r.timeSpan {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}

	return max
}

// Avg computes average value over the buckets in the last r.timeSpan seconds.
func (r *rollingCounter) Avg() float64 {
	return r.Sum() / float64(r.timeSpan)
}
