package heads

import (
	"math"
)

type Universe int
type Value float64

func (value Value) Valid() bool {
	return !math.IsNaN(float64(value))
}

var INVALID = Value(math.NaN())
