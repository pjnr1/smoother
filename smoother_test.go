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

func TestMultipleInstances(t *testing.T) {
	s1 := MakeSmoother(Exponential, 0.5, 0.0)
	s2 := MakeSmoother(Exponential, 0.5, 1.0)

	for i := 1; i < 10; i++ {
		s1.Next(float64(i))
		s2.Next(float64(i))
	}

	if s1.Get() == s2.Get() {
		t.Error("Values should not be equal!")
	}

}
