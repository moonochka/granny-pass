package processor

import "errors"

var (
	ErrWrongLetter = errors.New("letter is not set in distance matrix")
	ErrOpenFile    = errors.New("can not open file")
	ErrScanFile    = errors.New("can not scan file")
)

type Processor interface {
	PathLength(word string) (int, error)
	BigramPathLength(bigram string) (int, error)
	GapPathLen(word1, word2 string) (int, error)
	ReadFile(fileName string, needSort bool) ([]*wordMetric, error)

	KnapsackTable(k int, items []*wordMetric) [][]knapsack
	MaxChoice(bc [][]knapsack) (knapsack, int)

	KnapsackMinTable(items []*wordMetric) [][]knapsack
	MinChoice(bc [][]knapsack) (knapsack, int)
	MaxLenBackpack(b1, b2, b3 knapsack) knapsack
}

func New(m map[string]map[string]int, minLen, maxLen, wordCnt int) Processor {
	return &vocabulary{
		distanceMatrix: m,
		minLen:         minLen,
		maxLen:         maxLen,
		wordCnt:        wordCnt,
	}
}
