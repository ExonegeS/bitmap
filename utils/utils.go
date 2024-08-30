package utils

func GetPixel(data []byte, x, y int) (pixel [3]int) {
	if x < 0 || y < 0 || x >= Header.Width || y >= Header.Height {
		return [3]int{0, 0, 0}
	}
	pixel = [3]int{int(data[x*3+y*3*Header.Width+2]), int(data[x*3+y*3*Header.Width+1]), int(data[x*3+y*3*Header.Width])}
	return pixel
}

func SetPixel(data *[]byte, x, y int, pixel [3]int) (err error) {
	(*data)[Header.Offset+x*3+y*3*Header.Width+2] = byte(pixel[0])
	(*data)[Header.Offset+x*3+y*3*Header.Width+1] = byte(pixel[1])
	(*data)[Header.Offset+x*3+y*3*Header.Width] = byte(pixel[2])

	return nil
}
