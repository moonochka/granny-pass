package processor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

type vocab struct {
	distanceMap map[string]int
	minLen      int
	maxLen      int
	wordCnt     int
}

func (v *vocab) PathLen(word string) (int, error) {
	var bigram string
	sum := 0

	for i := 0; i < (len(word) - 1); i++ {
		bigram = string(word[i]) + string(word[i+1])

		pathLen, ok := v.distanceMap[bigram]
		if !ok {
			return 0, fmt.Errorf("unknown symbol in bigram:%s", bigram)
		}

		sum += pathLen
	}
	return sum, nil
}

//
//func (v *vocab) BigramPathLength(bigram string) (int, error) {
//	if len(bigram) < 2 {
//		return 0, nil
//	}
//
//	n, ok := v.distanceMatrix[string(bigram[0])][string(bigram[1])]
//	if !ok {
//		desc := fmt.Errorf("unknown symbol in bigram:%s", bigram)
//		return 0, errors.Join(ErrWrongLetter, desc)
//	}
//	return n, nil
//
//}

func (v *vocab) GapPathLen(word1, word2 string) (int, error) {
	l1, l2 := len(word1), len(word2)
	if l1 < 1 || l2 < 1 {
		return 0, nil
	}

	bigram := string(word1[l1-1]) + string(word2[0])

	n, ok := v.distanceMap[bigram]
	if !ok {
		return 0, fmt.Errorf("unknown symbol in bigram:%s", bigram)
	}

	return n, nil
}

//func splitRecursive(str string, size int) []string {
//	if len(str) <= size {
//		return []string{str}
//	}
//	return append([]string{str[0:size]}, splitRecursive(str[size-1:], size)...)
//}

func (v *vocab) ReadFile(fileName string, needSort bool) ([]*wordMetric, error) {
	var (
		word    string
		pathLen int
		res     []*wordMetric
	)

	file, err := os.Open(fileName)
	if err != nil {
		desc := fmt.Errorf("file name:%s", fileName)
		return nil, errors.Join(ErrOpenFile, desc)
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
			return len(res[i].word) > len(res[j].word)
		})
	}
	return res, nil
}
