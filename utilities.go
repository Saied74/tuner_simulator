package main

import (
	"log"
	"math"
)

const (
	cMin = 5.0e-12
	lMin = 10.0e-9
	z0   = 50.0
)

func makeSmith(home string) *smith {
	return &smith{
		outputFile:     home + "data.csv",
		minMaxFile:     home + "minMax.csv",
		point0:         &smithPoint{},
		point1:         &smithPoint{},
		gainTol:        0.15,
		phaseTol:       8.0,
		threshold:      1.2,
		iteration:      2,
		normalize:      "not",
		options:        "MMFitLC",
		minMax:         make(map[string]*lcMinMax),
		baseMaxSeries1: &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
	}
}

func (s *smith) resetSmith(home string) *smith {
	ss := makeSmith(home)
	ss.outputFile = s.outputFile
	ss.minMaxFile = s.minMaxFile
	ss.gainTol = s.gainTol
	ss.phaseTol = s.phaseTol
	ss.threshold = s.threshold
	ss.iteration = s.iteration
	ss.normalize = s.normalize
	ss.options = s.options
	ss.minMax = s.minMax
	return ss
}

func (s *smith) trueCalc() {
	s.locate()
	switch s.region {
	case 1:
		s.rotateRight()
	case 2:
		s.rotateLeft()
	default:
		log.Fatal("bad region")
	}
}

func (s *smith) locate() {
	theta := 2.0 * (s.theta / 360.0) * math.Pi
	s.point0.gammaReal = s.gamma * math.Cos(theta)
	s.point0.gammaImag = s.gamma * math.Sin(theta)
	gammaSq := s.gamma * s.gamma
	denom := 1 + gammaSq - 2*s.point0.gammaReal
	s.point0.r = (1 - gammaSq) / denom
	s.point0.x = (2 * s.point0.gammaImag) / denom
	zSq := s.point0.r*s.point0.r + s.point0.x*s.point0.x
	s.point0.g = s.point0.r / zSq
	s.point0.b = -s.point0.x / zSq
	if s.point0.r > 1.0 {
		s.region = 1
		return
	}
	if s.point0.x > 0.0 && s.point0.r < zSq {
		s.region = 1
		return
	}
	s.region = 2
	return
}

func (s *smith) reLocate() {
	zSq := s.point0.r*s.point0.r + s.point0.x*s.point0.x
	if s.point0.r > 1.0 {
		s.region = 1
		return
	}
	if s.point0.x > 0.0 && s.point0.r < zSq {
		s.region = 1
		return
	}
	s.region = 2
	return
}

func (s *smith) rotateRight() {
	s.point1.gammaReal = (1.0 - s.point0.g) / (1.0 + 3*s.point0.g)
	s.point1.gammaImag = -math.Sqrt(s.point1.gammaReal - s.point1.gammaReal*s.point1.gammaReal)
	s.calcEndPoint()
	s.parallelSuscep = s.point1.b - s.point0.b
	s.parallelReact = -1.0 / s.parallelSuscep
	s.seriesReact = -s.point1.x
	s.seriesSuscep = -1.0 / s.seriesReact
}

func (s *smith) rotateLeft() {
	s.point1.gammaReal = (s.point0.r - 1.0) / (3*s.point0.r + 1)
	s.point1.gammaImag = math.Sqrt(-s.point1.gammaReal - s.point1.gammaReal*s.point1.gammaReal)
	s.calcEndPoint()
	s.seriesReact = s.point1.x - s.point0.x
	s.seriesSuscep = -1.0 / s.seriesReact
	s.parallelSuscep = -s.point1.b
	s.parallelReact = -1.0 / s.parallelSuscep
}

func (s *smith) calcEndPoint() {
	gammaSq := s.point1.gammaReal*s.point1.gammaReal + s.point1.gammaImag*s.point1.gammaImag
	denom := 1 + gammaSq - 2*s.point1.gammaReal
	s.point1.r = (1 - gammaSq) / denom
	s.point1.x = 2 * s.point1.gammaImag / denom
	zSq := s.point1.r*s.point1.r + s.point1.x*s.point1.x
	s.point1.g = s.point1.r / zSq
	s.point1.b = -s.point1.x / zSq
}

