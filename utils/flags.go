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
	// TODO
	Flags["--filter"] = append(Flags["--filter"], "pixelate")
	Flags["--filter"] = append(Flags["--filter"], "blur")
	Flags["--filter"] = append(Flags["--filter"], "gaus")
	Flags["--filter"] = append(Flags["--filter"], "edge")
	Flags["--filter"] = append(Flags["--filter"], "sharp")

	// TODO
	// Mirror
	Flags["--mirror"] = append(Flags["--mirror"], "h")
	Flags["--mirror"] = append(Flags["--mirror"], "hor")
	Flags["--mirror"] = append(Flags["--mirror"], "horizontal")
	Flags["--mirror"] = append(Flags["--mirror"], "horizontally")
	Flags["--mirror"] = append(Flags["--mirror"], "v")
	Flags["--mirror"] = append(Flags["--mirror"], "ver")
	Flags["--mirror"] = append(Flags["--mirror"], "vertical")
	Flags["--mirror"] = append(Flags["--mirror"], "vertically")

	// TODO
	// Rotate
	Flags["--rotate"] = append(Flags["--rotate"], "right")
	Flags["--rotate"] = append(Flags["--rotate"], "90")
	Flags["--rotate"] = append(Flags["--rotate"], "180")
	Flags["--rotate"] = append(Flags["--rotate"], "270")
	Flags["--rotate"] = append(Flags["--rotate"], "left")
	Flags["--rotate"] = append(Flags["--rotate"], "-90")
	Flags["--rotate"] = append(Flags["--rotate"], "-180")
	Flags["--rotate"] = append(Flags["--rotate"], "-270")

	// TODO
	// Crop
	Flags["--crop"] = append(Flags["--crop"], "")
}

//
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp
// go run main.go header apply --a --b --c --d --e src.bmp dst.bmp
// go run main.go apply header--a --b --c --d --e src.bmp dst.bmp
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp apply  --a --b --c --d --e dst2.bmp
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp apply  --a --b --c --d --e dst2.bmp
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp apply  --a --b --c --d --e src2.bmp dst2.bmp
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp apply  --a --b --c --d --e dst2.bmp
// go run main.go apply --a --b --c --d --e src.bmp dst.bmp apply  --a --b --c --d --e src2.bmp dst2.bmp

func FlagsHandler(args []string) error {
	initFlags()

	if len(os.Args) < 3 {
		Helper(os.Args)
		return ErrHelperActivated
	}
	i := 0

	for j, arg := range args {
		if Contains([]string{"--help", "--h", "-h", "--helps", "--helper"}, arg) {
			Helper(os.Args)
			return ErrHelperActivated
		}
		if j == 0 {
			continue
		}
		if arg == "apply" {
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
		if arg == "header" {
			Prompts[i].Flags = append(Prompts[i].Flags, struct {
				Name  string
				Value string
			}{arg, arg})
			continue
		}

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			flag := parts[0]
			if _, ok := Flags[flag]; !ok {
				switch flag {
				case "--help", "--h", "--helps", "--helper":
					Helper(os.Args)
					return ErrHelperActivated
				default:
					return fmt.Errorf("flag %s is not defined.\n", flag)
				}
			}

			value := ""
			if len(parts) > 1 {
				value = parts[1]
				if len(value) > 0 {
					// Check is value is in Flags[flag] slice, if not print error
					if !Contains(Flags[flag], value) {
						return fmt.Errorf("flag %s value %s is not in allowed values.\n", flag, value)
					}
					if flag == "--rotate" {
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
						if Prompts[i].Flags[0].Name != "--rotate" {
							Prompts[i].Flags = append([]struct {
								Name  string
								Value string
							}{{Name: flag, Value: "0"}}, Prompts[i].Flags...)
						}
						valueInt, err := strconv.Atoi(Prompts[i].Flags[0].Value)
						if err != nil {
							// handle error
						}

						switch value {
						case "right", "90", "-270":
							valueInt = (valueInt + 90) % 360
						case "left", "-90", "270":
							valueInt = (valueInt - 90 + 360) % 360
						}
						valueInt = (valueInt / 90) * 90

						Prompts[i].Flags[0].Value = strconv.Itoa(valueInt)

						continue
					}
					Prompts[i].Flags = append(Prompts[i].Flags, struct {
						Name  string
						Value string
					}{flag, value})
				}
			}
			continue
		}
		// Update current prompt Destination and/or Source
		// Check that it has format .bmp
		if strings.HasSuffix(arg, ".bmp") {
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
					fmt.Printf("\n")
					fmt.Printf("The options are:\n")
					fmt.Printf("  -h, --help     		prints program usage information\n")
					fmt.Printf("  --filter <value>     applies filter to the image\n      ")
					for i, value := range Flags["--filter"] {
						if i != 0 {
							fmt.Printf("| ")
						}
						fmt.Printf("%s ", value)
					}
					fmt.Printf("\n  --mirror <value>     mirrors the image\n      ")
					fmt.Printf("Horizontal | Vertical")
					fmt.Printf("\n  --rotate <value>     rotates the image:\n      ")
					for i, value := range Flags["--rotate"] {
						if i != 0 {
							fmt.Printf("| ")
						}
						fmt.Printf("%s ", value)
					}
					fmt.Printf("\n  --crop <value>     	crops the image\n      ")
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
