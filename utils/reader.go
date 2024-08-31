package utils

import (
	"errors"
	"fmt"
	"os"
)

type BMPHeader struct {
	FileType         string
	FileSizeInBytes  int
	Offset           int
	DibHeaderSize    int
	Width            int
	Height           int
	PixelSize        int
	ImageSizeInBytes int
	Compression      int
}

var Header = &BMPHeader{}

func ReadFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return nil, fmt.Errorf("%v file does not exist", filename)
		case errors.Is(err, os.ErrPermission):
			return nil, fmt.Errorf("%v permission denied", filename)
		default:
			return nil, fmt.Errorf("reading %v :unknown error:", filename)
		}
	}
	return data, nil
}

func ReadHeader(data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	Header.FileType = string(data[:2])
	Header.FileSizeInBytes = BytesToInt(data[2:6])
	Header.Offset = BytesToInt(data[10:14])
	Header.DibHeaderSize = BytesToInt(data[14:18])
	Header.Width = BytesToInt(data[18:22])
	Header.Height = BytesToInt(data[22:26])
	Header.PixelSize = BytesToInt(data[28:30])
	Header.ImageSizeInBytes = BytesToInt(data[34:38])
	Header.Compression = BytesToInt(data[30:34])
	if Header.PixelSize != 24 {
		return fmt.Errorf("unsupported pixel size: %v", Header.PixelSize)
	}
	if Header.Compression != 0 {
		return fmt.Errorf("only uncompressed files allowed")
	}

	return nil

	// fmt.Printf("- reserved %v\n", data[6:10]) //#
	// fmt.Printf("- Planes %v\n", data[26:28])
	//
	// fmt.Printf("- XpixelsPerM %v\n", data[38:42])
	// fmt.Printf("- YpixelsPerM %v\n\n", data[42:46])
	// fmt.Printf("- Colors Used %v\n", data[46:50])
	// fmt.Printf("- Important Colors %v\n\n", data[50:54])
}

func PrintHeader() {
	fmt.Printf("BMP Header\n")
	fmt.Printf("- FileType %s\n", Header.FileType)
	fmt.Printf("- FileSizeInBytes %v\n", Header.FileSizeInBytes)
	fmt.Printf("- HeaderSize %v\n", Header.Offset)
	fmt.Printf("DIB Header:\n")
	fmt.Printf("- DibHeaderSize %v\n", Header.DibHeaderSize)
	fmt.Printf("- WidthInPixels %v\n", Header.Width)
	fmt.Printf("- HeightInPixels %v\n", Header.Height)
	fmt.Printf("- PixelSizeInBits %v\n", Header.PixelSize)
	fmt.Printf("- ImageSizeInBytes %v\n", Header.ImageSizeInBytes)
}

func CreateFile(data []byte, filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func BytesToInt(data []byte) int {
	result := int(0)
	for i := 0; i < len(data); i++ {
		result |= int(data[i]) << (8 * i)
	}
	return result
}
