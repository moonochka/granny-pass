package graph

import (
	"encoding/json"
	"io"
	"os"
)

func BigramDistanceArray(m map[string]map[string]int) []int {
	res := make([]int, 32*32)

	for k1, v1 := range m {
		for k2, v2 := range v1 {
			res[getIndex(k1[0], k2[0])] = v2
		}
	}
	return res
}

func getIndex(a, b uint8) int {
	a0 := int(a) - int('a')
	b0 := int(b) - int('a')
	return a0<<5 + b0
}

func SaveToJson(m []int, filename string) error {
	jsonData, err := json.Marshal(m)

	if err != nil {
		return err
	}

	jsonFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		_ = jsonFile.Close()
	}()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromJson(filename string) ([]int, error) {
	var res []int

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = jsonFile.Close()
	}()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
