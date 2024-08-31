package rotates

import (
	. "bitmap/utils"
	"fmt"
)

func Image_Transpose(data []byte) (newData []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	newData = make([]byte, len(data))
	oldH, oldW := Header.Height, Header.Width
	Header.Width, Header.Height = Header.Height, Header.Width
	w, h := data[18:22], make([]byte, 4)
	copy(h, data[22:26])
	_, _ = w, h
	data = append(data[:18], append(h, append(w, data[26:]...)...)...)
	copy(newData, data)

	for y := 0; y < oldH; y++ {
		for x := 0; x < oldW; x++ {
			pixel := [3]int{int(data[Header.Offset+x*3+y*3*oldW+2]), int(data[Header.Offset+x*3+y*3*oldW+1]), int(data[Header.Offset+x*3+y*3*oldW])}

			err = SetPixel(&newData, oldH-y-1, oldW-x-1, pixel)
			if err != nil {
				return nil, err
			}
		}
	}
	return newData, nil
}
