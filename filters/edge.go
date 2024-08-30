package filters

import . "bitmap/utils"

var kernelVertical = [][]float64{
	{1.0 / 4, 1.0 / 2, 1.0 / 4},
	{0.0, 0.0, 0.0},
	{-1.0 / 4, -1.0 / 2, -1.0 / 4},
}

var kernelHorizontal = [][]float64{
	{1.0 / 4, 0.0, -1.0 / 4},
	{1.0 / 2, 0.0, -1.0 / 2},
	{1.0 / 4, 0.0, -1.0 / 4},
}

func Filter_Edge(data []byte, direction bool, strength float64) (newData []byte, err error) {
	if direction {
		newData, err = EdgeConvolution(data, kernelVertical, strength)
	} else {
		newData, err = EdgeConvolution(data, kernelHorizontal, strength)
	}
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func EdgeConvolution(data []byte, kernel [][]float64, strength float64) (newData []byte, err error) {
	// Calculate the sum of all elements in the kernel
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := [3]float64{0, 0, 0}
			pixelSum := 0.0
			for i := -len(kernel) / 2; i <= len(kernel)/2; i++ {
				for j := -len(kernel[0]) / 2; j <= len(kernel[0])/2; j++ {
					if x+j >= 0 && x+j < Header.Width && y+i >= 0 && y+i < Header.Height {
						pixelKernel := GetPixel(data[54:], x+j, y+i)
						grey := (float64(pixelKernel[0]) + float64(pixelKernel[1]) + float64(pixelKernel[2])) / 3
						pixelSum += grey * kernel[i+len(kernel)/2][j+len(kernel[0])/2]
					}
				}
			}
			if pixelSum < -strength {
				pixel = [3]float64{255 * pixelSum, 0, 0}
			}
			if pixelSum > strength {
				pixel = [3]float64{0, 0, 255 * pixelSum}
			}

			err = SetPixel(&newData, x, y, [3]int{int(pixel[0]), int(pixel[1]), int(pixel[2])})
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
