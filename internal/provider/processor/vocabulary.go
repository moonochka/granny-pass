package processor

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type vocab struct {
	distanceArray []int
	minLen        int
	maxLen        int
	wordCnt       uint8
}

func symbolOffset(s uint8) int {
	return int(s) - int('a')
}

func getIndex(a, b int) int {
	res := int(uint16((a << 5) + b))
	return res
}

func getIndexBigram(s1, s2 uint8) int {
	return getIndex(symbolOffset(s1), symbolOffset(s2))
}

func (v *vocab) PathLen(word string) (int, error) {
	sum := 0

	for i := 0; i < (len(word) - 1); i++ {
		pathLen := v.distanceArray[getIndexBigram(word[i], word[i+1])]

		sum += pathLen
	}
	return sum, nil
}

func (v *vocab) GapPathLen(word1, word2 string) (int, error) {
	l1, l2 := len(word1), len(word2)
	if l1 < 1 || l2 < 1 {
		return 0, nil
	}

	n := v.distanceArray[getIndexBigram(word1[l1-1], word2[0])]
	return n, nil
}

func (v *vocab) ReadFile(fileName string, needSort bool) ([]*wordMetric, error) {
	var (
		word    string
		pathLen int
		res     []*wordMetric
	)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("file name:%s", fileName)
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)

	for Scanner.Scan() {
		word = Scanner.Text()

		pathLen, err = v.PathLen(word)
		if err != nil {
			return nil, err
		}

		wm := wordMetric{
			word:    word,
			pathLen: pathLen,
		}
		res = append(res, &wm)
	}

	if err = Scanner.Err(); err != nil {
		return nil, ErrScanFile
	}

	if needSort {
		sort.Slice(res, func(i, j int) bool {
			return len(res[i].word) < len(res[j].word)
		})
	}
	return res, nil
}
