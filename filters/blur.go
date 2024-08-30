package filters

import (
	. "bitmap/utils"
	"math"
)

func Filter_Blur(data []byte, radius int) (newData []byte, err error) {
	kernel := make([][]float64, 0)
	row := make([]float64, 0)
	val := 1 / float64((2*radius+1)*(2*radius+1))
	row = append(row, val)
	for j := 0; j < radius*2; j++ {
		row = append(row, val)
	}
	for i := 0; i < radius*2; i++ {
		kernel = append(kernel, row)
	}
	kernel = append(kernel, row)
	newData, err = Convolution(data, kernel)
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func Filter_GaussianBlur(data []byte, radius int, sigma float64) (newData []byte, err error) {
	// Create a 2D kernel with Gaussian distribution values
	kernel := make([][]float64, 2*radius+1)
	for i := 0; i < 2*radius+1; i++ {
		kernel[i] = make([]float64, 2*radius+1)
		for j := 0; j < 2*radius+1; j++ {
			x := float64(i - radius)
			y := float64(j - radius)
			kernel[i][j] = (1 / (2 * math.Pi * sigma * sigma)) * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
		}
	}

	// Apply the Gaussian blur using the Convolution function
	newData, err = Convolution(data, kernel)
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func Convolution(data []byte, kernel [][]float64) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := [3]float64{0, 0, 0}
			for i := -len(kernel) / 2; i <= len(kernel)/2; i++ {
				for j := -len(kernel[0]) / 2; j <= len(kernel[0])/2; j++ {
					if x+j >= 0 && x+j < Header.Width && y+i >= 0 && y+i < Header.Height {
						pixelKernel := GetPixel(data[54:], x+j, y+i)
						for k := 0; k < 3; k++ {
							pixel[k] += float64(pixelKernel[k]) * kernel[i+len(kernel)/2][j+len(kernel[0])/2]
						}
					}
				}
			}
			for k := 0; k < 3; k++ {
				pixel[k] = float64(int(math.Max(0, math.Min(255, pixel[k]))))
			}
			err = SetPixel(&newData, x, y, [3]int{int(pixel[0]), int(pixel[1]), int(pixel[2])})
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
