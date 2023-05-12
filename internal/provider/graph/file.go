package graph

import (
	"encoding/json"
	"io"
	"os"
)

func BigramDistanceMap(m map[string]map[string]int) map[string]int {
	res := make(map[string]int)

	for k1, v1 := range m {
		for k2, v2 := range v1 {
			res[k1+k2] = v2
		}
	}
	return res
}

func BigramDistanceArray(m map[string]map[string]int) []int {
	res := make([]int, 32*32)

	for k1, v1 := range m {
		for k2, v2 := range v1 {
			s1 := int(k1[0]) - int('a')
			s2 := int(k2[0]) - int('a')
			res[s1<<5+s2] = v2
		}
	}
	return res
}

func SaveToJson(m map[string]int, filename string) error {
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

func ReadFromJson(filename string) (map[string]int, error) {
	var res map[string]int

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
