package crops

import (
	. "bitmap/utils"
	"strconv"
)

func Crop(data []byte, size []string) (newData []byte, err error) {
	x, err := strconv.Atoi(size[0])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(size[1])
	if err != nil {
		return nil, err
	}
	w := Header.Width - x
	h := Header.Height - y
	if len(size) == 4 {
		w, err = strconv.Atoi(size[2])
		if err != nil {
			return nil, err
		}
		h, err = strconv.Atoi(size[3])
		if err != nil {
			return nil, err
		}
	}
	newData = make([]byte, Header.Offset+w*h*3)

	newHeader := append(data[:2], IntToBytes(Header.Offset+w*h*3)...)
	newHeader = append(newHeader, data[6:18]...)
	newHeader = append(newHeader, IntToBytes(w)...)
	newHeader = append(newHeader, IntToBytes(h)...)
	newHeader = append(newHeader, data[26:34]...)
	newHeader = append(newHeader, IntToBytes(w*h*3)...)
	newHeader = append(newHeader, data[38:54]...)
	oldW, oldH := Header.Width, Header.Height

	Header.Width = w
	Header.Height = h
	Header.FileSizeInBytes = len(newData)
	Header.ImageSizeInBytes = w * h

	copy(newData, newHeader)

	// BUG: Fix some shitty code here

	y = oldH - y
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			pixel := [3]int{int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW+2]), int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW+1]), int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW])}

			err := SetPixel(&newData, j, i, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
