package main

import (
	"log"
	"math"
	"os"

	"github.com/Saied74/cli"
)

type sensitivity struct {
	region            int
	parallelReactance float64
	seriesReactance   float64
}
type smithPoint struct {
	gammaReal float64
	gammaImag float64
	r         float64
	x         float64
	g         float64
	b         float64
}

type smith struct {
	outputFile     string
	s              float64
	gamma          float64
	gammaTemp      float64
	theta          float64
	thetaTemp      float64
	point0         *smithPoint
	point1         *smithPoint
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
	s := smith{
		outputFile: "data.csv",
		point0:     &smithPoint{},
		point1:     &smithPoint{},
	}
	c := cli.Command(&uiItems)
	for {
		item := <-c
		switch item.Name {
		case "Quit":
			os.Exit(1)
		case "noError":
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0555)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = writeImpedanceHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			for _, w := range swr {
				s.s = w
				s.gamma = (s.s - 1.0) / (s.s + 1.0)
				for i := 0; i < 360; i++ {
					s.theta = float64(i)
					s.trueCalc()
					err = s.writeImpedance(f)
					if err != nil {
						log.Fatal(err)
					}
					_, err = f.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case "tolerance":
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0555)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = writeImpedanceHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			err = writeToleranceHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			for _, w := range swr {
				s.s = w
				s.gamma = (s.s - 1.0) / (s.s + 1.0)
				for i := 0; i < 360; i++ {
					s.theta = float64(i)
					s.trueCalc()
					s.gammaTemp = s.gamma
					s.thetaTemp = s.theta
					s.calcTolerance()
					err = s.writeImpedance(f)
					if err != nil {
						log.Fatal(err)
					}
					err := s.writeTolerance(f)
					if err != nil {
						log.Fatal(err)
					}
					_, err = f.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}
					s.gamma = s.gammaTemp
					s.theta = s.thetaTemp
				}
			}
		case "distance":
		case "fileName":
			s.outputFile = item.Value
		default:
			continue
		}
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

//	s.writeImpedance(f)
// 			s.gammaTemp = s.gamma
// 			s.thetaTemp = s.theta
// 			s.calcTolerance()
// 			//	for i := 0; i < len(s.tolerance); i++ {
// 			s.writeImpedance(f)
// 			err := s.writeTolerance(f)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			_, err = f.WriteString("\n")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			//	}
// 			s.gamma = s.gammaTemp
// 			s.theta = s.thetaTemp
// 			// s.calcFreqs()
// 			// s.writeFreqs(f)
//
// 		}
// 	}
// }

func (s *smith) locate() {
	theta := 2.0 * (s.theta / 360.0) * math.Pi
	//fmt.Println(theta)
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

func (s *smith) calcFreqs() {
	s.freqs = []float64{}
	for _, freq := range freqList {
		c := -1.0 / (2.0 * math.Pi * freqs[freq] * s.parallelReact * 50.0)
		l := (s.seriesReact * 50.0) / (2.0 * math.Pi * freqs[freq])
		s.freqs = append(s.freqs, c, l)
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
