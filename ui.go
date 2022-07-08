package main

import (
	"strconv"
	"strings"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"fileName", "minMaxFile", "options", "gainTol", "phaseTol",
		"noError", "normalize", "threshold", "which", "tolerance", "bruteForce"},
	ItemList: map[string]*cli.Item{
		"noError": &cli.Item{
			Name:      "noError",
			Prompt:    "Run the simulation with no errors and write a csv file",
			Response:  "Do I need this 1?",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"tolerance": &cli.Item{
			Name:      "tolerance",
			Prompt:    "Run tolerance study according to the tolerance list (needs to be fixed)",
			Response:  "Do I need this? 2",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"bruteForce": &cli.Item{
			Name:      "bruteForce",
			Prompt:    "Brute force through the error and see",
			Response:  "Do I need this 3?",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"fileName": &cli.Item{
			Name:      "fileName",
			Prompt:    "Change the csv file name",
			Response:  "Do I need this 4?",
			Value:     "data.csv",
			Validator: filenameValidator,
		},
		"minMaxFile": &cli.Item{
			Name:      "minMaxFile",
			Prompt:    "Change the minMax file name",
			Response:  "Do I need this 4?",
			Value:     "minMax.csv",
			Validator: filenameValidator,
		},
		"gainTol": &cli.Item{
			Name:      "gainTol",
			Prompt:    "Specify the max gain error in percentage",
			Response:  "Do I need this 4?",
			Value:     "15%",
			Validator: gainValidator,
		},
		"phaseTol": &cli.Item{
			Name:      "phaseTol",
			Prompt:    "Specify the max phase error in degrees",
			Response:  "Do I need this 4?",
			Value:     "8",
			Validator: phaseValidator,
		},
		"which": &cli.Item{
			Name:      "which",
			Prompt:    "Select theta or gamma (or both) for tolerance study",
			Response:  "Do I need this 4?",
			Value:     "",
			Validator: whichValidator,
		},

		"threshold": &cli.Item{
			Name:      "threshold",
			Prompt:    "Max SWR threshold to stop brute force iteration",
			Response:  "Do I need this 4?",
			Value:     "1.2",
			Validator: thresholdValidator,
		},
		"normalize": &cli.Item{
			Name:      "normalize",
			Prompt:    "Normalize (or not), LC values",
			Response:  "Do I need this 4?",
			Value:     "normalize",
			Validator: normalValidator,
		},
		/*
			In all cases below the results are written to file.
			In the min/max case, they are written to the minMaxFile
			LC: Calculate the series inductance and parallel capacitance
			MMLC: Calculate minimum and maximum LC values
			FitLC: Approximate LC values using baseCap and baseInductor values
			MMFitLC: Calculate minimum and maximum approximated LC values
			DelFitNotFit:  Calculate the difference between true and approximated LC
			lcValues
			DelMMFitNotFit: Calculate minumum and maximum difference between true and
			approximated values
		*/
		"options": &cli.Item{
			Name:      "options",
			Prompt:    "Select: LC, MMLC, FitLC, MMFitLC, DelFitNotFit, DelMMFitNotFit, VI",
			Response:  "Do I need this 4?",
			Value:     "MMFitLC",
			Validator: optionValidator,
		},
	},
	ActionLines: []string{"Enter the number of item you would like to runn or q to quit",
		"If the option has no parameters, press enter a second time"},
}

var filenameValidator = cli.ItemValidator(func(x string) bool {
	y := strings.Split(x, ".")
	if len(y) != 2 {
		return false
	}
	if y[1] != "csv" {
		return false
	}
	return true
})

var gainValidator = cli.ItemValidator(func(x string) bool {
	if !strings.HasSuffix(x, "%") {
		return false
	}
	x = strings.TrimSuffix(x, "%")
	y, err := strconv.Atoi(x)
	if err != nil {
		return false
	}
	if float64(y)/100 > 1.0 {
		return false
	}
	return true
})

var phaseValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.Atoi(x)
	if err != nil {
		return false
	}
	return true
})

var whichValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "theta":
		return true
	case "gamma":
		return true
	case "both":
		return true
	default:
		return false
	}
})

var thresholdValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return false
	}
	return true
})

var normalValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "normalize":
		return true
	case "normal":
		return true
	case "norm":
		return true
	case "no":
		return true
	case "not":
		return true
	}
	return false
})

var optionValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "LC":
		return true
	case "MMLC":
		return true
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
})
