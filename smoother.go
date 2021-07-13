package smoother

import (
	"math"
)

type Method int

const (
	Exponential Method = iota
	DoubleExponential
	FilterFIR
)

type Smoother struct {
	state         float64
	initialState  float64
	previousState interface{}
	coefficients  interface{}
	method        Method
	t             int
}

/*
A Smoother must not be created and used without initialised with MakeSmoother
*/
func MakeSmoother(method Method, coefficients interface{}, initialState float64) *Smoother {
	t := 0
	if !math.IsNaN(initialState) {
		t++
	}
	switch method {
	case Exponential:
		if v, ok := coefficients.(float64); ok {
			return &Smoother{state: initialState, initialState: initialState, t: t, method: method, coefficients: v}
		}
	case DoubleExponential:
		if v, ok := coefficients.([2]float64); ok {
			return &Smoother{state: initialState, initialState: initialState, t: t, method: method, coefficients: [3]float64{v[0], v[1], 0}}
		}
	case FilterFIR:
		if v, ok := coefficients.([]float64); ok {
			N := len(v)
			return &Smoother{state: initialState, initialState: initialState, t: t, method: method, coefficients: coefficients, previousState: make([]float64, N)}
		}
	}
	return nil
}

func (s *Smoother) Get() float64 {
	return s.state
}

func (s *Smoother) Reset() {
	s.ResetWithState(s.initialState)
}
func (s *Smoother) ResetWithState(initialState float64) {
	s.t = 0
	s.state = initialState

	// Reset previous state buffer
	if ps, ok := s.previousState.([]float64); ok {
		for i := range ps {
			ps[i] = 0.0
		}
	}
}

func (s *Smoother) Next(x float64) float64 {
	nextMethods := map[Method]interface{}{
		Exponential:       s.exponentialNext,
		DoubleExponential: s.doubleExponentialNext,
		FilterFIR:         s.filterFirNext,
	}

	if f, ok := nextMethods[s.method]; ok {
		f.(func(float64) float64)(x)
	}

	s.t++
	return s.state
}

/*
	$s_t = s_{t-1} + \alpha (x_t - s_{t-1})$
	ref: https://en.wikipedia.org/wiki/Exponential_smoothing
*/
func (s *Smoother) exponentialNext(x float64) float64 {
	alpha := s.coefficients.(float64)
	if s.t == 0 {
		s.state = x
	} else {
		s.state = s.state + alpha*(x-s.state)
	}
	return s.state
}

/*
	$s_t = \alpha x_t + (1 - \alpha) (s_{t-1} + b_{t-1})$
	$b_t = \beta (s_t - s_{t-1}) + (1 - \beta) b_{t-1}$
	ref: https://en.wikipedia.org/wiki/Exponential_smoothing
*/
func (s *Smoother) doubleExponentialNext(x float64) float64 {
	c := s.coefficients.([3]float64)
	alpha := c[0]
	beta := c[1]
	b := c[2]
	if s.t == 0 {
		s.state = x
		b = 0
	} else {
		previousState := s.state
		s.state = alpha*x + (1-alpha)*(s.state+b)
		b = beta*(s.state-previousState) + (1-beta)*b
	}
	s.coefficients = [3]float64{alpha, beta, b}
	return s.state
}

func (s *Smoother) filterFirNext(x float64) float64 {
	coefficients := s.coefficients.([]float64)
	ps := s.previousState.([]float64)

	// Update previousSate buffer
	for i := len(ps) - 1; i > 0; i-- {
		ps[i] = ps[i-1]
	}
	ps[0] = x

	// Apply coefficients
	value := 0.0
	for i, c := range coefficients {
		value += ps[i] * c
	}

	// Update state and return
	s.state = value
	return s.state
}
