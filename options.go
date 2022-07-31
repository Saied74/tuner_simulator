package main

//functions in this file test run time conditions and return true or false for the most part

//noError is a case that goes with LC and FitLC case as an option.
//noError over rides the bruteforce solution.
func (s *smith) stepNoError() bool {
	if s.which == "noError" {
		return true
	}
	return false
}

//should L and C file be openned and written but only for LC condition
func (s *smith) stepLCFile() bool {
	switch s.options {
	case "LC":
		return true
	}
	return false
}

//should L and C file be openned and written but only for FitLC condition
func (s *smith) stepFitLCFile() bool {
	switch s.options {
	case "FitLC":
		return true
	}
	return false
}

//should the VI file be opened and written
func (s *smith) stepVIFile() bool {
	if s.stepVI() {
		return true
	}
	return false
}

//should the min and max of L & C values file be openned and written
// func (s *smith) stepMMFile() bool {
// 	if s.stepMMLC() {
// 		return true
// 	}
// 	if s.stepMMFitLC() {
// 		return true
// 	}
// 	if s.stepDelMMFitNotFit() {
// 		return true
// 	}
// 	return false
// }

//If true, approximate LC values using baseCap and baseInductor values
//also set up the condition for running MMFitLC, DelMMFitNotFit, and DelMMFitNotFit
func (s *smith) stepFitLC() bool {
	switch s.options {
	case "FitLC":
		return true
	case "MMFitLC":
		return true
	case "DelFitNotFit":
		return true
	case "DelMMFitNotFit":
		return true
	case "VI":
		return true
	}
	return false
}

func (s *smith) stepVI() bool {
	switch s.options {
	case "VI":
		return true
	}
	return false
}

//If true run min/max true LC values.  No more processing will proceed.
// func (s *smith) stepMMLC() bool {
// 	switch s.options {
// 	case "MMLC":
// 		return true
// 	}
// 	return false
// }

//if true calculate the min/max of the approximated LC values
//no more processing will proceed past this point
func (s *smith) stepMMFitLC() bool {
	switch s.options {
	case "MMFitLC":
		return true
	}
	return false
}

//If true, claculate the difference between approximated and true LC lcValues
//also set up the condition for running stepDelMMFitNotFit
// func (s *smith) stepDelFitNotFit() bool {
// 	switch s.options {
// 	case "DelFitNotFit":
// 		return true
// 	case "DelMMFitNotFit":
// 		return true
// 	}
// 	return false
// }

//if true, calculate the difference between minimum and maximum of true and
//approximated LC values.
// func (s *smith) stepDelMMFitNotFit() bool {
// 	switch s.options {
// 	case "DelMMFitNotFit":
// 		return true
// 	}
// 	return false
// }
