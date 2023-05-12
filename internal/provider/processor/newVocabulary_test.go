//go:build processorTest
// +build processorTest

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVocabulary(t *testing.T) {
	t.Run("test new vocabulary functions", func(t *testing.T) {
		var (
			w1, w2, w3, w4, w5 = "a", "of", "the", "cafe", "tanya"
			n                  int
			err                error
			wordMetrics        []*wordMetric
		)

		dist := getDistanceMapForTests()
		v := NewVocab(dist, 0, 0, 0)

		t.Run("PathLen", func(t *testing.T) {
			n, err = v.PathLen(w1)
			assert.NoError(t, err)
			assert.Equal(t, 0, n)

			n, err = v.PathLen(w2)
			assert.NoError(t, err)
			assert.Equal(t, 5, n)

			n, err = v.PathLen(w3)
			assert.NoError(t, err)
			assert.Equal(t, 6, n)

			n, err = v.PathLen(w4)
			assert.NoError(t, err)
			assert.Equal(t, 8, n)

			n, err = v.PathLen(w5)
			assert.NoError(t, err)
			assert.Equal(t, 17, n)

			//n, err = v.PathLen("q 0")
			//assert.Error(t, err)
			//
			//n, err = v.PathLen("12")
			//assert.Error(t, err)
			//
			//n, err = v.PathLen("?)")
			//assert.Error(t, err)
		})

		t.Run("GapPathLen", func(t *testing.T) {

			n, err = v.GapPathLen(w1, w2)
			assert.NoError(t, err)
			assert.Equal(t, 8, n)

			n, err = v.GapPathLen(w2, w3)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)

			n, err = v.GapPathLen(w3, w4)
			assert.NoError(t, err)
			assert.Equal(t, 2, n)

			n, err = v.GapPathLen(w4, w5)
			assert.NoError(t, err)
			assert.Equal(t, 2, n)

			n, err = v.GapPathLen("", "a")
			assert.NoError(t, err)
			assert.Equal(t, 0, n)

			n, err = v.GapPathLen("yes", "")
			assert.NoError(t, err)
			assert.Equal(t, 0, n)

			n, err = v.GapPathLen("q", "a")
			assert.NoError(t, err)
			assert.Equal(t, 1, n)

			n, err = v.GapPathLen("12", w1)
			assert.Error(t, err)

			n, err = v.GapPathLen(w1, "?)")
			assert.Error(t, err)
		})

		t.Run("ReadFile", func(t *testing.T) {
			wordMetrics, err = v.ReadFile("tests/test.txt", true)
			assert.NoError(t, err)

			length := wordMetrics[0].pathLen
			for _, wm := range wordMetrics {

				n, err = v.PathLen(wm.word)
				assert.NoError(t, err)
				assert.Equal(t, n, wm.pathLen)

				//check sorting
				assert.Equal(t, true, length <= len(wm.word))
			}
		})
	})
}
