package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"granny-pass/internal/provider/graph"
)

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
			)

			t.Run("Empty knapsack & empty word", func(t *testing.T) {
				kNew, err = v.FindBestCombination(k, wm1)
				assert.Error(t, err)
			})

			t.Run("Empty knapsack & filled word", func(t *testing.T) {
				wm1 = &wordMetric{
					word:    "zas",
					pathLen: 2,
				}

				kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zas", kNew.GetDescription())
			})

			t.Run("Knapsack from test1.txt & empty word", func(t *testing.T) {
				wordMetrics, err = v.ReadFile("tests/test1.txt", false)
				assert.NoError(t, err)

				k = knapsack{
					items: wordMetrics,
				}

				kNew, err = v.FindBestCombination(k, wm2)
				assert.Error(t, err)
			})

			t.Run("Filled knapsack & filled word", func(t *testing.T) {
				kNew, err = v.FindBestCombination(k, wm1)
				assert.NoError(t, err)
				assert.Equal(t, "zasaofthebike", kNew.GetDescription())

				wm2 = &wordMetric{
					word:    "dew",
					pathLen: 2,
				}

				kNew, err = v.FindBestCombination(k, wm2)
				assert.NoError(t, err)
				assert.Equal(t, "aofthebikedew", kNew.GetDescription())
			})

		})

		t.Run("ChooseCandidate", func(t *testing.T) {

			k11 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}},
				pathLen: 0,
				count:   1,
			}
			k21 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}},
				pathLen: 5,
				count:   1,
			}
			k31 := knapsack{
				items:   []*wordMetric{{word: "is", pathLen: 6}},
				pathLen: 6,
				count:   1,
			}

			k12 := knapsack{
				items:   []*wordMetric{{word: "a", pathLen: 0}, {word: "of", pathLen: 5}},
				pathLen: 13,
				count:   2,
			}
			k22 := knapsack{
				items:   []*wordMetric{{word: "of", pathLen: 5}, {word: "to", pathLen: 4}},
				pathLen: 13,
				count:   2,
			}
			k32 := knapsack{
				items:   []*wordMetric{{word: "to", pathLen: 4}, {word: "of", pathLen: 5}},
				pathLen: 9,
				count:   2,
			}

			ks1 := make(map[uint8]knapsack)
			ks2 := make(map[uint8]knapsack)
			ks3 := make(map[uint8]knapsack)

			ks1[1] = k11
			ks1[2] = k12

			ks2[1] = k21
			ks3[2] = k22

			ks3[1] = k31
			ks3[2] = k32

			t.Run("1,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks1, ks2, ks3)
				assert.Equal(t, k21, ksRes[1])
				assert.Equal(t, k32, ksRes[2])
			})

			k41 := knapsack{
				items:   []*wordMetric{{word: "the", pathLen: 6}},
				pathLen: 6,
				count:   1,
			}
			ks4 := make(map[uint8]knapsack)
			ks4[1] = k41

			t.Run("4,2,3", func(t *testing.T) {
				ksRes := v.ChooseCandidate(ks4, ks2, ks3)
				assert.Equal(t, k41, ksRes[1])
				_, ok := ksRes[2]
				assert.Equal(t, false, ok)
			})

		})

		t.Run("", func(t *testing.T) {

		})
	})
}

func getDistanceMapForTests() map[string]int {
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
	dist := graph.BigramDistanceMap(m)
	return dist
}
