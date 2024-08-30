package mirrors

import . "bitmap/utils"

func Mirror_Axis(data []byte, isVertical bool) (newData []byte, err error) {
	newData = make([]byte, len(data))
	copy(newData, data)
	for x := 0; x < Header.Width; x++ {
		for y := 0; y < Header.Height; y++ {
			pixel := GetPixel(data[Header.Offset:], x, y)
			w, h := Header.Width-x-1, Header.Height-y-1
			if isVertical {
				w = x
			} else {
				h = y
			}
			err = SetPixel(&newData, w, h, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
