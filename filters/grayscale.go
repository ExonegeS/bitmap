package filters

import . "bitmap/utils"

func Filter_Grayscale(data []byte) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			gray := (pixel[0] + pixel[1] + pixel[2]) / 3
			pixel[0] = gray
			pixel[1] = gray
			pixel[2] = gray
			err = SetPixel(&newData, x, y, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
