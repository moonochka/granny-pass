//go:build processorTest
// +build processorTest

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"granny-pass/internal/provider/graph"
)

func TestNewVocabulary(t *testing.T) {
	t.Run("test new vocabulary functions", func(t *testing.T) {
		var (
			w1, w2, w3, w4, w5, w6, w7 = "a", "of", "the", "cafe", "tanya", "az", "za"
			n                          int
			err                        error
			wordMetrics                []*wordMetric
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

			n, err = v.PathLen(w6)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)

			n, err = v.PathLen(w7)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)

			n, err = v.PathLen("q 0")
			assert.Error(t, err)

			n, err = v.PathLen("12")
			assert.Error(t, err)

			n, err = v.PathLen("?)")
			assert.Error(t, err)
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

		})

		t.Run("ReadFile", func(t *testing.T) {
			wordMetrics, err = v.ReadFile("testdata/test.txt", true)
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

func getDistanceMapForTests() []int {
	hash := func(v graph.Vertex) string {
		return v.Name
	}
	g := graph.New(hash)

	//add all key buttons
	for r := 'a'; r <= 'z'; r++ {
		_ = g.AddVertex(graph.Vertex{Name: string(r)})
	}

	//add all connections weight=1
	_ = g.AddEdge("q", "w")
	_ = g.AddEdge("w", "e")
	_ = g.AddEdge("e", "r")
	_ = g.AddEdge("r", "t")
	_ = g.AddEdge("t", "y")
	_ = g.AddEdge("y", "u")
	_ = g.AddEdge("u", "i")
	_ = g.AddEdge("i", "o")
	_ = g.AddEdge("o", "p")

	_ = g.AddEdge("a", "s")
	_ = g.AddEdge("s", "d")
	_ = g.AddEdge("d", "f")
	_ = g.AddEdge("f", "g")
	_ = g.AddEdge("g", "h")
	_ = g.AddEdge("h", "j")
	_ = g.AddEdge("j", "k")
	_ = g.AddEdge("k", "l")

	_ = g.AddEdge("z", "x")
	_ = g.AddEdge("x", "c")
	_ = g.AddEdge("c", "v")
	_ = g.AddEdge("v", "b")
	_ = g.AddEdge("b", "n")
	_ = g.AddEdge("n", "m")

	_ = g.AddEdge("q", "a")
	_ = g.AddEdge("w", "a")
	_ = g.AddEdge("w", "s")
	_ = g.AddEdge("e", "s")
	_ = g.AddEdge("e", "d")
	_ = g.AddEdge("r", "d")
	_ = g.AddEdge("r", "f")
	_ = g.AddEdge("t", "f")
	_ = g.AddEdge("t", "g")
	_ = g.AddEdge("y", "g")
	_ = g.AddEdge("y", "h")
	_ = g.AddEdge("u", "h")
	_ = g.AddEdge("u", "j")
	_ = g.AddEdge("i", "j")
	_ = g.AddEdge("i", "k")
	_ = g.AddEdge("o", "k")
	_ = g.AddEdge("o", "l")
	_ = g.AddEdge("p", "l")

	_ = g.AddEdge("a", "z")
	_ = g.AddEdge("s", "z")
	_ = g.AddEdge("s", "x")
	_ = g.AddEdge("d", "x")
	_ = g.AddEdge("d", "c")
	_ = g.AddEdge("f", "c")
	_ = g.AddEdge("f", "v")
	_ = g.AddEdge("g", "v")
	_ = g.AddEdge("g", "b")
	_ = g.AddEdge("h", "b")
	_ = g.AddEdge("h", "n")
	_ = g.AddEdge("j", "n")
	_ = g.AddEdge("j", "m")
	_ = g.AddEdge("k", "m")

	m, _ := g.WFI(20)
	dist := graph.BigramDistanceArray(m)
	return dist
}