//// TODO: this needs to fixed to account for both gamma and theta tolerance
func (s *smith) calcTolerance() {
	s.tolerance = []sensitivity{}
	for _, t := range tolerance {
		sen := sensitivity{}
		s.gamma += t * s.gamma

		s.locate()
		switch s.region {
		case 1:
			s.rotateRight()
		case 2:
			s.rotateLeft()
		default:
			log.Fatal("bad region")
		}
		sen.region = s.region
		sen.parallelReactance = s.parallelReact
		sen.seriesReactance = s.seriesReact
		s.tolerance = append(s.tolerance, sen)

		s.gamma = s.gammaTemp
		s.theta += t * s.theta
		s.locate()
		switch s.region {
		case 1:
			s.rotateRight()
		case 2:
			s.rotateLeft()
		default:
			log.Fatal("bad region")
		}
		sen.region = s.region
		sen.parallelReactance = s.parallelReact
		sen.seriesReactance = s.seriesReact
		s.tolerance = append(s.tolerance, sen)
		s.theta = s.thetaTemp
		//s.tolerance = append(s.tolerance, float64(s.region), s.parallelReact, s.seriesReact)
	}
}

func (s *smith) bruteIt() (float64, float64) {
	var r, x float64
	var inc = 0.01
	switch s.region {
	case 1:
		b := s.baseMaxSeries1.basePoint.b + s.parallelSuscep
		g := s.baseMaxSeries1.basePoint.g
		r, x = getImp(g, b)
		for i := 0; inc < 10; i++ {
			x += inc * s.seriesReact
			if calcSWR(r, x) < s.threshold {
				s.seriesReact = s.seriesReact * inc * float64(i)
				return r, x
			}
		}
	case 2:
		x = s.baseMaxSeries1.basePoint.x + s.seriesReact
		r = s.baseMaxSeries1.basePoint.r
		g, b := getAdm(r, x)
		for i := 0; i < 10; i++ {
			b += inc * s.parallelSuscep
			r, x = getImp(g, b)
			if calcSWR(r, x) < s.threshold {
				s.parallelReact = s.parallelReact * inc * float64(i)
				return r, x
			}
		}
	}
	return r, x //the program never reaches this point except when it fails
}

func getImp(g, b float64) (float64, float64) {
	ySq := b*b + g*g
	return g / ySq, -b / ySq
}

func getAdm(r, x float64) (float64, float64) {
	zSq := r*r + x*x
	return r / zSq, -x / zSq
}

func calcSWR(r, x float64) float64 {
	denom := (r+1)*(r+1) + x*x
	firstTerm := r*r + x*x - 1.0
	num := math.Sqrt(firstTerm*firstTerm + 4.0*x*x)
	absGamma := num / denom
	swr := (1 + absGamma) / (1 - absGamma)
	return swr
}

func (s *smith) calcFreqs() {
	s.freqs = []float64{}
	for _, freq := range freqList {
		c := -1.0 / (2.0 * math.Pi * freqs[freq] * s.parallelReact * z0)
		l := (s.seriesReact * z0) / (2.0 * math.Pi * freqs[freq])
		s.freqs = append(s.freqs, c, l)
	}
}

func (s *smith) copyExt(e *extreme) {

	e.s = s.s
	e.gamma = s.gamma
	e.theta = s.theta
	e.region = s.region
	//e.basePoint      *smithPoint
	e.basePoint.gammaReal = s.point0.gammaReal
	e.basePoint.gammaImag = s.point0.gammaImag
	e.basePoint.r = s.point0.r
	e.basePoint.x = s.point0.x
	e.basePoint.g = s.point0.g
	e.basePoint.b = s.point0.b
	e.parallelReact = s.parallelReact
	e.parallelSuscep = s.parallelSuscep
	e.seriesReact = s.seriesReact
	e.seriesSuscep = s.seriesSuscep
}

