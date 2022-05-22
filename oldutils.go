package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func (s *smith) bruteTwo() (float64, float64) {
	var r, x float64
	switch s.region {
	case 1:
		b := s.baseMaxSeries1.basePoint.b // + s.parallelSuscep
		g := s.baseMaxSeries1.basePoint.g
		r, x = getImp(g, b)
		x += s.seriesReact
		return r, x
	case 2:
		x = s.baseMaxSeries1.basePoint.x // + s.seriesReact
		r = s.baseMaxSeries1.basePoint.r
		g, b := getAdm(r, x)
		b += s.parallelSuscep
		r, x = getImp(g, b)
		return r, x
	}
	return r, x //the program never reaches this point except when it fails
}

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

func (s *smith) calcGammaTheta(r, x float64) {
	denom := (r+1)*(r+1) + x*x
	firstTerm := r*r + x*x - 1.0
	s.point0.gammaReal = firstTerm / denom
	s.point0.gammaImag = 4 * x * x / denom
	s.gamma = math.Sqrt(s.point0.gammaReal*s.point0.gammaReal +
		s.point0.gammaImag*s.point0.gammaImag)
	s.theta = math.Atan(s.point0.gammaImag / s.point0.gammaReal)
}

//writes the header for the tolerance study.  It generates the text from
//the tolerance values.  G stands for Gamma, T stands for Theta
func writeToleranceHeader(f *os.File) error {

	for _, t := range tolerance {
		tt := int(t * 100)
		item := fmt.Sprintf("Region G %d,Region T %d,", tt, tt)
		_, err := f.WriteString(item)
		if err != nil {
			return err
		}
	}
	return nil
}

//writes header for the maximum and minimum values of parallel and series impedances
func writeDistanceHeader(f *os.File) error {
	_, err := f.WriteString("Struct,Varied,SWR,Gamma,Theta,Reigon,r,x,Parallel X,Series X")
	if err != nil {
		return err
	}
	return nil
}
