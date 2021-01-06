package week06

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRollingCounter(t *testing.T) {

	counter := NewRollingCounter(RollingCounterOpts{TimeSpan: 10})
	for _, request := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		counter.Increment(request)
		time.Sleep(1 * time.Second)
	}

	assert.Equal(t, float64(45), counter.Sum())
	assert.Equal(t, float64(4.5), counter.Avg())
	assert.Equal(t, float64(1), counter.Min())
	assert.Equal(t, float64(9), counter.Max())

	t.Logf("before: sum is %v, min is %v,max is %v, avg is %v", counter.Sum(), counter.Min(), counter.Max(), counter.Avg())
	time.Sleep(1 * time.Second)
	counter.UpdateMax(20)
	assert.Equal(t, float64(64), counter.Sum())
	assert.Equal(t, float64(6.4), counter.Avg())
	assert.Equal(t, float64(2), counter.Min())
	assert.Equal(t, float64(20), counter.Max())

	t.Logf("after: sum is %v, min is %v,max is %v, avg is %v", counter.Sum(), counter.Min(), counter.Max(), counter.Avg())
}
