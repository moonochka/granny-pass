//go:build processorTest
// +build processorTest

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnapsack(t *testing.T) {
	t.Run("test knapsack functions", func(t *testing.T) {
		var (
			pathLen             int
			err                 error
			wordMetrics, wordMs []*wordMetric
			b, b1               knapsack
		)

		dist := getDistanceMapForTests()
		v := NewVocab(dist, 0, 0, 0)

		t.Run("test.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("testdata/test.txt", false)
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

			t.Run("firstWord & lastWord & Length", func(t *testing.T) {

				t.Run("empty knapsack", func(t *testing.T) {
					assert.Equal(t, "", b1.firstWord())
					assert.Equal(t, "", b1.lastWord())
					assert.Equal(t, 0, b1.Length())
					assert.Equal(t, true, b1.isEmpty())
				})

				t.Run("1 word knapsack", func(t *testing.T) {
					word := "hello"
					pathLen, err = v.PathLen(word)
					wm := wordMetric{
						word:    word,
						pathLen: pathLen,
					}
					wordMs = append(wordMs, &wm)
					b1 = knapsack{
						items: wordMs,
					}
					assert.Equal(t, word, b1.firstWord())
					assert.Equal(t, word, b1.lastWord())
					assert.Equal(t, 5, b1.Length())
					assert.Equal(t, false, b1.isEmpty())
				})

				t.Run("knapsack from test.txt", func(t *testing.T) {
					b = knapsack{
						items: wordMetrics,
					}

					assert.Equal(t, "a", b.firstWord())
					assert.Equal(t, "biotransformation", b.lastWord())
					assert.Equal(t, 28, b.Length())
					assert.Equal(t, false, b1.isEmpty())
				})
			})
		})

		t.Run("test1.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("testdata/test1.txt", false)
			assert.NoError(t, err)

			t.Run("knapsack functions", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "aofthebike", b.GetDescription())
				assert.Equal(t, 10, b.Length())
				assert.Equal(t, false, b1.isEmpty())
				assert.Equal(t, "a", b.firstWord())
				assert.Equal(t, "bike", b.lastWord())
			})

		})

		t.Run("test2.txt", func(t *testing.T) {

			wordMetrics, err = v.ReadFile("testdata/test2.txt", false)
			assert.NoError(t, err)

			t.Run("knapsack functions", func(t *testing.T) {
				b = knapsack{
					items: wordMetrics,
				}

				assert.Equal(t, "aofbiotransformationbikeoatmeals", b.GetDescription())
				assert.Equal(t, 32, b.Length())
				assert.Equal(t, false, b1.isEmpty())
				assert.Equal(t, "a", b.firstWord())
				assert.Equal(t, "oatmeals", b.lastWord())
			})
		})
	})
}
