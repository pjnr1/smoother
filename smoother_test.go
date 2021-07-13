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

func TestFilterFir_GetAndNext(t *testing.T) {
	values := []float64{1.0, 0.0, 0.0, 0.5, 0.0, 0.0, 1.0, 0.5, 0.0, 0.0, 0.0}
	expect := []float64{1.0, 0.5, 0.25, 0.5, 0.25, 0.125, 1.0, 1.0, 0.5, 0.125, 0.0}

	s := MakeSmoother(FilterFIR, []float64{1.0, 0.5, 0.25}, 0.0)

	for i, v := range values {
		if got := s.Next(v); got != expect[i] {
			t.Errorf("s.Next(%f) = %f, expected %f", v, got, expect[i])
		}
	}
}
