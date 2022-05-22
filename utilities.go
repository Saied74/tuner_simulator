package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func makeSmith() *smith {
	return &smith{
		outputFile:       "data.csv",
		point0:           &smithPoint{},
		point1:           &smithPoint{},
		gainTol:          0.15,
		phaseTol:         8.0,
		threshold:        1.2,
		iteration:        2,
		normalize:        "not",
		baseMaxSeries1:   &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		baseMaxParallel1: &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		baseMinSeries1:   &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		baseMinParallel1: &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		baseMaxSeries2:   &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		baseMaxParallel2: &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		baseMinSeries2:   &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		baseMinParallel2: &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},

		tolMaxSeries1:   &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		tolMaxParallel1: &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		tolMinSeries1:   &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		tolMinParallel1: &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		tolMaxSeries2:   &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		tolMaxParallel2: &extreme{parallelReact: -100000.0, seriesReact: -100000.0, basePoint: &smithPoint{}},
		tolMinSeries2:   &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
		tolMinParallel2: &extreme{parallelReact: 100000.0, seriesReact: 100000.0, basePoint: &smithPoint{}},
	}
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

func (s *smith) bruteIt() (float64, float64) {
	var r, x float64
	var inc = 0.1
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
		c := -1.0 / (2.0 * math.Pi * freqs[freq] * s.parallelReact * 50.0)
		l := (s.seriesReact * 50.0) / (2.0 * math.Pi * freqs[freq])
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

func (s *smith) calcWriteLCValues(f *os.File) error {
	var line, cFix, lFix = "", "", ""
	var l, c float64
	for _, val := range freqList {
		f, ok := freqs[val]
		if !ok {
			return fmt.Errorf("bad index into freqList")
		}
		if math.Abs(s.parallelReact) > math.SmallestNonzeroFloat64 {
			c = (-1.0) / (2.0 * math.Pi * f * s.parallelReact * 50)
		}
		if s.seriesReact > math.SmallestNonzeroFloat64 {
			l = (s.seriesReact * 50) / (2 * math.Pi * f)
		}
		if strings.HasPrefix(s.normalize, "norm") {
			c, cFix = normalize(c)
			cFix += "F"
			l, lFix = normalize(l)
			lFix += "H"
			line += fmt.Sprintf("%.2f %s,%.2f %s,", c, cFix, l, lFix)
		}
		if !strings.HasPrefix(s.normalize, "norm") {
			line += fmt.Sprintf("%e,%e,", c, l)
		}

	}
	_, err := f.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}

func normalize(lc float64) (float64, string) {
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
