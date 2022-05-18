package main

import (
	"strconv"
	"strings"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"noError", "tolerance", "distance", "bruteForce",
		"fileName", "pointTol", "which"},
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
			Prompt:    "Run tolerance study according to the tolerance list",
			Response:  "Do I need this? 2",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"distance": &cli.Item{
			Name:      "distance",
			Prompt:    "Calculate the minimum and maximum distances",
			Response:  "Do I need this 3?",
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
		"pointTol": &cli.Item{
			Name:      "pointTol",
			Prompt:    "Specify a single tolerance value for min/max calculatoin",
			Response:  "Do I need this 4?",
			Value:     "15%",
			Validator: toleranceValidator,
		},
		"which": &cli.Item{
			Name:      "which",
			Prompt:    "Select theta or gamma for tolerance study",
			Response:  "Do I need this 4?",
			Value:     "",
			Validator: whichValidator,
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

var toleranceValidator = cli.ItemValidator(func(x string) bool {
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

var whichValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "theta":
		return true
	case "gamma":
		return true
	default:
		return false
	}
	return false
})
