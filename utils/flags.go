package utils

import (
	"fmt"
	"os"
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

var Flags = make(map[string][]string) // map of flag name to flag values

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
		return fmt.Errorf("not enough arguments in calling  the program.\n")
	}
	i := 0

	for j, arg := range args {
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
				return fmt.Errorf("flag %s is not defined.\n", flag)
			}

			value := ""
			if len(parts) > 1 {
				value = parts[1]
				if len(value) > 0 {
					// Check is value is in Flags[flag] slice, if not print error
					if !Contains(Flags[flag], value) {
						return fmt.Errorf("flag %s value %s is not in allowed values.\n", flag, value)
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
