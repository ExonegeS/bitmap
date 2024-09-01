package utils

func GetPixel(data []byte, x, y int) (pixel [3]int) {
	if x < 0 || y < 0 || x >= Header.Width || y >= Header.Height {
		return [3]int{0, 0, 0}
	}
	location := (x + y*Header.Width) * 3
	pixel = [3]int{int(data[location+2]), int(data[location+1]), int(data[location+0])}
	return pixel
}

func SetPixel(data *[]byte, x, y int, pixel [3]int) (err error) {
	location := (x * 3) + (y * (Header.Width*3 + ((4 - (Header.Width*3)%4) % 4)))
	(*data)[Header.Offset+location+2] = byte(pixel[0])
	(*data)[Header.Offset+location+1] = byte(pixel[1])
	(*data)[Header.Offset+location+0] = byte(pixel[2])

	return nil
}
