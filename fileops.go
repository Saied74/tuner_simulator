package main

import (
	"fmt"
	"log"
	"math/cmplx"
	"os"
	"strings"
)

func (s *smith) openTwoFiles() (*os.File, *os.File) {
	f1, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.OpenFile(s.minMaxFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f1, f2
}

func (s *smith) writeLCandFitHeaders(f1, f2 *os.File) {
	err := writeImpedanceHeader(f1)
	if err != nil {
		log.Fatal(err)
	}
	err = writeLCHeader(f1)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f1.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *smith)writeVIHeaders(f1 *os.File) {
    err := writeImpedanceHeader(f1)
    if err != nil {
        log.Fatal(err)
    }

    err = writeVIHeader(f1)
    if err != nil {
        log.Fatal(err)
    }
}


func (s *smith) writeSimpleLCValues(l, c float64, f *os.File) {
	_, err := f.WriteString(fmt.Sprintf("%e,%e,", c, l))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *smith) writeSimpleMMValues(f *os.File) {
	var line = ",C,L\n"
	_, err := f.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range freqList {
		line = val + ","
		mm := *s.minMax[val]
		line += fmt.Sprintf("%e,%e,", mm.maxC, mm.maxL)
		line += "\n"
		_, err = f.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// writes the header for the base case with no errors
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
	_, err := f.WriteString("r,x,New Region,New Series,New Parallel,swr,")
	if err != nil {
		return err
	}
	return nil
}

// writes the base case (no errors) data
func (s *smith) writeImpedance(f *os.File) {
    var swr, r, x float64
	line := fmt.Sprintf("%.1f,%0.0f,%0.2f,%0.2f,%0.2f,%0.2f,%d,%0.2f,%0.2f,",
		s.s, s.theta, s.point0.r, s.point0.x, s.point1.r, s.point1.x, s.region, s.parallelReact, s.seriesReact)
	_, err := f.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
    swr = 1
	_, err = f.WriteString(fmt.Sprintf("%.3f,%.3f,%d,%.3f,%.3f,%.2f,",
					        r, x, s.region, s.seriesReact, s.parallelReact, swr))
	if err != nil {
    	log.Fatal(err)
	}

}

// writes VI file header for all the current through C and voltage across L values
func writeVIHeader(f *os.File) error {
	var h string
	for i, item := range lcValues {

		if i%2 == 0 {
			h += item + " VCap,"
			for _, cBase := range baseCap {
				c, fix := normalizeLC(cBase)
				h += fmt.Sprintf("%s %.2f%sF,", item, c, fix)
			}
		}
		if i%2 != 0 {
			h += item + " IInd,"
			for _, lBase := range baseInductor {
				l, fix := normalizeLC(lBase)
				h += fmt.Sprintf("%s %.2f%sH,", item, l, fix)
			}
		}
	}
	h = strings.TrimSuffix(h, ",")
	h += "\n"
	_, err := f.WriteString(h)
	return err
}


// header for when the actual value of L and C are calculated
// past use, may not have any future use
func writeLCHeader(f *os.File) error {
	for _, item := range lcValues {
		_, err := f.WriteString(item + ",")
		if err != nil {
			return err
		}
	}
	return nil
}

func writeMMHeader(f *os.File) error {
	_, err := f.WriteString("MinMax,")
	if err != nil {
		return err
	}
	err = writeLCHeader(f)
	if err != nil {
		return err
	}
	return nil
}

func (s *smith) writeMMValues(f *os.File) error {
	// _, err := f.WriteString("Max,")
	// if err != nil {
	// 	return err
	// }
	var line = ",C,L\n"
	_, err := f.WriteString(line)
	if err != nil {
		return err
	}

	for _, val := range freqList {
		line = val + ","
		mm := *s.minMax[val]
		if strings.HasPrefix(s.normalize, "norm") {
			c, cFix := normalizeLC(mm.maxC)
			cFix += "F"
			l, lFix := normalizeLC(mm.maxL)
			lFix += "H"
			line += fmt.Sprintf("%.2f %s,%.2f %s,", c, cFix, l, lFix)
		}
		if !strings.HasPrefix(s.normalize, "norm") {
			line += fmt.Sprintf("%e,%e,", mm.maxC, mm.maxL)
		}
		line += "\n"
		_, err = f.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// writes the actual Ls and Cs based on freequency of bands
// past use, may not have any future use
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

func (s *smith) writeLCValues(l, c float64, f *os.File) error {
	var line, cFix, lFix = "", "", ""
	if strings.HasPrefix(s.normalize, "norm") {
		c, cFix = normalizeLC(c)
		cFix += "F"
		l, lFix = normalizeLC(l)
		lFix += "H"
		line += fmt.Sprintf("%.2f %s,%.2f %s,", c, cFix, l, lFix)
	}
	if !strings.HasPrefix(s.normalize, "norm") {
		line += fmt.Sprintf("%e,%e,", c, l)
	}
	_, err := f.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}

func (m maxVI) writeMaxVI(f1 *os.File) {
	line := " ,"
	line += "Cap Voltage,"
	for _, c := range baseCap {
		cNorm, suffix := normalizeLC(c)
		line += fmt.Sprintf("%.2f %sF,", cNorm, suffix)
	}
	line = strings.TrimSuffix(line, ",")
	line += "\n"
	_, err := f1.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
	for _, freqVal := range freqList {
		key := freqVal + " C"
		line = key + ","
		for i, cMatch := range m[key] {
			if i == 0 {
				line += fmt.Sprintf("%.0f,", cmplx.Abs(cMatch.vAcross))
			} else {
				line += fmt.Sprintf("%.2f,", cmplx.Abs(cMatch.iThrough))
			}
		}
		line = strings.TrimSuffix(line, ",")
		line += "\n"
		_, err := f1.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}

	}
	_, err = f1.WriteString("\n\n\n\n")

	line = " ,"
	line += "Ind Current,"
	for _, l := range baseInductor {
		lNorm, suffix := normalizeLC(l)
		line += fmt.Sprintf("%.1f %sH,", lNorm, suffix)
	}
	line = strings.TrimSuffix(line, ",")
	line += "\n"
	_, err = f1.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
	for _, freqVal := range freqList {
		key := freqVal + " L"
		line = key + ","
		for i, lMatch := range m[key] {
			if i == 0 {
				line += fmt.Sprintf("%.1f,", cmplx.Abs(lMatch.iThrough))
			} else {
				line += fmt.Sprintf("%.0f,", cmplx.Abs(lMatch.vAcross))
			}
		}
		line = strings.TrimSuffix(line, ",")
		line += "\n"
		_, err := f1.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}
