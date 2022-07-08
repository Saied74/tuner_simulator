package main

import (
	"fmt"
	"math"
	"os"
)

//this file is probably as useless as the useless.go file.

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

func (s *smith) calcGammaTheta(r, x float64) {
	denom := (r+1)*(r+1) + x*x
	firstTerm := r*r + x*x - 1.0
	s.point0.gammaReal = firstTerm / denom
	s.point0.gammaImag = 4 * x * x / denom
	s.gamma = math.Sqrt(s.point0.gammaReal*s.point0.gammaReal +
		s.point0.gammaImag*s.point0.gammaImag)
	s.theta = math.Atan(s.point0.gammaImag / s.point0.gammaReal)
}

//writes header for the maximum and minimum values of parallel and series impedances
func writeDistanceHeader(f *os.File) error {
	_, err := f.WriteString("Struct,Varied,SWR,Gamma,Theta,Reigon,r,x,Parallel X,Series X")
	if err != nil {
		return err
	}
	return nil
}

func (s *smith) writeDistance(f *os.File) error {
	line := fmt.Sprintf("baseMaxSeries1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMaxSeries1.s, s.baseMaxSeries1.gamma, s.baseMaxSeries1.theta,
		s.baseMaxSeries1.region, s.baseMaxSeries1.basePoint.r, s.baseMaxSeries1.basePoint.x,
		s.baseMaxSeries1.parallelReact, s.baseMaxSeries1.seriesReact)
	_, err := f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMinSeries1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMinSeries1.s, s.baseMinSeries1.gamma, s.baseMinSeries1.theta,
		s.baseMinSeries1.region, s.baseMinSeries1.basePoint.r, s.baseMinSeries1.basePoint.x,
		s.baseMinSeries1.parallelReact, s.baseMinSeries1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMaxParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMaxParallel1.s, s.baseMaxParallel1.gamma, s.baseMaxParallel1.theta,
		s.baseMaxParallel1.region, s.baseMaxParallel1.basePoint.r, s.baseMaxParallel1.basePoint.x,
		s.baseMaxParallel1.parallelReact, s.baseMaxParallel1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMinParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMinParallel1.s, s.baseMinParallel1.gamma, s.baseMinParallel1.theta,
		s.baseMinParallel1.region, s.baseMinParallel1.basePoint.r, s.baseMinParallel1.basePoint.x,
		s.baseMinParallel1.parallelReact, s.baseMinParallel1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMaxSeries2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMaxSeries2.s, s.baseMaxSeries2.gamma, s.baseMaxSeries2.theta,
		s.baseMaxSeries2.region, s.baseMaxSeries2.basePoint.r, s.baseMaxSeries2.basePoint.x,
		s.baseMaxSeries2.parallelReact, s.baseMaxSeries2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMinSeries2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMinSeries2.s, s.baseMinSeries2.gamma, s.baseMinSeries2.theta,
		s.baseMinSeries2.region, s.baseMinSeries2.basePoint.r, s.baseMinSeries2.basePoint.x,
		s.baseMinSeries2.parallelReact, s.baseMinSeries2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMaxParallel2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMaxParallel2.s, s.baseMaxParallel2.gamma, s.baseMaxParallel2.theta,
		s.baseMaxParallel2.region, s.baseMaxParallel2.basePoint.r, s.baseMaxParallel2.basePoint.x,
		s.baseMaxParallel2.parallelReact, s.baseMaxParallel2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("baseMinParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.baseMinParallel2.s, s.baseMinParallel2.gamma, s.baseMinParallel2.theta,
		s.baseMinParallel2.region, s.baseMinParallel2.basePoint.r, s.baseMinParallel2.basePoint.x,
		s.baseMinParallel2.parallelReact, s.baseMinParallel2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}

	line = fmt.Sprintf("tolMaxSeries1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMaxSeries1.s, s.tolMaxSeries1.gamma, s.tolMaxSeries1.theta,
		s.tolMaxSeries1.region, s.tolMaxSeries1.basePoint.r, s.tolMaxSeries1.basePoint.x,
		s.tolMaxSeries1.parallelReact, s.tolMaxSeries1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMinSeries1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMinSeries1.s, s.tolMinSeries1.gamma, s.tolMinSeries1.theta,
		s.tolMinSeries1.region, s.tolMinSeries1.basePoint.r, s.tolMinSeries1.basePoint.x,
		s.tolMinSeries1.parallelReact, s.tolMinSeries1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMaxParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMaxParallel1.s, s.tolMaxParallel1.gamma, s.tolMaxParallel1.theta,
		s.tolMaxParallel1.region, s.tolMaxParallel1.basePoint.r, s.tolMaxParallel1.basePoint.x,
		s.tolMaxParallel1.parallelReact, s.tolMaxParallel1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMinParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMinParallel1.s, s.tolMinParallel1.gamma, s.tolMinParallel1.theta,
		s.tolMinParallel1.region, s.tolMinParallel1.basePoint.r, s.tolMinParallel1.basePoint.x,
		s.tolMinParallel1.parallelReact, s.tolMinParallel1.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMaxSeries2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMaxSeries2.s, s.tolMaxSeries2.gamma, s.tolMaxSeries2.theta,
		s.tolMaxSeries2.region, s.tolMaxSeries2.basePoint.r, s.tolMaxSeries2.basePoint.x,
		s.tolMaxSeries2.parallelReact, s.tolMaxSeries2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMinSeries2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMinSeries2.s, s.tolMinSeries2.gamma, s.tolMinSeries2.theta,
		s.tolMinSeries2.region, s.tolMinSeries2.basePoint.r, s.tolMinSeries2.basePoint.x,
		s.tolMinSeries2.parallelReact, s.tolMinSeries2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMaxParallel2,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMaxParallel2.s, s.tolMaxParallel2.gamma, s.tolMaxParallel2.theta,
		s.tolMaxParallel2.region, s.tolMaxParallel2.basePoint.r, s.tolMaxParallel2.basePoint.x,
		s.tolMaxParallel2.parallelReact, s.tolMaxParallel2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}
	line = fmt.Sprintf("tolMinParallel1,%s,%.1f,%.3f,%0.0f,%d,%0.2f,%0.2f,%0.2f,%0.2f\n",
		s.which, s.tolMinParallel2.s, s.tolMinParallel2.gamma, s.tolMinParallel2.theta,
		s.tolMinParallel2.region, s.tolMinParallel2.basePoint.r, s.tolMinParallel2.basePoint.x,
		s.tolMinParallel2.parallelReact, s.tolMinParallel2.seriesReact)
	_, err = f.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}
