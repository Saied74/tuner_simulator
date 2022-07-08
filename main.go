package main

import (
	"fmt"
	"log"
	"math"
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
	minMaxFile       string
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
	options          string
	minMax           map[string]*lcMinMax
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

type lcMinMax struct {
	freq float64
	minC float64
	maxC float64
	minL float64
	maxL float64
}

var swrList = []float64{1.5, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

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
	"6m low C", "6m low L", "6m high C", "6m high L",
}

var freqList = []string{"160m low", "160m high", "80m low", "80m high",
	"40m low", "40m high", "20m low", "20m high", "17m low", "17m high",
	"15m low", "15m high", "12m low", "12m high", "10m low", "10m high",
	"6m low", "6m high"}

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
	"6m low":    50.0e6,
	"6m high":   54.0e6,
}

var tolerance = []float64{0.01, -0.01, 0.02, -0.02, 0.05, -0.05, 0.10, -0.10,
	0.15, -0.15, 0.20, -0.20, 0.25, -0.25}

var baseCap = []float64{3000.0e-12, 1500.0e-12, 750.0e-12, 360.0e-12, 180.0e-12,
	91.0e-12, 43.0e-12, 22.0e-12, 11.0e-12}

var baseInductor = []float64{6400.0e-9, 3200.0e-9, 1600.0e-9, 800.0e-9, 400.0e-9,
	200.0e-9, 100.0e-9, 50.0e-9, 25.0e-9}

func main() {
	home := os.Getenv("HOME")
	home += "/Documents/hamradio/Antennas/tuner/Simulation_output/"
	s := makeSmith(home)
	c := cli.Command(&uiItems)
	for {
		item := <-c
		switch item.Name {
		case "Quit":
			os.Exit(1)
		case "noError":
			s := s.resetSmith(home)
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
			for _, w := range swrList {
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
			// TODO: If the tolerance option to be used, it needs to be fixed.
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
			for _, w := range swrList {
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
		case "bruteForce":

			//s := s.resetSmith()
			f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = writeImpedanceHeader(f)
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
			f1, err := os.OpenFile(s.minMaxFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer f1.Close()
			err = writeMMHeader(f1)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f1.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
			for _, w := range swrList {
				var swr, r, x float64
				for i := 0; i < 360; i++ {

					s := s.resetSmith(home)
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
					for _, val := range freqList {
						var l, c, l0, c0 float64
						var ok bool
						freq, ok := freqs[val]
						if !ok {
							log.Fatal(fmt.Errorf("bad index into freqList"))
						}
						//this is a series of fall through conditions.
						//first, the true value of L & C matching elements are calculated
						//this is for all cases below
						l, c = s.calcLCValues(freq)
						//If called for, the min/max values of true L&C calculated
						//and written to file (at the end)
						//no other conditions will be met after this.
						if s.stepMMLC() {
							s.calcMinMax(l, c, freq, val)
						}
						//if called for, true LC values are approximated
						//this also sets up the condition for all the cases below
						if s.stepFitLC() {
							l0, c0 = l, c
							c, ok = fitLC(c0, baseCap)
							if !ok {
								c = 495.0
							}
							l, ok = fitLC(l0, baseInductor)
							if !ok {
								l = 495.0
							}
						}
						//if called for, the min/max values of the approximated LC are
						//calculated.
						//no other conditions will be met after this.
						if s.stepMMFitLC() {
							s.calcMinMax(l, c, freq, val)
						}
						//if called for, the difference between the true and approximate
						//LC values are calculated.
						//This is also a set up for the case that follows.
						if s.stepDelFitNotFit() {
							c = math.Abs(c0 - c)
							l = math.Abs(l0 - l)
						}
						if s.stepDelMMFitNotFit() {
							s.calcMinMax(l, c, freq, val)
						}

						err := s.writeLCValues(l, c, f)
						if err != nil {
							log.Fatal(err)
						}
					}
					_, err = f.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}

					s.gamma = s.gammaTemp
					s.theta = s.thetaTemp
				}
			}
			//condition for writing min/max values
			if s.stepMMLC() || s.stepMMFitLC() || s.stepDelMMFitNotFit() {
				err = s.writeMMValues(f1)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "fileName":
			s.outputFile = home + item.Value
		case "minMaxFile":
			s.minMaxFile = home + item.Value
		case "gainTol":
			x := strings.TrimSuffix(item.Value, "%")
			y, _ := strconv.Atoi(x)
			s.gainTol = float64(y) / 100
		case "phaseTol":
			y, _ := strconv.Atoi(item.Value)
			s.phaseTol = float64(y)
		case "which":
			s.which = item.Value
		case "threshold":
			thresh, _ := strconv.ParseFloat(item.Value, 64)
			s.threshold = thresh
		case "normalize":
			s.normalize = item.Value
		case "options":
			s.options = item.Value
		default:
			log.Fatal("Bad parameter passed", item.Name, item.Value)
		}
	}
}
