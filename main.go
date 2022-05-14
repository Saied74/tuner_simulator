package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

type sensitivity struct {
	region            int
	parallelReactance float64
	seriesReactance   float64
}

type smith struct {
	s              float64
	gamma          float64
	gammaTemp      float64
	theta          float64
	thetaTemp      float64
	gammaReal0     float64
	gammaImag0     float64
	r0             float64
	x0             float64
	g0             float64
	b0             float64
	gammaReal1     float64
	gammaImag1     float64
	r1             float64
	x1             float64
	g1             float64
	b1             float64
	region         int
	parallelReact  float64
	parallelSuscep float64
	seriesReact    float64
	seriesSuscep   float64
	freqs          []float64
	tolerance      []sensitivity
}

var swr = []float64{1.5, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

var angle = []float64{0, 20, 40, 60, 80, 100, 120, 180, 160, 180,
	200, 210, 240, 260, 280, 300, 320, 340, 360}

var impedance = []string{"swr", "theta", "r0", "x0", "r1", "x1", "region",
	"parallel", "series"}

var rcValues = []string{"160m low C", "160m low L", "160m high C", "160m high L",
	"80m low C", "80m low L", "80m high C", "80m high L",
	"40m low C", "40m low L", "40m high C", "40m high L",
	"20m low C", "20m low L", "20m high C", "20m high L",
	"17m low C", "17m low L", "17m high C", "17m high L",
	"15m low C", "15m low L", "15m high C", "15m high L",
	"12m low C", "12m low L", "12m high C", "12m high L",
	"10m low C", "12m low L", "10m high C", "12m high L",
}

var freqList = []string{"160m low", "160m high", "80m low", "80m high",
	"40m low", "40m high", "20m low", "20m high", "17m low", "17m high",
	"15m low", "15m high", "12m low", "12m high", "10m low", "10m hihg"}

var freqs = map[string]float64{
	"160m low":  1.8e6,
	"160m high": 2.0e6,
	"80m low":   3.5e6,
	"80m high":  4.0e6,
	"40m low":   7.0e6,
	"40m high":  7.3e6,
	"20m low":   14.0e6,
	"20m high":  14.350e6,
	"17m low":   18.068e6,
	"17m high":  18.168e6,
	"15m low":   21.0e6,
	"15m high":  21.450e6,
	"12m low":   24.89e6,
	"12m high":  24.990e6,
	"10m low":   28.0e6,
	"10m hihg":  29.7e6,
}

var tolerance = []float64{0.01, -0.01, 0.02, -0.02, 0.05, -0.05, 0.10, -0.10,
	0.15, -0.15, 0.20, -0.20, 0.25, -0.25}

func main() {
	f, err := os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = writeImpedanceHeader(f)
	if err != nil {
		log.Fatal(err)
	}
	err = writeToleranceHeader(f)
	// err = writeFreqHeader(f)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
	s := smith{}
	for _, w := range swr {
		s.s = w
		s.gamma = (s.s - 1.0) / (s.s + 1.0)
		// for _, a := range angle {
		for i := 0; i < 360; i++ {
			a := float64(i)
			s.theta = a
			s.locate()
			switch s.region {
			case 1:
				s.rotateRight()
			case 2:
				s.rotateLeft()
			default:
				log.Fatal("bad region")
			}
			//	s.writeImpedance(f)
			s.gammaTemp = s.gamma
			s.thetaTemp = s.theta
			s.calcTolerance()
			//	for i := 0; i < len(s.tolerance); i++ {
			s.writeImpedance(f)
			err := s.writeTolerance(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			//	}
			s.gamma = s.gammaTemp
			s.theta = s.thetaTemp
			// s.calcFreqs()
			// s.writeFreqs(f)

		}
	}
}

func (s *smith) locate() {
	theta := 2.0 * (s.theta / 360.0) * math.Pi
	//fmt.Println(theta)
	s.gammaReal0 = s.gamma * math.Cos(theta)
	s.gammaImag0 = s.gamma * math.Sin(theta)
	//	fmt.Println(s.gammaReal0, s.gammaImag0)
	gammaSq := s.gamma * s.gamma
	denom := 1 + gammaSq - 2*s.gammaReal0
	s.r0 = (1 - gammaSq) / denom
	s.x0 = (2 * s.gammaImag0) / denom
	zSq := s.r0*s.r0 + s.x0*s.x0
	s.g0 = s.r0 / zSq
	s.b0 = -s.x0 / zSq
	if s.r0 > 1.0 {
		s.region = 1
		return
	}
	if s.x0 > 0.0 && s.r0 < zSq {
		s.region = 1
		return
	}
	s.region = 2
	return
}

func (s *smith) rotateRight() {
	s.gammaReal1 = (1.0 - s.g0) / (1.0 + 3*s.g0)
	s.gammaImag1 = -math.Sqrt(s.gammaReal1 - s.gammaReal1*s.gammaReal1)
	s.calcEndPoint()
	s.parallelSuscep = s.b1 - s.b0
	s.parallelReact = -1.0 / s.parallelSuscep
	s.seriesReact = -s.x1
	s.seriesSuscep = -1.0 / s.seriesReact
}

func (s *smith) rotateLeft() {
	s.gammaReal1 = (s.r0 - 1.0) / (3*s.r0 + 1)
	s.gammaImag1 = math.Sqrt(-s.gammaReal1 - s.gammaReal1*s.gammaReal1)
	s.calcEndPoint()
	s.seriesReact = s.x1 - s.x0
	s.seriesSuscep = -1.0 / s.seriesReact
	s.parallelSuscep = -s.b1
	s.parallelReact = -1.0 / s.parallelSuscep
}

func (s *smith) calcEndPoint() {
	gammaSq := s.gammaReal1*s.gammaReal1 + s.gammaImag1*s.gammaImag1
	denom := 1 + gammaSq - 2*s.gammaReal1
	//	fmt.Println(s.gammaReal1, s.gammaImag1, gammaSq, denom)
	s.r1 = (1 - gammaSq) / denom
	s.x1 = 2 * s.gammaImag1 / denom
	zSq := s.r1*s.r1 + s.x1*s.x1
	s.g1 = s.r1 / zSq
	s.b1 = -s.x1 / zSq
}

func (s *smith) calcFreqs() {
	s.freqs = []float64{}
	for _, freq := range freqList {
		c := -1.0 / (2.0 * math.Pi * freqs[freq] * s.parallelReact * 50.0)
		l := (s.seriesReact * 50.0) / (2.0 * math.Pi * freqs[freq])
		s.freqs = append(s.freqs, c, l)
	}
}

func writeImpedanceHeader(f *os.File) error {
	for _, item := range impedance {
		_, err := f.WriteString(item + ",")
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFreqHeader(f *os.File) error {
	for _, item := range rcValues {
		_, err := f.WriteString(item + ",")
		if err != nil {
			return err
		}
	}
	return nil
}

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

func (s *smith) writeImpedance(f *os.File) {
	line := fmt.Sprintf("%.1f,%0.0f,%0.2f,%0.2f,%0.2f,%0.2f,%d,%0.2f,%0.2f,",
		s.s, s.theta, s.r0, s.x0, s.r1, s.x1, s.region, s.parallelReact, s.seriesReact)
	f.WriteString(line)
}

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

func (s *smith) writeTolerance(f *os.File) error {
	for _, item := range s.tolerance {
		_, err := f.WriteString(fmt.Sprintf("%d,", item.region))
		if err != nil {
			return err
		}
		// _, err = f.WriteString(fmt.Sprintf("%0.2f,", item.parallelReactance))
		// if err != nil {
		// 	return err
		// }
		// _, err = f.WriteString(fmt.Sprintf("%0.2f,", item.seriesReactance))
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
