package main

import (
	"fmt"
	"os"
)

//writes the header for the base case with no errors
func writeImpedanceHeader(f *os.File) error {
	//csv first line for the base case with no errors
	var impedance = []string{"swr", "theta", "r0", "x0", "r1", "x1", "region",
		"parallel", "series"}

	for _, item := range impedance {
		_, err := f.WriteString(item + ",")
		if err != nil {
			return err
		}
	}
	return nil
}

//header for when the actual value of L and C are calculated
//past use, may not have any future use
func writeLCHeader(f *os.File) error {
	for _, item := range lcValues {
		_, err := f.WriteString(item + ",")
		if err != nil {
			return err
		}
	}
	return nil
}

//writes the base case (no errors) data
func (s *smith) writeImpedance(f *os.File) error {
	line := fmt.Sprintf("%.1f,%0.0f,%0.2f,%0.2f,%0.2f,%0.2f,%d,%0.2f,%0.2f,",
		s.s, s.theta, s.point0.r, s.point0.x, s.point1.r, s.point1.x, s.region, s.parallelReact, s.seriesReact)
	_, err := f.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}

//write the results of the tolerance study (region number)
func (s *smith) writeTolerance(f *os.File) error {
	for _, item := range s.tolerance {
		_, err := f.WriteString(fmt.Sprintf("%d,", item.region))
		if err != nil {
			return err
		}
	}
	return nil
}

//writes the actual Ls and Cs based on freequency of bands
//past use, may not have any future use
func (s *smith) writeFreqs(f *os.File) {
	var m float64
	for i, freq := range s.freqs {
		m = 1e9
		if i%2 == 0 {
			m = 1e12
		}
		f.WriteString(fmt.Sprintf("%f,", freq*m))
	}
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
