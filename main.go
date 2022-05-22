package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

type extreme struct {
	s              float64
	gamma          float64
	theta          float64
	region         int
	basePoint      *smithPoint
	parallelReact  float64
	parallelSuscep float64
	seriesReact    float64
	seriesSuscep   float64
}

type smith struct {
	outputFile       string
	s                float64
	gamma            float64
	gammaTemp        float64
	theta            float64
	thetaTemp        float64
	point0           *smithPoint
	point1           *smithPoint
	region           int
	parallelReact    float64
	parallelSuscep   float64
	seriesReact      float64
	seriesSuscep     float64
	freqs            []float64
	tolerance        []sensitivity
	which            string
	iteration        int
	threshold        float64
	gainTol          float64
	phaseTol         float64
	normalize        string
	baseMaxSeries1   *extreme
	baseMaxParallel1 *extreme
	baseMinSeries1   *extreme
	baseMinParallel1 *extreme
	baseMaxSeries2   *extreme
	baseMaxParallel2 *extreme
	baseMinSeries2   *extreme
	baseMinParallel2 *extreme
	tolMaxSeries1    *extreme
	tolMaxParallel1  *extreme
	tolMinSeries1    *extreme
	tolMinParallel1  *extreme
	tolMaxSeries2    *extreme
	tolMaxParallel2  *extreme
	tolMinSeries2    *extreme
	tolMinParallel2  *extreme
}

var swr = []float64{1.5, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

// var angle = []float64{0, 20, 40, 60, 80, 100, 120, 180, 160, 180,
// 	200, 210, 240, 260, 280, 300, 320, 340, 360}

var lcValues = []string{"160m low C", "160m low L", "160m high C", "160m high L",
	"80m low C", "80m low L", "80m high C", "80m high L",
	"40m low C", "40m low L", "40m high C", "40m high L",
	"20m low C", "20m low L", "20m high C", "20m high L",
	"17m low C", "17m low L", "17m high C", "17m high L",
	"15m low C", "15m low L", "15m high C", "15m high L",
	"12m low C", "12m low L", "12m high C", "12m high L",
	"10m low C", "10m low L", "10m high C", "10m high L",
}

var freqList = []string{"160m low", "160m high", "80m low", "80m high",
	"40m low", "40m high", "20m low", "20m high", "17m low", "17m high",
	"15m low", "15m high", "12m low", "12m high", "10m low", "10m high"}

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
	"10m high":  29.7e6,
}

var tolerance = []float64{0.01, -0.01, 0.02, -0.02, 0.05, -0.05, 0.10, -0.10,
	0.15, -0.15, 0.20, -0.20, 0.25, -0.25}

