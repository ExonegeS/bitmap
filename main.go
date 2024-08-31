package main

import (
	"bitmap/crops"
	"bitmap/filters"
	"bitmap/mirrors"
	"bitmap/rotates"
	. "bitmap/utils"
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	err := FlagsHandler(os.Args)
	if errors.Is(err, ErrHelperActivated) {
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	for _, prompt := range Prompts {
		dataSrc, err := ReadFile(prompt.Src)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		dataDest := make([]byte, 0)
		for _, flag := range prompt.Flags {
			ReadHeader(dataSrc)
			if len(dataDest) > 0 {
				dataSrc = dataDest
			}
			switch flag.Name {
			case "header":
				PrintHeader()
				break
			case "--filter":
				{
					switch flag.Value {
					case "blue":
						dataDest, err = filters.Filter_Blue_Channel(dataSrc)
						break
					case "red":
						dataDest, err = filters.Filter_Red_Channel(dataSrc)
						break
					case "green":
						dataDest, err = filters.Filter_Green_Channel(dataSrc)
						break
					case "negative":
						dataDest, err = filters.Filter_Negative(dataSrc)
						break
					case "grayscale":
						dataDest, err = filters.Filter_Grayscale(dataSrc)
						break
					case "pixelate":
						dataDest, err = filters.Filter_Pixelate(dataSrc, 20)
						break
					case "blur":
						dataDest, err = filters.Filter_Blur(dataSrc, 10)
						break
					case "gaus":
						dataDest, err = filters.Filter_GaussianBlur(dataSrc, 10, 5)
						break
					case "edge":
						dataDest, err = filters.Filter_Edge(dataSrc, true, 10.0)
						break
					case "sharp":
						dataDest, err = filters.Filter_Sharp(dataSrc)
						break
					}
				}
			case "--mirror":
				switch flag.Value {
				case "v", "vert", "vertical", "vertically":
					dataDest, err = mirrors.Mirror_Axis(dataSrc, true)
					break
				case "h", "hor", "horizontal", "horizontally":
					dataDest, err = mirrors.Mirror_Axis(dataSrc, false)
					break
				}
			case "--rotate":
				switch flag.Value {
				case "90":
					dataDest, err = rotates.Image_Transpose(dataSrc)
					dataDest, err = mirrors.Mirror_Axis(dataDest, false)
					break
				case "180", "-180":
					dataDest, err = mirrors.Mirror_Axis(dataSrc, false)
					dataDest, err = mirrors.Mirror_Axis(dataDest, true)
					break
				case "270":
					dataDest, err = rotates.Image_Transpose(dataSrc)
					dataDest, err = mirrors.Mirror_Axis(dataDest, true)
					break
				}
			case "--crop":
				switch len(strings.SplitN(flag.Value, "-", 4)) {
				case 2, 4:
					dataDest, err = crops.Crop(dataSrc, strings.SplitN(flag.Value, "-", 4))
					dataDest, err = mirrors.Mirror_Axis(dataDest, true)
					break
				default:
					fmt.Println("Error: Crop requires two or four arguments.")
					return
				}
			}
		}
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
		if len(dataDest) < 1 {
			CreateFile(dataSrc, prompt.Dest)
		} else {
			CreateFile(dataDest, prompt.Dest)
		}
	}
}
