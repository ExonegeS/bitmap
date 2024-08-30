package main

import (
	"bitmap/filters"
	"bitmap/mirrors"
	. "bitmap/utils"
	"fmt"
	"os"
)

func main() {
	err := FlagsHandler(os.Args)
	if err != nil {
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

		/*
			switch Flags["branch"][0] {
			case "header":
				{
					err = (ReadHeader(data))
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						os.Exit(1)
					}
				}
				break
			case "apply":
				{
					for
					newData, err := filters.Filter_Negative(data)
					if err != nil {
						fmt.Printf("Error: %v", err)
						os.Exit(1)
					}
					CreateFile(newData, Tail[1])
				}
			}

		*/
	}
}
