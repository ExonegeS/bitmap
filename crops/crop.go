package crops

import (
	"fmt"
	"strconv"

	. "bitmap/utils"
)

func Crop(data []byte, size []string) (newData []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
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
	w = TrimNumber(w, 0, Header.Width-x)
	h = TrimNumber(h, 0, Header.Height-y)
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("Crop size is invalid\n")
	}
	padding := (4 - (w*3)%4) % 4
	newData = make([]byte, Header.Offset+(w*3+padding)*h)

	newHeader := append(data[:2], IntToBytes(Header.Offset+(w*3+padding)*h)...)
	newHeader = append(newHeader, data[6:18]...)
	newHeader = append(newHeader, IntToBytes(w)...)
	newHeader = append(newHeader, IntToBytes(h)...)
	newHeader = append(newHeader, data[26:34]...)
	newHeader = append(newHeader, IntToBytes(w*h*3)...)
	newHeader = append(newHeader, data[38:54]...)
	oldW, oldH := Header.Width, Header.Height
	oldPadding := (4 - (oldW*3)%4) % 4
	Header.Width = w
	Header.Height = h
	Header.FileSizeInBytes = len(newData)
	Header.ImageSizeInBytes = (w*3 + padding) * h

	copy(newData, newHeader)

	y = oldH - y - 1

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			pixel := [3]int{int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW+2+oldPadding*(y-i)]), int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW+1+oldPadding*(y-i)]), int(data[Header.Offset+(x+j)*3+(y-i)*3*oldW+oldPadding*(y-i)])}
			if (i+j)%10 == 0 && false {
				pixel[0] = 255
			}
			err := SetPixel(&newData, j, Header.Height-i-1, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}

func TrimNumber(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
