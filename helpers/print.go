package helpers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"
)

// PrintInfo print Info
func PrintInfo() {
	f := figure.NewColorFigure("ia", "big", "red", true)
	figletStr := f.String()
	fmt.Println(figletStr)
	fmt.Println()
}

func PrintFormat(obj interface{}, format string) error {
	var err error
	switch format {
	case "json":
		err = PrintJSON(obj)
	case "yaml":
		err = PrintYAML(obj)

	case "csv":
		err = PrintCSV(obj)

	default:
		return errors.New("unknown format")
	}
	if err != nil {
		return err
	}
	return nil
}

func PrintJSON(obj interface{}) error {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintYAML(obj interface{}) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintCSV(obj interface{}) error {
	csv, err := gocsv.MarshalBytes(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(csv))
	return nil
}
