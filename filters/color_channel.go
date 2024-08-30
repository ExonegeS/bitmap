package filters

import (
	. "bitmap/utils"
)

func Filter_Red_Channel(data []byte) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			pixel[1] = 0
			pixel[2] = 0
			err = SetPixel(&newData, x, y, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}

func Filter_Green_Channel(data []byte) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			pixel[0] = 0
			pixel[2] = 0
			err = SetPixel(&newData, x, y, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}

func Filter_Blue_Channel(data []byte) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for y := 0; y < Header.Height; y++ {
		for x := 0; x < Header.Width; x++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			pixel[0] = 0
			pixel[1] = 0
			err = SetPixel(&newData, x, y, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
