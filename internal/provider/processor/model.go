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
	MaxLenKnapsack(b1, b2, b3 knapsack) knapsack
}

func New(m map[string]map[string]int, minLen, maxLen, wordCnt int) Processor {
	return &vocabulary{
		distanceMatrix: m,
		minLen:         minLen,
		maxLen:         maxLen,
		wordCnt:        wordCnt,
	}
}

type NewProcessor interface {
	PathLen(word string) (int, error)
	GapPathLen(word1, word2 string) (int, error)
	ReadFile(fileName string, needSort bool) ([]*wordMetric, error)

	calcSet(i, j int, wm *wordMetric, kt *[][]map[uint8]knapsack) error
	FindBestCombination(k knapsack, wm *wordMetric) (bool, knapsack, error)
	ChooseCandidate(candidateKs, upKs, leftKs map[uint8]knapsack) map[uint8]knapsack

	NewKnapsackTable(items []*wordMetric) *[][]map[uint8]knapsack
	MinChoice(kt *[][]map[uint8]knapsack) (knapsack, int)
}

func NewVocab(m []int, minLen, maxLen int, wordCnt uint8) NewProcessor {
	return &vocab{
		distanceArray: m,
		minLen:        minLen,
		maxLen:        maxLen,
		wordCnt:       wordCnt,
	}
}
