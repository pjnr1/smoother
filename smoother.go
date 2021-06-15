package smoother

import (
	"math"
)

type Method int

const (
	Exponential Method = iota
	DoubleExponential
)

type Smoother struct {
	state        float64
	coefficients interface{}
	method       Method
	t            int
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
			return &Smoother{state: initialState, coefficients: v, t: t, method: method}
		}
	case DoubleExponential:
		if v, ok := coefficients.([2]float64); ok {
			return &Smoother{state: initialState, coefficients: [3]float64{v[0], v[1], 0}, t: t, method: method}
		}
	}
	return nil
}

func (s *Smoother) Get() float64 {
	return s.state
}

func (s *Smoother) Reset() {
	s.t = 0
}

func (s *Smoother) Next(x float64) float64 {
	switch s.method {
	case Exponential:
		s.exponentialNext(x)
	case DoubleExponential:
		s.doubleExponentialNext(x)
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
