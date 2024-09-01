package main

import (
	"bitmap/crops"
	"bitmap/filters"
	"bitmap/mirrors"
	"bitmap/rotates"
	"errors"
	"fmt"
	"os"
	"strings"

	. "bitmap/utils"
)

func main() {
	// Handle flags
	err := FlagsHandler(os.Args)
	if errors.Is(err, ErrHelperActivated) {
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Process each prompt
	for _, prompt := range Prompts {
		// Read source file
		dataSrc, err := ReadFile(prompt.Src)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Initialize destination data and print header function
		dataDest := make([]byte, 0)
		var printHeaderFunc func()

		// Process each flag
		for _, flag := range prompt.Flags {
			// Read header
			ReadHeader(dataSrc)
			if len(dataDest) > 0 {
				dataSrc = dataDest
			}

			// Handle flag
			switch flag.Name {
			case "header":
				printHeaderFunc = PrintHeader
			case "--filter":
				dataDest, err = applyFilter(dataSrc, flag.Value)
			case "--mirror":
				dataDest, err = applyMirror(dataSrc, flag.Value)
			case "--rotate":
				dataDest, err = applyRotate(dataSrc, flag.Value)
			case "--crop":
				dataDest, err = applyCrop(dataSrc, flag.Value)
			}

			// Check for error
			if err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		}

		// Print header if necessary
		if printHeaderFunc != nil {
			printHeaderFunc()
		}

		// Create destination file
		if len(dataDest) < 1 {
			CreateFile(dataSrc, prompt.Dest)
		} else {
			CreateFile(dataDest, prompt.Dest)
		}
	}
}

func applyFilter(dataSrc []byte, value string) ([]byte, error) {
	switch value {
	case "blue":
		return filters.Filter_Blue_Channel(dataSrc)
	case "red":
		return filters.Filter_Red_Channel(dataSrc)
	case "green":
		return filters.Filter_Green_Channel(dataSrc)
	case "negative":
		return filters.Filter_Negative(dataSrc)
	case "grayscale":
		return filters.Filter_Grayscale(dataSrc)
	case "pixelate":
		return filters.Filter_Pixelate(dataSrc, 20)
	case "blur":
		return filters.Filter_Blur(dataSrc, 10)
	case "gaus":
		return filters.Filter_GaussianBlur(dataSrc, 10, 5)
	case "edge":
		return filters.Filter_Edge(dataSrc, true, 10.0)
	case "sharp":
		return filters.Filter_Sharp(dataSrc)
	default:
		return nil, errors.New("unknown filter")
	}
}

func applyMirror(dataSrc []byte, value string) ([]byte, error) {
	switch value {
	case "v", "vert", "vertical", "vertically":
		return mirrors.Mirror_Axis(dataSrc, true)
	case "h", "hor", "horizontal", "horizontally":
		return mirrors.Mirror_Axis(dataSrc, false)
	default:
		return nil, errors.New("unknown mirror")
	}
}

func applyRotate(dataSrc []byte, value string) ([]byte, error) {
	fmt.Println(value)
	switch value {
	case "0":
		return dataSrc, nil
	case "90":
		dataDest, err := rotates.Image_Transpose(dataSrc)
		if err != nil {
			return nil, err
		}
		return mirrors.Mirror_Axis(dataDest, false)
	case "180", "-180":
		dataDest, err := mirrors.Mirror_Axis(dataSrc, false)
		if err != nil {
			return nil, err
		}
		return mirrors.Mirror_Axis(dataDest, true)
	case "270":
		dataDest, err := rotates.Image_Transpose(dataSrc)
		if err != nil {
			return nil, err
		}
		return mirrors.Mirror_Axis(dataDest, true)
	default:
		return nil, errors.New("unknown rotation")
	}
}

func applyCrop(dataSrc []byte, value string) ([]byte, error) {
	split := strings.SplitN(value, "-", 4)
	if len(split) != 2 && len(split) != 4 {
		return nil, errors.New("invalid crop arguments")
	}
	return crops.Crop(dataSrc, split)
}
