package main

import (
	"strconv"
	"strings"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"fileName", "minMaxFile", "vSource", "capQ", "indQ", "simpleLC", "fitLC", "calcVI"},
	ItemList: map[string]*cli.Item{
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
		"vSource": &cli.Item{
			Name:      "vSource",
			Prompt:    "Source voltage in volts",
			Response:  "useless",
			Value:     "223",
			Validator: voltageValidator,
		},
		"capQ": &cli.Item{
			Name:      "capQ",
			Prompt:    "capacitor Q",
			Response:  "useless",
			Value:     "1000",
			Validator: qValidator,
		},
		"indQ": &cli.Item{
			Name:      "indQ",
			Prompt:    "inductor Q",
			Response:  "useless",
			Value:     "100",
			Validator: qValidator,
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
		"simpleLC": &cli.Item{
			Name:      "simpleLC",
			Prompt:    "Calculate simple LC and LC min max calculation (enter \"sim\" to run",
			Response:  "Do I need this 4?",
			Value:     "",
			Validator: simpleValidator,
		},
		"fitLC": &cli.Item{
			Name:      "fitLC",
			Prompt:    "Fit LC to standard values and calculate max values (enter \"fit\" to run)",
			Response:  "Do I need this 4?",
			Value:     "",
			Validator: fitValidator,
		},
		"calcVI": &cli.Item{
			Name:      "calcVI",
			Prompt:    "calculate current and voltage across/through L & C and maximum values (enter \"calc\" to run)",
			Response:  "Do I need this 4?",
			Value:     "",
			Validator: calcValidator,
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

var voltageValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.Atoi(x)
	if err != nil {
		return false
	}
	return true
})


var qValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.Atoi(x)
	if err != nil {
		return false
	}
	return true
})


var simpleValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "simple":
		return true
	case "Simple":
		return true
	case "s":
		return true
	case "sim":
		return true
	case "Sim":
		return true
	}
	return false
})


var fitValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "fit":
		return true
	case "Fit":
		return true
	case "f":
		return true
	}
	return false
})

var calcValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "calc":
		return true
	case "cal":
		return true
	case "c":
		return true
	}
	return false
})
