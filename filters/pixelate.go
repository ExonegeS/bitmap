package filters

import (
	. "bitmap/utils"
)

func Filter_Pixelate(data []byte, blockSize int) (newData []byte, err error) {
	if blockSize <= 0 {
		blockSize = 20 // default block size
	}

	newData = make([]byte, len(data))
	copy(newData, data)

	for y := blockSize / 2; y < Header.Height+blockSize/2; y += blockSize {
		for x := blockSize / 2; x < Header.Width+blockSize/2; x += blockSize {
			// calculate the average color of the block
			var r, g, b int
			for dy := 0; dy < blockSize; dy++ {
				for dx := 0; dx < blockSize; dx++ {
					pixel := GetPixel(data[Header.Offset:], x+dx-blockSize/2, y+dy-blockSize/2)
					r += pixel[0]
					g += pixel[1]
					b += pixel[2]
				}
			}
			r /= blockSize * blockSize
			g /= blockSize * blockSize
			b /= blockSize * blockSize
			pixel := [3]int{r, g, b}
			for dy := 0; dy < blockSize; dy++ {
				for dx := 0; dx < blockSize; dx++ {
					if x+dx-blockSize/2 < 0 || x+dx-blockSize/2 >= Header.Width || y+dy-blockSize/2 < 0 || y+dy-blockSize/2 >= Header.Height {
						continue
					}
					SetPixel(&newData, x+dx-blockSize/2, y+dy-blockSize/2, pixel)
				}
			}
		}
	}

	return newData, nil
}
