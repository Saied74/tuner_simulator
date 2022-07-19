package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
)

const (
	cMin    = 5.0e-12
	lMin    = 10.0e-9
	z0      = 50.0
	epsilon = math.SmallestNonzeroFloat64 + math.SmallestNonzeroFloat64
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
		matchC: []*matchParts{ //ordered from the highest value to the lowest value
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
		},
		matchL: []*matchParts{
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
			&matchParts{},
		},
		vSource: 235.0,
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
	ss.matchC = []*matchParts{ //ordered from the highest value to the lowest value
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
	}
	ss.matchL = []*matchParts{
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
	}
	ss.vSource = s.vSource
	ss.which = s.which
	return ss
}

func (s *smith) trueCalc() {
	s.locate()
	switch s.region {
	case 1:
		s.rotateRight()
	case 2:
		s.rotateLeft()
	case 3:
		s.noParallel()
	case 4:
		s.noSeries()
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
	if math.Abs(s.point0.r-1.0) < epsilon {
		s.region = 3
		return
	}
	if math.Abs(s.point0.g-1.0) < epsilon {
		s.region = 4
		return
	}
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

//noParallel is the case when the original point (point 0) is on the r=1 circle
//and only a series inductance is needed to mach the impedance
func (s *smith) noParallel() {
	s.point1.gammaReal = (1.0 - s.point0.g) / (1.0 + 3*s.point0.g)
	s.point1.gammaImag = -math.Sqrt(s.point1.gammaReal - s.point1.gammaReal*s.point1.gammaReal)
	s.calcEndPoint()
	// s.parallelSuscep = s.point1.b - s.point0.b
	// s.parallelReact = -1.0 / s.parallelSuscep
	s.seriesReact = -s.point1.x
	s.seriesSuscep = -1.0 / s.seriesReact
}

//noSeries is the case when the original point (point 0) is on the g=1 circle
//and only a parallel parallel capacitance is needed to match the impedance
func (s *smith) noSeries() {
	s.point1.gammaReal = (s.point0.r - 1.0) / (3*s.point0.r + 1)
	s.point1.gammaImag = math.Sqrt(-s.point1.gammaReal - s.point1.gammaReal*s.point1.gammaReal)
	s.calcEndPoint()
	// s.seriesReact = s.point1.x - s.point0.x
	// s.seriesSuscep = -1.0 / s.seriesReact
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
func fitLC(lc float64, base []float64) (float64, []*matchParts, bool) {
	match := []*matchParts{ //ordered from the highest value to the lowest value
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
		&matchParts{},
	}
	var biggest float64
	for _, b := range base {
		biggest += b
	}
	if lc > biggest {
		return lc, match, false
	}
	smallest := base[len(base)-1]
	y := 0.0
	for i, item := range base {
		if lc > item {
			match[i].inPlay = true
			match[i].value = item
			lc -= item
			y += item
		}
	}
	if math.Abs(lc-smallest) < smallest/2 {
		y += smallest
	}
	return y, match, true
}

func (s *smith) calcFittedLC() (float64, float64) {
	var l, c float64
	for _, lItem := range s.matchL {
		if lItem.inPlay {
			l += lItem.value
		}
	}
	for _, cItem := range s.matchC {
		if cItem.inPlay {
			l += cItem.value
		}
	}
	return l, c
}

func (s *smith) addCCurrent(line string) string {
	for _, c := range s.matchC {
		line += fmt.Sprintf("%.2f,", cmplx.Abs(c.iThrough))
	}
	return line
}

func (s *smith) addLVoltage(line string) string {
	for _, l := range s.matchL {
		line += fmt.Sprintf("%.2f,", cmplx.Abs(l.vAcross))
	}
	return line
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

func (s *smith) calcYLoad() complex128 {
	r := s.point0.r * z0
	x := s.point0.x * z0
	// r := s.baseMaxSeries1.basePoint.r * z0
	// x := s.baseMaxSeries1.basePoint.x * z0
	//fmt.Println(r, x)
	g, b := getAdm(r, x)
	return complex(g, b)
}

func (s *smith) calcZLoad() complex128 {
	r := s.baseMaxSeries1.basePoint.r * z0
	x := s.baseMaxSeries1.basePoint.x * z0
	return complex(r, x)
}

func calcZfromY(y complex128) complex128 {
	g := real(y)
	b := imag(y)
	r, x := getImp(g, b)
	return complex(r, x)
}

func calcYfromZ(z complex128) complex128 {
	r := real(z)
	x := imag(z)
	g, b := getAdm(r, x)
	return complex(g, b)
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

func (s *smith) calcImpedance(f float64) {
	for i := range s.matchC {
		if s.matchC[i].inPlay {
			s.matchC[i].impedance = complex(0.0, (-1.0)/(2*math.Pi*s.matchC[i].value*f))
		}
	}
	for i := range s.matchL {
		if s.matchL[i].inPlay {
			s.matchL[i].impedance = complex(0.0, 2*math.Pi*s.matchL[i].value*f)
		}
	}
}

func (s *smith) sumLC() (float64, float64) {
	var l, c float64
	for _, match := range s.matchL {
		if match.inPlay {
			l += match.value
		}
	}
	for _, match := range s.matchC {
		if match.inPlay {
			c += match.value
		}
	}
	return l, c
}

//also returns the load current
func (s *smith) capCurrent(vParallel complex128) {
	for i := range s.matchC {
		if s.matchC[i].inPlay {
			s.matchC[i].iThrough = vParallel / s.matchC[i].impedance
		}
	}
}

func (s *smith) indVoltage(iSeries complex128) {
	for i := range s.matchL {
		if s.matchL[i].inPlay {
			s.matchL[i].vAcross = iSeries * s.matchL[i].impedance
		}
	}
}
