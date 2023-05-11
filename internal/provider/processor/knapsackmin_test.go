//go:build processorTest
// +build processorTest

package processor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnapsackmin(t *testing.T) {
	t.Run("test knapsackmin functions", func(t *testing.T) {
		var (
			maxPathLen, p int
			err           error
			wordMetrics   []*wordMetric
			b             knapsack
			bt            [][]knapsack
		)

		tests := []testParam{
			{
				fileName: "tests/test3.txt",
				minLen:   4,
				maxLen:   6,
				wordCnt:  2,
			},
			{
				fileName: "tests/out5a.txt",
				minLen:   20,
				maxLen:   24,
				wordCnt:  4,
			},
			{
				fileName: "tests/out5b.txt",
				minLen:   20,
				maxLen:   24,
				wordCnt:  4,
			},
		}

		for i, param := range tests {
			dist := getDistanceMatrixForTests()
			t.Run(fmt.Sprintf("Test %d, from file %s", i, param.fileName), func(t *testing.T) {

				v := New(dist, param.minLen, param.maxLen, param.wordCnt)
				wordMetrics, err = v.ReadFile(param.fileName, true)
				assert.NoError(t, err)

				t.Run("KnapsackTable and MinChoice", func(t *testing.T) {
					n := len(wordMetrics)
					bt = v.KnapsackMinTable(wordMetrics)

					if i == 0 {
						printTable(bt, n, param.maxLen)
					}

					for n1, x := range bt {
						for n2, b1 := range x {
							p, err = v.PathLength(b1.GetDescription())
							assert.NoError(t, err)
							if p != b1.pathLen {
								fmt.Printf("AAAA! [%v %v] %v \n %v", n1, n2, b1.GetDescriptionWithSpace(), b1)
							}
							//проверяем правильность pathLen
							assert.Equal(t, p, b1.pathLen)

							//длина слова не больше номера столбца
							assert.Equal(t, true, len(b1.GetDescription()) <= n2)

							// по столбцу значение pathLen не может увеличиваться, только если слово удлиннилось
							if n1 != 0 && n2 != 0 {
								assert.Equal(t, true, len(b1.GetDescription()) >= len(bt[n1-1][n2].GetDescription()))
								if len(b1.GetDescription()) == len(bt[n1-1][n2].GetDescription()) {
									assert.Equal(t, true, b1.pathLen <= bt[n1-1][n2].pathLen)
								}
							}
						}
					}

					b, maxPathLen = v.MinChoice(bt)
					fmt.Printf("pathLen: %d\npassword: %s\nwords: %s\n", maxPathLen, b.GetDescription(), b.GetDescriptionWithSpace())

					if i == 0 {
						assert.Equal(t, 5, maxPathLen)
						assert.Equal(t, "lokip", b.GetDescription())
					}
				})
			})
		}
	})

	t.Run("test MaxLenKnapsack functions", func(t *testing.T) {
		dist := getDistanceMatrixForTests()
		v := New(dist, 0, 0, 0)

		b1 := knapsack{
			items: []*wordMetric{
				{word: "of", pathLen: 5},
				{word: "a", pathLen: 0},
			},
			pathLen: 8,
			count:   2,
		}
		b2 := knapsack{
			items: []*wordMetric{
				{word: "of", pathLen: 5},
			},
			pathLen: 5,
			count:   1,
		}
		b3 := knapsack{
			items: []*wordMetric{
				{word: "a", pathLen: 0},
				{word: "of", pathLen: 5},
			},
			pathLen: 13,
			count:   2,
		}
		assert.Equal(t, b1, v.MaxLenKnapsack(b1, b2, b3))
		assert.Equal(t, b1, v.MaxLenKnapsack(b1, b3, b2))
		assert.Equal(t, b1, v.MaxLenKnapsack(b2, b1, b3))
		assert.Equal(t, b1, v.MaxLenKnapsack(b2, b3, b1))
		assert.Equal(t, b1, v.MaxLenKnapsack(b3, b1, b2))
		assert.Equal(t, b1, v.MaxLenKnapsack(b3, b2, b1))

		b4 := knapsack{
			items: []*wordMetric{
				{word: "a", pathLen: 0},
				{word: "ass", pathLen: 1},
			},
			pathLen: 1,
			count:   2,
		}
		assert.Equal(t, b4, v.MaxLenKnapsack(b1, b4, b3))
		assert.Equal(t, b4, v.MaxLenKnapsack(b1, b3, b4))
		assert.Equal(t, b4, v.MaxLenKnapsack(b4, b1, b3))
		assert.Equal(t, b4, v.MaxLenKnapsack(b4, b3, b1))
		assert.Equal(t, b4, v.MaxLenKnapsack(b3, b1, b4))
		assert.Equal(t, b4, v.MaxLenKnapsack(b3, b4, b1))

	})
}
