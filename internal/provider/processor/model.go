package processor

import "errors"

var (
	ErrOpenFile = errors.New("can not open file")
	ErrScanFile = errors.New("can not scan file")
)

type NewProcessor interface {
	PathLen(word string) (int, error)
	GapPathLen(word1, word2 string) (int, error)
	ReadFile(fileName string, needSort bool) ([]*wordMetric, error)

	calcSet(i, j int, wm *wordMetric, kt *[][][]knapsack) error
	FindBestCombination(k knapsack, wm *wordMetric) (bool, knapsack, error)
	ChooseCandidate(candidateKs, upKs, leftKs []knapsack) []knapsack

	KnapsackTable(items []*wordMetric) *[][][]knapsack
	MinChoice(kt *[][][]knapsack) (knapsack, int)
}

func NewVocab(m []int, minLen, maxLen int, wordCnt uint8) NewProcessor {
	return &vocab{
		distanceArray: m,
		minLen:        minLen,
		maxLen:        maxLen,
		wordCnt:       wordCnt,
	}
}