func main() {
	s := makeSmith()
	c := cli.Command(&uiItems)
	for {
		item := <-c
		switch item.Name {
		case "Quit":
			os.Exit(1)
		case "noError":
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
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
		case "oneError":
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = writeImpedanceHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("swr")
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			for _, w := range swr {
				s.s = w
				gamma := (s.s - 1.0) / (s.s + 1.0)
				s.gamma = gamma
				if s.which == "gamma" {
					s.gamma += s.gainTol * gamma
				}
				for i := 0; i < 360; i++ {
					theta := float64(i)
					s.theta = theta
					if s.which == "theta" {
						s.theta += s.phaseTol
					}
					s.trueCalc()
					err = s.writeImpedance(f)
					if err != nil {
						log.Fatal(err)
					}
					swr := calcSWR(s.point1.r, s.point1.x)
					_, err = f.WriteString(fmt.Sprintf("%.2f", swr))
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
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
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
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			err = writeDistanceHeader(f)
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
					switch s.region {
					case 1:
						if s.seriesReact > s.baseMaxSeries1.seriesReact {
							s.copyExt(s.baseMaxSeries1)
						}
						if s.seriesReact < s.baseMinSeries1.seriesReact {
							s.copyExt(s.baseMinSeries1)
						}
						if s.parallelReact > s.baseMaxParallel1.parallelReact {
							s.copyExt(s.baseMaxParallel1)
						}
						if s.parallelReact < s.baseMinParallel1.parallelReact {
							s.copyExt(s.baseMinParallel1)
						}
					case 2:
						if s.seriesReact > s.baseMaxSeries2.seriesReact {
							s.copyExt(s.baseMaxSeries2)
						}
						if s.seriesReact < s.baseMinSeries2.seriesReact {
							s.copyExt(s.baseMinSeries2)
						}
						if s.parallelReact > s.baseMaxParallel2.parallelReact {
							s.copyExt(s.baseMaxParallel2)
						}
						if s.parallelReact < s.baseMinParallel2.parallelReact {
							s.copyExt(s.baseMinParallel2)
						}
					}
					s.gammaTemp = s.gamma
					s.thetaTemp = s.theta
					switch s.which {
					case "theta":
						s.theta += s.phaseTol
					case "gamma":
						s.gamma += s.gamma * s.gainTol
					}
					s.trueCalc()
					switch s.region {
					case 1:
						if s.seriesReact > s.tolMaxSeries1.seriesReact {
							s.copyExt(s.tolMaxSeries1)
						}
						if s.seriesReact < s.tolMinSeries1.seriesReact {
							s.copyExt(s.tolMinSeries1)
						}
						if s.parallelReact > s.tolMaxParallel1.parallelReact {
							s.copyExt(s.tolMaxParallel1)
						}
						if s.parallelReact < s.tolMinParallel1.parallelReact {
							s.copyExt(s.tolMinParallel1)
						}
					case 2:
						if s.seriesReact > s.tolMaxSeries2.seriesReact {
							s.copyExt(s.tolMaxSeries2)
						}
						if s.seriesReact < s.tolMinSeries2.seriesReact {
							s.copyExt(s.tolMinSeries2)
						}
						if s.parallelReact > s.tolMaxParallel2.parallelReact {
							s.copyExt(s.tolMaxParallel2)
						}
						if s.parallelReact < s.tolMinParallel2.parallelReact {
							s.copyExt(s.tolMinParallel2)
						}
					}
					s.gamma = s.gammaTemp
					s.theta = s.thetaTemp
				}
			}
			err = s.writeDistance(f)
			if err != nil {
				log.Fatal(err)
			}
		case "bruteForce":
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = writeImpedanceHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("r,x,New Region,New Series,New Parallel,swr,")
			if err != nil {
				log.Fatal(err)
			}
			err = writeLCHeader(f)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			for _, w := range swr {
				var swr, r, x float64
				for i := 0; i < 360; i++ {
					s = makeSmith()
					s.s = w
					s.gamma = (s.s - 1.0) / (s.s + 1.0)
					s.theta = float64(i)
					s.trueCalc()
					err = s.writeImpedance(f)
					if err != nil {
						log.Fatal(err)
					}
					s.gammaTemp = s.gamma
					s.thetaTemp = s.theta
					s.copyExt(s.baseMaxSeries1)
					switch s.which {
					case "gamma":
						s.gamma += s.gainTol * s.gamma
					case "theta":
						s.theta += s.phaseTol
					case "both":
						s.gamma += s.gainTol * s.gamma
						s.theta += s.phaseTol
					}
					s.trueCalc()
					r, x = s.bruteIt()
					swr = calcSWR(r, x)

					_, err = f.WriteString(fmt.Sprintf("%.3f,%.3f,%d,%.3f,%.3f,%.2f,",
						r, x, s.region, s.seriesReact, s.parallelReact, swr))
					if err != nil {
						log.Fatal(err)
					}
					err = s.calcWriteLCValues(f)
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
		case "fileName":
			s.outputFile = item.Value
		case "gainTol":
			x := strings.TrimSuffix(item.Value, "%")
			y, _ := strconv.Atoi(x)
			s.gainTol = float64(y) / 100
		case "phaseTol":
			y, _ := strconv.Atoi(item.Value)
			s.phaseTol = float64(y)
		case "which":
			s.which = item.Value
		case "iterations":
			iter, _ := strconv.Atoi(item.Value)
			s.iteration = iter
			fmt.Println("Iteration: ", s.iteration)
		case "threshold":
			thresh, _ := strconv.ParseFloat(item.Value, 64)
			s.threshold = thresh
		case "normalize":
			s.normalize = item.Value
		default:
			log.Fatal("Bad parameter passed", item.Name, item.Value)
		}
	}
}
