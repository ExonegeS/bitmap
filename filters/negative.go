package filters

import . "bitmap/utils"

func Filter_Negative(data []byte) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			pixel[0] = 255 - pixel[0]
			pixel[1] = 255 - pixel[1]
			pixel[2] = 255 - pixel[2]
			err = SetPixel(&newData, x, y, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
