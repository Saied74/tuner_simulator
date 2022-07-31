package main

import (
	"fmt"
	"math/cmplx"
	"os"
	"strings"
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
	_, err := f.WriteString("r,x,New Region,New Series,New Parallel,swr,")
	if err != nil {
		return err
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

//writes VI file header for all the current through C and voltage across L values
func writeVIHeader(f *os.File) error {
	var h string
	h += "swr,"
	h += "theta,"
	h += "series,"
	h += "parallel,"
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

//writes the header for the tolerance study.  It generates the text from
//the tolerance values.  G stands for Gamma, T stands for Theta
// func writeToleranceHeader(f *os.File) error {
//
// 	for _, t := range tolerance {
// 		tt := int(t * 100)
// 		item := fmt.Sprintf("Region G %d,Region T %d,", tt, tt)
// 		_, err := f.WriteString(item)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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
	// _, err = f.WriteString("\n")
	// if err != nil {
	// 	return err
	// }
	// _, err = f.WriteString("Min,")
	// if err != nil {
	// 	return err
	// }
	// line = ""
	// for _, val := range freqList {
	// 	mm := *s.minMax[val]
	// 	if strings.HasPrefix(s.normalize, "norm") {
	// 		c, cFix := normalizeLC(mm.minC)
	// 		cFix += "F"
	// 		l, lFix := normalizeLC(mm.minL)
	// 		lFix += "H"
	// 		line += fmt.Sprintf("%.2f %s,%.2f %s,", c, cFix, l, lFix)
	// 	}
	// 	if !strings.HasPrefix(s.normalize, "norm") {
	// 		line += fmt.Sprintf("%e,%e,", mm.minC, mm.minL)
	// 	}
	// }
	// _, err = f.WriteString(line)
	// if err != nil {
	// 	return err
	// }
	// _, err = f.WriteString("\n")
	// if err != nil {
	// 	return err
	// }
	return nil
}

// //write the results of the tolerance study (region number)
// func (s *smith) writeTolerance(f *os.File) error {
// 	for _, item := range s.tolerance {
// 		_, err := f.WriteString(fmt.Sprintf("%d,", item.region))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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

func (m maxVI) writeMaxVI(f1 *os.File) error {
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
		return err
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
			return err
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
		return err
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
			return err
		}
	}
	return nil
}
