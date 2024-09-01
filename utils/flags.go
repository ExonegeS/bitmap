package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Prompts = make([]struct {
	Flags []struct {
		Name  string
		Value string
	}
	Src  string
	Dest string
}, 0)

var (
	ErrHelperActivated = errors.New("Helper activated")
	Flags              = make(map[string][]string) // map of flag name to flag values
)

func initFlags() {
	Prompts = append(Prompts, struct {
		Flags []struct {
			Name  string
			Value string
		}
		Src  string
		Dest string
	}{
		Flags: make([]struct {
			Name  string
			Value string
		}, 0),
		Src:  "",
		Dest: "",
	})
	// Filter
	Flags["--filter"] = append(Flags["--filter"], "blue")
	Flags["--filter"] = append(Flags["--filter"], "red")
	Flags["--filter"] = append(Flags["--filter"], "green")
	Flags["--filter"] = append(Flags["--filter"], "grayscale")
	Flags["--filter"] = append(Flags["--filter"], "negative")
	Flags["--filter"] = append(Flags["--filter"], "pixelate")
	Flags["--filter"] = append(Flags["--filter"], "blur")
	Flags["--filter"] = append(Flags["--filter"], "gaus")
	Flags["--filter"] = append(Flags["--filter"], "edge")
	Flags["--filter"] = append(Flags["--filter"], "sharp")

	// Mirror
	Flags["--mirror"] = append(Flags["--mirror"], "h")
	Flags["--mirror"] = append(Flags["--mirror"], "hor")
	Flags["--mirror"] = append(Flags["--mirror"], "horizontal")
	Flags["--mirror"] = append(Flags["--mirror"], "horizontally")
	Flags["--mirror"] = append(Flags["--mirror"], "v")
	Flags["--mirror"] = append(Flags["--mirror"], "ver")
	Flags["--mirror"] = append(Flags["--mirror"], "vertical")
	Flags["--mirror"] = append(Flags["--mirror"], "vertically")

	// Rotate
	Flags["--rotate"] = append(Flags["--rotate"], "right")
	Flags["--rotate"] = append(Flags["--rotate"], "90")
	Flags["--rotate"] = append(Flags["--rotate"], "180")
	Flags["--rotate"] = append(Flags["--rotate"], "270")
	Flags["--rotate"] = append(Flags["--rotate"], "left")
	Flags["--rotate"] = append(Flags["--rotate"], "-90")
	Flags["--rotate"] = append(Flags["--rotate"], "-180")
	Flags["--rotate"] = append(Flags["--rotate"], "-270")

	// Crop
	Flags["--crop"] = append(Flags["--crop"], "")
}

