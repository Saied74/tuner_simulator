package main

//The usage of this simulator is described in my blog :
//https://ad2cc.blogspot.com/2022/07/500-watt-antenna-tuner-part-3-simulator.html
//I have also copied the contents of the blog page in the README.md file
//The strucuture of the program is very simple.  In the file ui.go, I generate
//the user command line interface using my cli package:
//https://github.com/Saied74/cli/blob/master/cli.go
//Function main invkes the cli package and recieves the user commands through a
//channel.  Then through a switch statement, it selects what command to run.
//The house keeping commands like quit and file name changes are self explantery.
//The three calculations (simpleLC, fitLC, and calcVI use the analytical model
//that is also available in the blog.
import (
	"fmt"
	"log"
	"math/cmplx"
	"os"
	"strconv"
	//	"strings"

	"github.com/Saied74/cli"
)

const (
	bigImpedance = complex(0, 100000)
	zSource      = complex(50, 0)
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

type matchParts struct {
	inPlay     bool
	value      float64
	resistance float64
	impedance  complex128
	vAcross    complex128
	iThrough   complex128
	vToGround  complex128
}

type smith struct {
	outputFile     string
	minMaxFile     string
	s              float64 //swr
	gamma          float64 //magnitude of the reflection coefficient
	gammaTemp      float64
	theta          float64 //phase of the reflection coefficient
	thetaTemp      float64
	point0         *smithPoint
	point1         *smithPoint
	region         int
	parallelReact  float64
	parallelSuscep float64
	seriesReact    float64
	seriesSuscep   float64
	freqs          []float64
	// tolerance      []sensitivity
	which          string //use gain error, phase error, or both
	iteration      int
	threshold      float64
	gainTol        float64
	phaseTol       float64
	normalize      string //convert floating point numbers to pF, nH, etc.
	options        string
	minMax         map[string]*lcMinMax
	baseMaxSeries1 *extreme
	matchC         []*matchParts //ordered from the highest value to the lowest value
	matchL         []*matchParts //ordered from the highest value to the lowest value
	capQ           float64
	indQ           float64
	vSource        float64
	power          float64
}

type lcMinMax struct {
	freq float64
	minC float64
	maxC float64
	minL float64
	maxL float64
}

type maxVI map[string][]*matchParts

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

// Important Note:  Capacitance values are listed from the largest to the smallest

var baseCap = []float64{6000.0e-12, 3000.0e-12, 1410.0e-12, 690.0e-12, 330.0e-12, 180.0e-12,
	86.0e-12, 43.0e-12, 22.0e-12, 12.0e-12}

/*
var baseCap = []float64{6000.0e-12, 3000.0e-12, 1500.0e-12, 750.0e-12, 375.0e-12, 187.0e-12,
	94.0e-12, 47.0e-12, 22.0e-12, 12.0e-12}
*/

// Important note:  Inductance values are listed from the largest to the smallest
var baseInductor = []float64{12800.0e-9, 6400.0e-9, 3200.0e-9, 1600.0e-9, 800.0e-9, 400.0e-9,
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
		case "simpleLC":
			//open the full data set and min/max files
			f1, f2 := s.openTwoFiles()
			defer f1.Close()
			defer f2.Close()
			s.writeLCandFitHeaders(f1, f2)
			for _, w := range swrList {
				for i := 0; i < 360; i++ {
					s := s.resetSmith(home)
					s.s = w
					s.gamma = (s.s - 1.0) / (s.s + 1.0)
					s.theta = float64(i)
					s.trueCalc()
					s.writeImpedance(f1)
					for _, freqVal := range freqList {
						var l, c float64
						freq, ok := freqs[freqVal]
						if !ok {
							log.Fatal(fmt.Errorf("bad index into freqList"))
						}
						l, c = s.calcLCValues(freq) //write once per swr/theta line
						s.writeSimpleLCValues(l, c, f1)
						s.calcMinMax(l, c, freq, freqVal) //write once at the end
					}
					_, err := f1.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			s.writeSimpleMMValues(f2)

		case "fitLC":
			f1, f2 := s.openTwoFiles()
			defer f1.Close()
			defer f2.Close()
			s.writeLCandFitHeaders(f1, f2)
			for _, w := range swrList {
				for i := 0; i < 360; i++ {
					s := s.resetSmith(home)
					s.s = w
					s.gamma = (s.s - 1.0) / (s.s + 1.0)
					s.theta = float64(i)
					s.trueCalc()
					s.writeImpedance(f1)
					for _, freqVal := range freqList {
						var l, c float64
						var matchC, matchL []*matchParts
						freq, ok := freqs[freqVal]
						if !ok {
							log.Fatal(fmt.Errorf("bad index into freqList"))
						}
						l, c = s.calcLCValues(freq) //write once per swr/theta line

						c, matchC, ok = fitLC(c, baseCap) //not ok indicates could not fit
						if !ok {
							log.Fatal("Capacitor was too big")
						}
						l, matchL, ok = fitLC(l, baseInductor)
						if !ok {
							log.Fatal("Inductor was too big")
						}
						s.matchC = matchC
						s.matchL = matchL
						s.writeSimpleLCValues(l, c, f1)
						s.calcMinMax(l, c, freq, freqVal) //write once at the end
					}
					_, err := f1.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			s.writeSimpleMMValues(f2)

		case "calcVI":
			f1, f2 := s.openTwoFiles()
			defer f1.Close()
			defer f2.Close()
			s.writeVIHeaders(f1) //note writeVIHeaders is different than writeVIHeader
			//minmax VI header is written at the same time as data.
			m := makeMaxVI()
			for _, w := range swrList {
				for i := 0; i < 360; i++ {
					s := s.resetSmith(home)
					s.s = w
					s.gamma = (s.s - 1.0) / (s.s + 1.0)
					s.theta = float64(i)
					s.trueCalc()
					s.writeImpedance(f1)
					for _, freqVal := range freqList {
						var l, c float64
						var matchC, matchL []*matchParts
						freq, ok := freqs[freqVal]
						if !ok {
							log.Fatal(fmt.Errorf("bad index into freqList"))
						}
						l, c = s.calcLCValues(freq) //write once per swr/theta line

						c, matchC, ok = fitLC(c, baseCap) //not ok indicates could not fit
						if !ok {
							log.Fatal("Capacitor was too big")
						}
						l, matchL, ok = fitLC(l, baseInductor)
						if !ok {
							log.Fatal("Inductor was too big")
						}
						s.matchC = matchC
						s.matchL = matchL

						var line string
						var parallelY, seriesZ complex128
						zSource := complex(z0, 0)
						vSource := complex(s.vSource, 0)
						s.copyExt(s.baseMaxSeries1)
						seriesZ, parallelY = s.calcImpedance(freq) //matching circuit components

						switch s.region {
						case 1:
							line = ""
							zLoad := s.calcZLoad()
							yLoad := 1.0 / zLoad      // complex(s.point0.r * z0, s.point0.x * z0)
							yTwo := yLoad + parallelY //add load and parallalel capacitor admittances
							zTwo := 1.0 / yTwo
							zThree := zTwo + seriesZ
							zFour := zThree + zSource
							iSeries := vSource / zFour
							vParallel := iSeries * zTwo
							s.capCurrent(vParallel) //calculate capacitor currents
							s.indVoltage(iSeries)   //calculate the voltage across each inductor
							line += fmt.Sprintf("%.2f,", cmplx.Abs(vParallel))
							line = s.addCCurrent(line)
							line += fmt.Sprintf("%f,", cmplx.Abs(iSeries))
							line = s.addLVoltage(line)
							m.calcMaxEngaged(s, vParallel, iSeries, freqVal)
						case 2:
							line = ""
							zLoad := s.calcZLoad()  //calculate load impedance
							zTwo := zLoad + seriesZ //looking into load in series with the tuning inductor
							yTwo := 1.0 / zTwo
							yThree := yTwo + parallelY //looking into zTwo in parallel with tuning capacitor
							zThree := 1.0 / yThree
							zFour := zSource + zThree  //total load on the generator
							iSource := vSource / zFour //total current provided by the generator
							vParallel := iSource * zThree
							s.capCurrent(vParallel)
							iSeries := vParallel * yTwo
							s.indVoltage(iSeries)
							line += fmt.Sprintf("%.2f,", s.vSource)
							line = s.addCCurrent(line)
							line += fmt.Sprintf("%f,", cmplx.Abs(iSeries))
							line = s.addLVoltage(line)
							m.calcMaxEngaged(s, vParallel /*complex(s.vSource, 0)*/, iSeries, freqVal)
						case 3:
							line = ""
							zLoad := s.calcZLoad()
							zThree := zLoad + seriesZ
							zFour := zSource + zThree
							iSeries := vSource / zFour
							s.indVoltage(iSeries)
							line += fmt.Sprintf("%.2f,", s.vSource)
							line = s.addCCurrent(line)
							line += fmt.Sprintf("%f,", cmplx.Abs(iSeries))
							line = s.addLVoltage(line)
							m.calcMaxEngaged(s, complex(0, 0), iSeries, freqVal)
						}
						_, err := f1.WriteString(line)
						if err != nil {
							log.Fatal(err)
						}

						s.calcMinMax(l, c, freq, freqVal) //write once at the end
					}
					_, err := f1.WriteString("\n")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			m.writeMaxVI(f2)

		case "fileName":
			s.outputFile = home + item.Value
		case "minMaxFile":
			s.minMaxFile = home + item.Value
		case "vSource":
			y, _ := strconv.Atoi(item.Value)
			s.vSource = float64(y)
		case "capQ":
			y, _ := strconv.Atoi(item.Value)
			s.capQ = float64(y)
		case "indQ":
			y, _ := strconv.Atoi(item.Value)
			s.indQ = float64(y)
		default:
			log.Fatal("Bad parameter passed", item.Name, item.Value)
		}
	}
}
