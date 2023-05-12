package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {

	t.Run("test file functions", func(t *testing.T) {
		var (
			dist, distRes []int
			m             map[string]map[string]int
			err           error
			filename      = "testdata/test1.json"
			s             string
			bigrams       []string
		)

		t.Run("BigramDistanceMap", func(t *testing.T) {
			m = getDistMapForTest()
			dist = BigramDistanceArray(m)

			bigrams = []string{"aa", "ss", "dd", "as", "sa", "sd", "ds", "ad", "da"}
			for _, s = range bigrams {
				assert.Equal(t, m[string(s[0])][string(s[1])], dist[getIndex(s[0], s[1])])
				if s[0] == s[1] {
					assert.Equal(t, true, dist[getIndex(s[0], s[1])] == 0)
				} else if s == "ad" || s == "da" {
					assert.Equal(t, true, dist[getIndex(s[0], s[1])] == 2)
				} else {
					assert.Equal(t, true, dist[getIndex(s[0], s[1])] == 1)
				}
			}
		})

		t.Run("SaveToJson", func(t *testing.T) {
			err = SaveToJson(dist, filename)
			assert.NoError(t, err)
		})

		t.Run("ReadFromJson", func(t *testing.T) {
			distRes, err = ReadFromJson(filename)
			assert.NoError(t, err)

			for k := range dist {
				assert.Equal(t, dist[k], distRes[k])
			}
		})
	})
}

func getDistMapForTest() map[string]map[string]int {
	var (
		edge       Edge[string]
		v1, v2, v3 = Vertex{Name: "a"}, Vertex{Name: "s"}, Vertex{Name: "d"}
		maxN       = 20
	)

	hash := func(v Vertex) string {
		return v.Name
	}
	g := newUndirected(hash, newMemoryStorage[string]())

	_ = g.AddVertex(v1)
	_ = g.AddVertex(v2)
	_ = g.AddVertex(v3)

	edge = Edge[string]{
		v1: g.hash(v1),
		v2: g.hash(v2),
	}
	_ = g.AddEdge(edge.v1, edge.v2)

	edge = Edge[string]{
		v1: g.hash(v2),
		v2: g.hash(v3),
	}

	_ = g.AddEdge(edge.v1, edge.v2)

	m, _ := g.WFI(maxN)

	return m
}
