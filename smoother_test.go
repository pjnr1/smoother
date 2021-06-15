package smoother

import (
	"testing"
)

func TestExponential_GetAndNext(t *testing.T) {
	initValue := 0.0
	smoother := MakeSmoother(Exponential, 0.5, initValue)

	if got := smoother.Get(); got != initValue {
		t.Errorf("smoother.Get() with initialState=%f and t=0, returned %f", initValue, got)
	}

	expected := 0.5
	if got := smoother.Next(1.0); got != expected {
		t.Errorf("smoother.Next(1.0) with initialState=%f and t=0, returned %f", expected, got)
	}
}
