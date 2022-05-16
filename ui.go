package main

import (
	"strings"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"noError", "tolerance", "distance", "fileName"},
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
		"fileName": &cli.Item{
			Name:      "fileName",
			Prompt:    "Change the csv file name",
			Response:  "Do I need this 4?",
			Value:     "data.csv",
			Validator: filenameValidator,
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
