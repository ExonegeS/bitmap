package filters

var kernelSharp = [][]float64{
	{0, 0, -1.0 / 5, 0, 0},
	{0, -1.0 / 5, -1.0 / 2, -1.0 / 5, 0},
	{-1.0 / 5, -1.0 / 2, 4.6, -1.0 / 2, -1.0 / 5},
	{0, -1.0 / 5, -1.0 / 2, -1.0 / 5, 0},
	{0, 0, -1.0 / 5, 0, 0},
}

func Filter_Sharp(data []byte) (newData []byte, err error) {
	newData, err = Convolution(data, kernelSharp)
	if err != nil {
		return nil, err
	}
	return newData, nil
}