func (s *smith) calcLCValues(f float64) (float64, float64) {
	var l, c float64
	if math.Abs(s.parallelReact) > math.SmallestNonzeroFloat64 {
		c = (-1.0) / (2.0 * math.Pi * f * s.parallelReact * 50)
	}
	if s.seriesReact > math.SmallestNonzeroFloat64 {
		l = (s.seriesReact * 50) / (2 * math.Pi * f)
	}
	return l, c
}

func normalizeLC(lc float64) (float64, string) {
	const uLimit = 1.0e-6
	const nLimit = 1.0e-9
	if lc > uLimit {
		return lc * 1.0e6, "u"
	}
	if lc > nLimit {
		return lc * 1.0e9, "n"
	}
	return lc * 1.0e12, "p"
}

//approximate LC values using baseCap and baseInductor values
func fitLC(lc float64, base []float64) (float64, bool) {
	var biggest float64
	for _, b := range base {
		biggest += b
	}
	if lc > biggest {
		return lc, false
	}
	smallest := base[len(base)-1]
	y := 0.0
	for _, item := range base {
		if lc > item {
			lc -= item
			y += item
		}
	}
	if math.Abs(lc-smallest) < smallest/2 {
		y += smallest
	}
	return y, true
}

//If true run min/max true LC values.  No more processing will proceed.
func (s *smith) stepMMLC() bool {
	switch s.options {
	case "MMLC":
		return true
	}
	return false
}

//If true, approximate LC values using baseCap and baseInductor values
//also set up the condition for running MMFitLC, DelMMFitNotFit, and DelMMFitNotFit
func (s *smith) stepFitLC() bool {
	switch s.options {
	case "FitLC":
		return true
	case "MMFitLC":
		return true
	case "DelFitNotFit":
		return true
	case "DelMMFitNotFit":
		return true
	}
	return false
}

//if true calculate the min/max of the approximated LC values
//no more processing will proceed past this point
func (s *smith) stepMMFitLC() bool {
	switch s.options {
	case "MMFitLC":
		return true
	}
	return false
}

//If true, claculate the difference between approximated and true LC lcValues
//also set up the condition for running stepDelMMFitNotFit
func (s *smith) stepDelFitNotFit() bool {
	switch s.options {
	case "DelFitNotFit":
		return true
	case "DelMMFitNotFit":
		return true
	}
	return false
}

//if true, calculate the difference between minimum and maximum of true and
//approximated LC values.
func (s *smith) stepDelMMFitNotFit() bool {
	switch s.options {
	case "DelMMFitNotFit":
		return true
	}
	return false
}

func (s *smith) calcMinMax(l, c, freq float64, val string) {
	pmm, ok := s.minMax[val]
	if !ok {
		s.minMax[val] = &lcMinMax{
			freq: freq,
			minC: c,
			maxC: c,
			minL: l,
			maxL: l,
		}
	} else {
		mm := *pmm
		if c > cMin && c < mm.minC {
			mm.minC = c
		}
		if c > mm.maxC {
			mm.maxC = c
		}
		if l > lMin && l < mm.minL {
			mm.minL = l
		}
		if l > mm.maxL {
			mm.maxL = l
		}
		s.minMax[val] = &mm
	}
}

func (s *smith) calcZLoad() complex128 {
	r := s.baseMaxSeries1.basePoint.r * z0
	x := s.baseMaxSeries1.basePoint.x * z0
	return complex(r, x)
}

//for region 1 capacitor parallel with the load
func (s *smith) calcRegion1Z(z complex128) complex128 {
	z1 := complex(0, s.parallelReact*z0)
	z2 := (z1 * z) / (z1 + z)
	return z2 + complex(0, s.seriesReact*z0)
}

func (s *smith) calcRegion2Z(z complex128) complex128 {
	z1 := z + complex(0, s.seriesReact)
	z2 := complex(0, s.parallelReact*z0)
	return (z1 * z2) / (z1 + z2)
}
