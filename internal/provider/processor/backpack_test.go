//go:build processorTest
// +build processorTest

package processor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackpack(t *testing.T) {
	t.Run("test knapsack functions", func(t *testing.T) {
		var (
			maxPathLen, pathLen int
			err                 error
			wordMetrics, wordMs []*wordMetric
			b, b1               knapsack
			bt                  [][]knapsack
		)

		dist := getDistanceMatrixForTests()
		v := New(dist, 0, 0, 0)

		t.Run("test.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("tests/test.txt", false)
			assert.NoError(t, err)

			t.Run("GetDescription", func(t *testing.T) {
				t.Run("empty knapsack", func(t *testing.T) {
					assert.Equal(t, "", b.GetDescription())
				})
				t.Run("knapsack from test.txt", func(t *testing.T) {
					b = knapsack{
						items: wordMetrics,
					}

					assert.Equal(t, "aofbikeoatsbiotransformation", b.GetDescription())
				})
			})

			t.Run("lastWord", func(t *testing.T) {

				t.Run("empty knapsack", func(t *testing.T) {
					assert.Equal(t, "", b1.lastWord())
				})

				t.Run("1 word knapsack", func(t *testing.T) {
					word := "hello"
					pathLen, err = v.PathLength(word)
					wm := wordMetric{
						word:    word,
						len:     len(word),
						pathLen: pathLen,
					}
					wordMs = append(wordMs, &wm)
					b1 = knapsack{
						items: wordMs,
					}
					assert.Equal(t, word, b1.lastWord())
				})

				t.Run("knapsack from test.txt", func(t *testing.T) {
					b = knapsack{
						items: wordMetrics,
					}

					assert.Equal(t, "biotransformation", b.lastWord())
				})
			})

			t.Run("KnapsackTable and MaxChoice", func(t *testing.T) {
				k := 12 //24

				bt = v.KnapsackTable(k, wordMetrics)

				//printTable(bt, n, k)

				b, maxPathLen = v.MaxChoice(bt)
				assert.Equal(t, 30, maxPathLen)
				assert.Equal(t, "aofbikeoats", b.GetDescription())
			})
		})

		t.Run("test1.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("tests/test1.txt", false)
			assert.NoError(t, err)

			t.Run("GetDescription", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "aofthebike", b.GetDescription())
			})

			t.Run("lastWord", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "bike", b.lastWord())
			})

			t.Run("KnapsackTable and MaxChoice", func(t *testing.T) {
				k := 5 // 24

				bt = v.KnapsackTable(k, wordMetrics)

				//printTable(bt, n, k)

				assert.Equal(t, "", bt[0][0].GetDescription())
				assert.Equal(t, 0, bt[0][0].pathLen)
				assert.Equal(t, 0, bt[0][0].count)

				assert.Equal(t, "", bt[0][1].GetDescription())
				assert.Equal(t, 0, bt[0][1].pathLen)

				assert.Equal(t, "", bt[0][2].GetDescription())
				assert.Equal(t, 0, bt[0][2].pathLen)

				assert.Equal(t, "", bt[0][3].GetDescription())
				assert.Equal(t, 0, bt[0][3].pathLen)

				assert.Equal(t, "", bt[0][4].GetDescription())
				assert.Equal(t, 0, bt[0][4].pathLen)

				assert.Equal(t, "", bt[0][5].GetDescription())
				assert.Equal(t, 0, bt[0][5].pathLen)

				/*========*/

				assert.Equal(t, "", bt[1][0].GetDescription())
				assert.Equal(t, 0, bt[1][0].pathLen)

				assert.Equal(t, "a", bt[1][1].GetDescription())
				assert.Equal(t, 0, bt[1][1].pathLen)
				assert.Equal(t, 1, bt[1][1].count)

				assert.Equal(t, "a", bt[1][2].GetDescription())
				assert.Equal(t, 0, bt[1][2].pathLen)

				assert.Equal(t, "a", bt[1][3].GetDescription())
				assert.Equal(t, 0, bt[1][3].pathLen)

				assert.Equal(t, "a", bt[1][4].GetDescription())
				assert.Equal(t, 0, bt[1][4].pathLen)

				assert.Equal(t, "a", bt[1][5].GetDescription())
				assert.Equal(t, 0, bt[1][5].pathLen)

				/*========*/

				assert.Equal(t, "", bt[2][0].GetDescription())
				assert.Equal(t, 0, bt[2][0].pathLen)

				assert.Equal(t, "a", bt[2][1].GetDescription())
				assert.Equal(t, 0, bt[2][1].pathLen)
				assert.Equal(t, 1, bt[2][1].count)

				assert.Equal(t, "of", bt[2][2].GetDescription())
				assert.Equal(t, 5, bt[2][2].pathLen)
				assert.Equal(t, 1, bt[2][2].count)

				assert.Equal(t, "aof", bt[2][3].GetDescription())
				assert.Equal(t, 5, bt[2][3].pathLen)
				assert.Equal(t, 2, bt[2][3].count)

				assert.Equal(t, "aof", bt[2][4].GetDescription())
				assert.Equal(t, 5, bt[2][4].pathLen)

				assert.Equal(t, "aof", bt[2][5].GetDescription())
				assert.Equal(t, 5, bt[2][5].pathLen)

				/*========*/

				assert.Equal(t, "", bt[3][0].GetDescription())
				assert.Equal(t, 0, bt[3][0].pathLen)
				assert.Equal(t, 0, bt[3][0].count)

				assert.Equal(t, "a", bt[3][1].GetDescription())
				assert.Equal(t, 0, bt[3][1].pathLen)
				assert.Equal(t, 1, bt[3][1].count)

				assert.Equal(t, "of", bt[3][2].GetDescription())
				assert.Equal(t, 5, bt[3][2].pathLen)
				assert.Equal(t, 1, bt[3][2].count)

				assert.Equal(t, "the", bt[3][3].GetDescription())
				assert.Equal(t, 6, bt[3][3].pathLen)
				assert.Equal(t, 1, bt[3][3].count)

				assert.Equal(t, "athe", bt[3][4].GetDescription())
				assert.Equal(t, 6, bt[3][4].pathLen)
				assert.Equal(t, 2, bt[3][4].count)

				assert.Equal(t, "ofthe", bt[3][5].GetDescription())
				assert.Equal(t, 11, bt[3][5].pathLen)
				assert.Equal(t, 2, bt[3][5].count)

				/*========*/

				assert.Equal(t, "", bt[4][0].GetDescription())
				assert.Equal(t, 0, bt[4][0].pathLen)
				assert.Equal(t, 0, bt[4][0].count)

				assert.Equal(t, "a", bt[4][1].GetDescription())
				assert.Equal(t, 0, bt[4][1].pathLen)
				assert.Equal(t, 1, bt[4][1].count)

				assert.Equal(t, "of", bt[4][2].GetDescription())
				assert.Equal(t, 5, bt[4][2].pathLen)
				assert.Equal(t, 1, bt[4][2].count)

				assert.Equal(t, "the", bt[4][3].GetDescription())
				assert.Equal(t, 6, bt[4][3].pathLen)
				assert.Equal(t, 1, bt[4][3].count)

				assert.Equal(t, "bike", bt[4][4].GetDescription())
				assert.Equal(t, 10, bt[4][4].pathLen)
				assert.Equal(t, 1, bt[4][4].count)

				assert.Equal(t, "ofthe", bt[4][5].GetDescription())
				assert.Equal(t, 11, bt[4][5].pathLen)
				assert.Equal(t, 2, bt[4][5].count)

				/*========*/

				b, maxPathLen = v.MaxChoice(bt)
				assert.Equal(t, 11, maxPathLen)
				assert.Equal(t, "ofthe", b.GetDescription())
			})
		})

		t.Run("test2.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("tests/test2.txt", false)
			assert.NoError(t, err)

			t.Run("GetDescription", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "aofbiotransformationbikeoatmeals", b.GetDescription())
			})

			t.Run("lastWord", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "oatmeals", b.lastWord())
			})

			t.Run("KnapsackTable and MaxChoice", func(t *testing.T) {
				k := 20 //24

				bt = v.KnapsackTable(k, wordMetrics)

				//printTable(bt, n, k)

				b, maxPathLen = v.MaxChoice(bt)
				assert.Equal(t, 63, maxPathLen)
				assert.Equal(t, "aofbiotransformation", b.GetDescription())
			})
		})
	})
}

func printTable(bt [][]knapsack, n, k int) {
	for i := 0; i < n+1; i++ {
		for j := 0; j < k+1; j++ {
			fmt.Printf("[%v=%v=%v \t]", bt[i][j].pathLen, bt[i][j].GetDescription(), bt[i][j].count)
		}
		fmt.Println()
	}
}
