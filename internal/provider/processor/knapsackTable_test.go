//go:build processorTest
// +build processorTest

package processor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParam struct {
	fileName string
	minLen   int
	maxLen   int
	wordCnt  int
}

func TestNewKnapsack(t *testing.T) {
	t.Run("test NewKnapsack functions", func(t *testing.T) {

		dist := getDistanceMapForTests()
		v := NewVocab(dist, 0, 0, 0)

		t.Run("FindBestCombination", func(t *testing.T) {

			var (
				err         error
				k, kNew     knapsack
				wm1, wm2    *wordMetric
				wordMetrics []*wordMetric
				needStop    bool
			)

			t.Run("Empty knapsack & empty word", func(t *testing.T) {
				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.Error(t, err)
			})

			t.Run("Empty knapsack & filled word", func(t *testing.T) {
				wm1 = &wordMetric{
					word:    "zas",
					pathLen: 2,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zas", kNew.GetDescription())
				assert.Equal(t, true, needStop)

			})

			t.Run("Knapsack from test1.txt & empty word", func(t *testing.T) {
				wordMetrics, err = v.ReadFile("testdata/test1.txt", false)
				assert.NoError(t, err)

				k = knapsack{
					items: wordMetrics,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm2)
				assert.Error(t, err)
			})

			t.Run("Filled knapsack & filled word", func(t *testing.T) {
				needStop, kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zasaofthebike", kNew.GetDescription())
				assert.Equal(t, false, needStop)

				wm2 = &wordMetric{
					word:    "eew",
					pathLen: 2,
				}

				needStop, kNew, err = v.FindBestCombination(k, wm2)
				assert.NoError(t, err)
				assert.Equal(t, "aofthebikeeew", kNew.GetDescription())
				assert.Equal(t, true, needStop)
			})

		})

		t.Run("ChooseCandidate", func(t *testing.T) {

			k11 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}},
				pathLen: 0,
			}
			k21 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}},
				pathLen: 5,
			}
			k31 := knapsack{
				items:   []*wordMetric{{word: "is", pathLen: 6}},
				pathLen: 6,
			}

			k12 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}, {word: "of", pathLen: 5}},
				pathLen: 13,
			}
			k22 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}, {word: "to", pathLen: 4}},
				pathLen: 13,
			}
			k32 := knapsack{
				items:   []*wordMetric{{word: "to", pathLen: 4}, {word: "of", pathLen: 5}},
				pathLen: 9,
			}

			ks1 := make([]knapsack, 3)
			ks2 := make([]knapsack, 3)
			ks3 := make([]knapsack, 3)

			ks1[1] = k11
			ks1[2] = k12

			ks2[1] = k21
			ks3[2] = k22

			ks3[1] = k31
			ks3[2] = k32

			t.Run("1,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks1, ks2, ks3)
				assert.Equal(t, k21.GetDescription(), ksRes[1].GetDescription())
				assert.Equal(t, k32.GetDescription(), ksRes[2].GetDescription())
			})

			k41 := knapsack{
				items:   []*wordMetric{{word: "the", pathLen: 6}},
				pathLen: 6,
			}
			ks4 := make([]knapsack, 3)
			ks4[1] = k41

			t.Run("4,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks4, ks2, ks3)
				assert.Equal(t, k41.GetDescription(), ksRes[1].GetDescription())
				assert.Equal(t, true, ksRes[2].isEmpty())
			})

		})

		t.Run("KnapsackTable and MinChoice", func(t *testing.T) {
			var (
				maxPathLen, p int
				err           error
				wordMetrics   []*wordMetric
				k             knapsack
				kt            *[][][]knapsack
			)

			tests := []testParam{
				{
					fileName: "testdata/test3.txt",
					minLen:   4,
					maxLen:   6,
					wordCnt:  2,
				},
				{
					fileName: "testdata/out5a.txt",
					minLen:   20,
					maxLen:   24,
					wordCnt:  4,
				},
				//{
				//	fileName: "testdata/out5b.txt",
				//	minLen:   20,
				//	maxLen:   24,
				//	wordCnt:  4,
				//},
			}

			for i, param := range tests {
				dist = getDistanceMapForTests()
				t.Run(fmt.Sprintf("Test %d, from file %s", i, param.fileName), func(t *testing.T) {

					v = NewVocab(dist, param.minLen, param.maxLen, uint8(param.wordCnt))
					wordMetrics, err = v.ReadFile(param.fileName, true)
					assert.NoError(t, err)

					n := len(wordMetrics)
					kt = v.KnapsackTable(wordMetrics)

					if i == 0 {
						printMap(kt, n, param.maxLen, uint8(param.wordCnt))
					}

					for n1, x := range *kt {
						for n2, y := range x {
							for cnt, k1 := range y {
								if !k1.isEmpty() {
									// кол-во слов = индексу map
									assert.Equal(t, len(k1.items), cnt)
								}
								p, err = v.PathLen(k1.GetDescription())
								assert.NoError(t, err)
								if p != k1.pathLen {
									fmt.Printf("AAAA! [%v %v] %v \n %v", n1, n2, k1.GetDescriptionWithSpace(), k1)
								}
								//проверяем правильность pathLen
								assert.Equal(t, p, k1.pathLen)

								//длина слова не больше номера столбца
								assert.Equal(t, true, len(k1.GetDescription()) <= n2)

								// по столбцу значение pathLen не может увеличиваться, только если слово удлиннилось
								if n1 != 0 && n2 != 0 {
									k2 := (*kt)[n1-1][n2][cnt]
									assert.Equal(t, true, len(k1.GetDescription()) >= len(k2.GetDescription()))
									if len(k1.GetDescription()) == len(k2.GetDescription()) {
										assert.Equal(t, true, k1.pathLen <= k2.pathLen)
									}
								}
							}
						}
					}

					k, maxPathLen = v.MinChoice(kt)
					fmt.Printf("pathLen: %d\npassword: %s\nwords: %s\n", maxPathLen, k.GetDescription(), k.GetDescriptionWithSpace())

					if i == 0 {
						assert.Equal(t, 4, maxPathLen)
						assert.Equal(t, "ploki", k.GetDescription())
					}
				})
			}
		})
	})
}