func FlagsHandler(args []string) error {
	// Initialize flags
	initFlags()

	// Check if the number of arguments is less than 3
	if len(os.Args) < 3 {
		// Call the Helper function and return ErrHelperActivated
		Helper(os.Args)
		return ErrHelperActivated
	}
	i := 0

	// Iterate over the arguments
	for j, arg := range args {
		// Check if the argument is a help flag
		if Contains([]string{"--help", "--h", "-h", "--helps", "--helper"}, arg) {
			// Call the Helper function and return ErrHelperActivated
			Helper(os.Args)
			return ErrHelperActivated
		}
		// Skip the first argument
		if j == 0 {
			continue
		}
		// Check if the argument is "apply"
		if arg == "apply" {
			// If the current prompt has a destination, create a new prompt
			if len(Prompts[i].Dest) > 0 {
				Prompts = append(Prompts, struct {
					Flags []struct {
						Name  string
						Value string
					}
					Src  string
					Dest string
				}{
					Flags: make([]struct {
						Name  string
						Value string
					}, 0),
					Src: Prompts[i].Dest,
				})
				i++
			}
			continue
		}
		// Check if the argument is "header"
		if arg == "header" {
			// Add the header flag to the current prompt
			Prompts[i].Flags = append(Prompts[i].Flags, struct {
				Name  string
				Value string
			}{arg, arg})
			continue
		}

		// Check if the argument is a flag
		if strings.HasPrefix(arg, "--") {
			// Split the flag into name and value
			parts := strings.SplitN(arg, "=", 2)
			flag := parts[0]
			// Check if the flag is valid
			if _, ok := Flags[flag]; !ok {
				switch flag {
				case "--help", "--h", "--helps", "--helper":
					// Call the Helper function and return ErrHelperActivated
					Helper(os.Args)
					return ErrHelperActivated
				default:
					// Return an error if the flag is not defined
					return fmt.Errorf("flag %s is not defined.\n", flag)
				}
			}

			// Get the value of the flag
			value := ""
			if len(parts) < 2 {
				// Return an error if the flag's value is not defined
				return fmt.Errorf("%s flag's value is not defined.\n", flag)
			}
			value = parts[1]
			if len(value) < 1 {
				// Return an error if the flag's value is empty
				return fmt.Errorf("%s flag's value is not defined.\n", flag)
			}
			if len(value) > 0 {
				// Check if the value is valid for the flag
				if flag == "--crop" {
					// Split the value into parts
					n := strings.SplitN(value, "-", 4)
					// Check if the value has the correct format
					if len(n) != 4 && len(n) != 2 {
						return fmt.Errorf("invalid crop value length: %s (%v)\n", value, n)
					}
					// Check if each part of the value is a valid crop value
					for _, v := range n {
						if !isValidCropValue(v) {
							return fmt.Errorf("invalid crop value: %s\n", value)
						}
					}
					// Add the flag to the current prompt
					Prompts[i].Flags = append(Prompts[i].Flags, struct {
						Name  string
						Value string
					}{flag, value})
					continue
				}
				if flag == "--rotate" {
					// Check if the current prompt has a rotate flag
					if len(Prompts) < 1 {
						Prompts[i].Flags = append(Prompts[i].Flags, struct {
							Name  string
							Value string
						}{flag, "0"})
					}
					if len(Prompts[i].Flags) < 1 {
						Prompts[i].Flags = append(Prompts[i].Flags, struct {
							Name  string
							Value string
						}{flag, "0"})
					}
					// Check if the first flag is a rotate flag
					if Prompts[i].Flags[0].Name != "--rotate" {
						// Add a new rotate flag to the beginning of the flags
						Prompts[i].Flags = append([]struct {
							Name  string
							Value string
						}{{Name: flag, Value: "0"}}, Prompts[i].Flags...)
					}
					// Calculate the new rotate value
					newValue := 0
					valueInt, err := strconv.Atoi(Prompts[i].Flags[0].Value)
					if err != nil {
						// Handle error
					}
					switch value {
					case "right":
						newValue = 90
					case "left":
						newValue = -90
					default:
						newValue, err = strconv.Atoi(value)
						if err != nil {
							// Return an error if the rotate value is invalid
							return fmt.Errorf("invalid rotate value: %s\n", value)
						}
					}
					valueInt = (((newValue+valueInt)%360+360)%360 + 45) / 90 * 90

					// Update the rotate flag's value
					Prompts[i].Flags[0].Value = strconv.Itoa(valueInt)

					continue
				}
				// Check if the value is in the allowed values for the flag
				if !Contains(Flags[flag], value) {
					return fmt.Errorf("flag %s value %s is not in allowed values.\n", flag, value)
				}
				// Add the flag to the current prompt
				Prompts[i].Flags = append(Prompts[i].Flags, struct {
					Name  string
					Value string
				}{flag, value})
			}
			continue
		}
		// Check if the argument is a bitmap file
		if strings.HasSuffix(arg, ".bmp") {
			// Update the source and destination of the current prompt
			if Prompts[i].Src == "" {
				Prompts[i].Src = strings.ReplaceAll(arg, "\"", "")
			}
			if Prompts[i].Dest == "" {
				Prompts[i].Dest = strings.ReplaceAll(arg, "\"", "")
			} else {
				Prompts[i].Src = Prompts[i].Dest
				Prompts[i].Dest = strings.ReplaceAll(arg, "\"", "")
			}
		} else {
			// Return an error if the argument is not a bitmap file
			return fmt.Errorf("%v is not bitmap file\n", arg)
		}
	}
	return nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isValidCropValue(v string) bool {
	i, err := strconv.Atoi(v)
	return err == nil && i >= 0
}

func Helper(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: bitmap <command> [arguments]")
		fmt.Println("")
		fmt.Println("The commands are:")
		fmt.Println("  header    prints bitmap file header information")
		fmt.Println("  apply     applies processing to the image and saves it to the file")
	}
	if len(args) > 1 {
		switch args[1] {
		case "header":
			fmt.Println("Usage:")
			fmt.Println("  bitmap header <source_file>")
			fmt.Println("")
			fmt.Println("Description:")
			fmt.Println("  Prints bitmap file header information")

		case "apply":
			for _, arg := range args {
				if Contains([]string{"--help", "--h", "-h", "--helps", "--helper"}, arg) {
					fmt.Printf("Usage: bitmap apply [options] <source_file> <output_file>\n")
					fmt.Printf("	Multiple apply's may be used, in that case previous <output_file>\n")
					fmt.Printf("	will be used as <source_file> for the next apply\n")
					fmt.Printf("	To change the new source file for new apply use both <source_file> <output_file>\n")

					fmt.Printf("\n")
					fmt.Printf("The options are:\n")
					fmt.Printf("  -h, --help	prints program usage information\n")
					fmt.Printf("  --filter=<value>     applies filter to the image\n      ")
					for i, value := range Flags["--filter"] {
						if i != 0 {
							fmt.Printf("| ")
						}
						fmt.Printf("%s ", value)
					}
					fmt.Printf("\n  --mirror=<value>     mirrors the image\n      ")
					fmt.Printf("Horizontal | Vertical")
					fmt.Printf("\n  --rotate=<value>     rotates the image:\n      ")
					for i, value := range Flags["--rotate"] {
						if i != 0 {
							fmt.Printf("| ")
						}
						fmt.Printf("%s ", value)
					}
					fmt.Printf("\n  --crop=<value>     	crops the image\n      ")
					fmt.Printf("OffsetX-OffsetY | OffsetX-OffsetY-Width-Height\n")
					return
				}
			}
			fmt.Println("Usage: bitmap apply [options] <source_file> <output_file>")
			fmt.Println("")
			fmt.Println("The options are:")
			fmt.Println("  -h, --help     prints program usage information")
			fmt.Println("  --filter       applies filter to the image")
			fmt.Println("  --mirror       mirrors the image")
			fmt.Println("  --rotate       rotates the image")
			fmt.Println("  --crop         crops the image")
		default:
			fmt.Println("Usage: bitmap <command> [arguments]")
			fmt.Println("")
			fmt.Println("The commands are:")
			fmt.Println("  header    prints bitmap file header information")
			fmt.Println("  apply     applies processing to the image and saves it to the file")
		}
	}
}
