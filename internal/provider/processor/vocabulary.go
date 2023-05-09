package processor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

type vocabulary struct {
	distanceMatrix map[string]map[string]int
	minLen         int
	maxLen         int
	wordCnt        int
}

func (v *vocabulary) PathLength(word string) (int, error) {
	sum := 0
	bigrams := splitRecursive(word, 2)

	for _, bigram := range bigrams {
		n, err := v.BigramPathLength(bigram)
		if err != nil {
			return 0, err
		}
		sum += n
	}
	return sum, nil
}

func (v *vocabulary) BigramPathLength(bigram string) (int, error) {
	if len(bigram) < 2 {
		return 0, nil
	}

	n, ok := v.distanceMatrix[string(bigram[0])][string(bigram[1])]
	if !ok {
		desc := fmt.Errorf("unknown symbol in bigram:%s", bigram)
		return 0, errors.Join(ErrWrongLetter, desc)
	}
	return n, nil

}

func (v *vocabulary) GapPathLen(word1, word2 string) (int, error) {
	l1, l2 := len(word1), len(word2)
	if l1 < 1 || l2 < 1 {
		return 0, nil
	}

	n, err := v.BigramPathLength(string(word1[l1-1]) + string(word2[0]))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func splitRecursive(str string, size int) []string {
	if len(str) <= size {
		return []string{str}
	}
	return append([]string{str[0:size]}, splitRecursive(str[size-1:], size)...)
}

func (v *vocabulary) ReadFile(fileName string, needSort bool) ([]*wordMetric, error) {
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

		pathLen, err = v.PathLength(word)
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