/*
func TestMain1(t *testing.T) {

		t.Run("test like main", func(t *testing.T) {
			var (
				maxPathLen  int
				err         error
				wordMetrics []*wordMetric
				k           knapsack
				kt          *[][][]knapsack
			)

			testdata := []testParam{
				{
					fileName: "testdata/10000.txt",
					minLen:   20,
					maxLen:   24,
					wordCnt:  4,
				},
			}

			for i, param := range testdata {
				dist := getDistanceMapForTests()
				t.Run(fmt.Sprintf("Test %d, from file %s", i, param.fileName), func(t *testing.T) {

					v := NewVocab(dist, param.minLen, param.maxLen, uint8(param.wordCnt))
					wordMetrics, err = v.ReadFile(param.fileName, true)
					assert.NoError(t, err)

					kt = v.KnapsackTable(wordMetrics)

					k, maxPathLen = v.MinChoice(kt)
					fmt.Printf("pathLen: %d\npassword: %s\nwords: %s\n", maxPathLen, k.GetDescription(), k.GetDescriptionWithSpace())
				})
			}
		})
	}
*/
func printMap(kt *[][][]knapsack, n, k int, wordCnt uint8) {
	for i := 0; i < n+1; i++ {
		for j := 0; j < k+1; j++ {
			s := ""
			for cnt := uint8(1); cnt <= wordCnt; cnt++ {
				if b := (*kt)[i][j][uint8(cnt)]; !b.isEmpty() {
					s = fmt.Sprintf("%s[%d:%d=%v]", s, cnt, b.pathLen, b.GetDescription())
				} else {
					s = fmt.Sprintf("%s[%d:_=_]", s, cnt)
				}
			}
			fmt.Printf("%s \t", s)
		}
		fmt.Println()
	}
}
