//go:build graphTest
// +build graphTest

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUndirected(t *testing.T) {

	t.Run("test undirected graph functions", func(t *testing.T) {
		var (
			v1, v2, v3, v4 = Vertex{Name: "a"}, Vertex{Name: "b"}, Vertex{Name: "c"}, Vertex{Name: "d"}
			r1, r2         Vertex
			edge           Edge[string]
			rEdge          Edge[Vertex]
			err            error
			order          int
		)

		hash := func(v Vertex) string {
			return v.Name
		}

		g := newUndirected(hash, newMemoryStorage[string]())

		t.Run("AddVertex", func(t *testing.T) {

			t.Run("add first", func(t *testing.T) {
				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 0, order)

				err = g.AddVertex(v1)
				assert.NoError(t, err)

				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 1, order)

				r1, err = g.Vertex(hash(v1))
				assert.NoError(t, err)
				assert.Equal(t, v1, r1)
			})

			t.Run("add second", func(t *testing.T) {
				err = g.AddVertex(v2)
				assert.NoError(t, err)

				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 2, order)

				r2, err = g.Vertex(hash(v2))
				assert.NoError(t, err)
				assert.Equal(t, v2, r2)
			})

			t.Run("check not added", func(t *testing.T) {
				r2, err = g.Vertex(hash(v3))
				assert.Error(t, err)

				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 2, order)
			})
		})

		t.Run("AddEdge", func(t *testing.T) {
			t.Run("first edge", func(t *testing.T) {
				edge = Edge[string]{
					v1: g.hash(v1),
					v2: g.hash(v2),
				}

				err = g.AddEdge(edge.v1, edge.v2)
				assert.NoError(t, err)

				t.Run("check direct", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v1), g.hash(v2))
					assert.NoError(t, err)
					assert.Equal(t, edge.v1, rEdge.v1.Name)
					assert.Equal(t, edge.v2, rEdge.v2.Name)
				})

				t.Run("check reverse", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v2), g.hash(v1))
					assert.NoError(t, err)
					assert.Equal(t, edge.v1, rEdge.v2.Name)
					assert.Equal(t, edge.v1, rEdge.v2.Name)
				})

				t.Run("check nonexistent", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v1), g.hash(v3))
					assert.Error(t, err)
				})

				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 2, order)
			})

			t.Run("second edge", func(t *testing.T) {
				err = g.AddVertex(v3)
				assert.NoError(t, err)

				edge = Edge[string]{
					v1: g.hash(v2),
					v2: g.hash(v3),
				}

				err = g.AddEdge(edge.v1, edge.v2)
				assert.NoError(t, err)

				t.Run("check direct", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v2), g.hash(v3))
					assert.NoError(t, err)
					assert.Equal(t, edge.v1, rEdge.v1.Name)
					assert.Equal(t, edge.v2, rEdge.v2.Name)
				})

				t.Run("check reverse", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v3), g.hash(v2))
					assert.NoError(t, err)
					assert.Equal(t, edge.v1, rEdge.v2.Name)
					assert.Equal(t, edge.v1, rEdge.v2.Name)
				})

				t.Run("check nonexistent", func(t *testing.T) {
					rEdge, err = g.Edge(g.hash(v1), g.hash(v3))
					assert.Error(t, err)
				})

				order, err = g.Order()
				assert.NoError(t, err)
				assert.Equal(t, 3, order)
			})
		})

		t.Run("AdjacencyMapWithMaxWeight", func(t *testing.T) {
			var (
				n  = 20
				m1 map[string]map[string]int
			)
			m1, err = g.AdjacencyMapWithMaxWeight(n)
			assert.NoError(t, err)
			assert.Equal(t, 0, m1[hash(v1)][hash(v1)])
			assert.Equal(t, 0, m1[hash(v2)][hash(v2)])
			assert.Equal(t, 0, m1[hash(v3)][hash(v3)])

			assert.Equal(t, 1, m1[hash(v1)][hash(v2)])
			assert.Equal(t, 1, m1[hash(v2)][hash(v1)])
			assert.Equal(t, 1, m1[hash(v2)][hash(v3)])
			assert.Equal(t, 1, m1[hash(v3)][hash(v2)])

			assert.Equal(t, n, m1[hash(v1)][hash(v3)])
			assert.Equal(t, n, m1[hash(v3)][hash(v1)])

			t.Run("check after adding unconnected vertex", func(t *testing.T) {
				err = g.AddVertex(v4)
				assert.NoError(t, err)

				m1, err = g.AdjacencyMapWithMaxWeight(n)
				assert.NoError(t, err)
				assert.Equal(t, 0, m1[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m1[hash(v2)][hash(v2)])
				assert.Equal(t, 0, m1[hash(v3)][hash(v3)])
				assert.Equal(t, 0, m1[hash(v4)][hash(v4)])

				assert.Equal(t, 1, m1[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m1[hash(v2)][hash(v1)])
				assert.Equal(t, 1, m1[hash(v2)][hash(v3)])
				assert.Equal(t, 1, m1[hash(v3)][hash(v2)])

				assert.Equal(t, n, m1[hash(v1)][hash(v3)])
				assert.Equal(t, n, m1[hash(v3)][hash(v1)])

				assert.Equal(t, n, m1[hash(v1)][hash(v4)])
				assert.Equal(t, n, m1[hash(v2)][hash(v4)])
				assert.Equal(t, n, m1[hash(v3)][hash(v4)])

				assert.Equal(t, n, m1[hash(v4)][hash(v1)])
				assert.Equal(t, n, m1[hash(v4)][hash(v2)])
				assert.Equal(t, n, m1[hash(v4)][hash(v3)])

			})

			t.Run("check after adding edge", func(t *testing.T) {
				edge = Edge[string]{
					v1: g.hash(v3),
					v2: g.hash(v4),
				}

				err = g.AddEdge(edge.v1, edge.v2)
				assert.NoError(t, err)

				m1, err = g.AdjacencyMapWithMaxWeight(n)
				assert.NoError(t, err)
				assert.Equal(t, 0, m1[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m1[hash(v2)][hash(v2)])
				assert.Equal(t, 0, m1[hash(v3)][hash(v3)])
				assert.Equal(t, 0, m1[hash(v4)][hash(v4)])

				assert.Equal(t, 1, m1[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m1[hash(v2)][hash(v1)])
				assert.Equal(t, 1, m1[hash(v2)][hash(v3)])
				assert.Equal(t, 1, m1[hash(v3)][hash(v2)])
				assert.Equal(t, 1, m1[hash(v3)][hash(v4)])
				assert.Equal(t, 1, m1[hash(v4)][hash(v3)])

				assert.Equal(t, n, m1[hash(v1)][hash(v3)])
				assert.Equal(t, n, m1[hash(v3)][hash(v1)])

				assert.Equal(t, n, m1[hash(v1)][hash(v4)])
				assert.Equal(t, n, m1[hash(v2)][hash(v4)])

				assert.Equal(t, n, m1[hash(v4)][hash(v1)])
				assert.Equal(t, n, m1[hash(v4)][hash(v2)])
			})
		})

		t.Run("test WFI for undirected graph", func(t *testing.T) {

			var (
				m          map[string]map[string]int
				maxN       = 20
				v1, v2, v3 = Vertex{Name: "aaa"}, Vertex{Name: "bbb"}, Vertex{Name: "ccc"}
				edge       Edge[string]
				err        error
			)

			hash := func(v Vertex) string {
				return v.Name
			}

			g := newUndirected(hash, newMemoryStorage[string]())

			t.Run("empty graph", func(t *testing.T) {
				m, err = g.WFI(maxN)
				assert.Error(t, err)
			})

			t.Run("1 vertex graph", func(t *testing.T) {
				err = g.AddVertex(v1)
				assert.NoError(t, err)

				m, err = g.WFI(maxN)
				assert.NoError(t, err)
				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
			})

			t.Run("2 vertex graph", func(t *testing.T) {
				err = g.AddVertex(v2)
				assert.NoError(t, err)

				m, err = g.WFI(maxN)
				assert.NoError(t, err)

				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m[hash(v2)][hash(v2)])

				assert.Equal(t, maxN, m[hash(v1)][hash(v2)])
				assert.Equal(t, maxN, m[hash(v2)][hash(v1)])
			})

			t.Run("2 vertex + 1 edge graph", func(t *testing.T) {
				edge = Edge[string]{
					v1: g.hash(v1),
					v2: g.hash(v2),
				}

				err = g.AddEdge(edge.v1, edge.v2)
				assert.NoError(t, err)

				m, err = g.WFI(maxN)
				assert.NoError(t, err)

				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m[hash(v2)][hash(v2)])

				assert.Equal(t, 1, m[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m[hash(v2)][hash(v1)])
			})

			t.Run("3 vertex + 1 edge graph", func(t *testing.T) {
				err = g.AddVertex(v3)
				assert.NoError(t, err)

				m, err = g.WFI(maxN)
				assert.NoError(t, err)

				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m[hash(v2)][hash(v2)])
				assert.Equal(t, 0, m[hash(v3)][hash(v3)])

				assert.Equal(t, 1, m[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m[hash(v2)][hash(v1)])

				assert.Equal(t, maxN, m[hash(v1)][hash(v3)])
				assert.Equal(t, maxN, m[hash(v3)][hash(v1)])
				assert.Equal(t, maxN, m[hash(v2)][hash(v3)])
				assert.Equal(t, maxN, m[hash(v3)][hash(v2)])
			})

			t.Run("3 vertex + 2 edge graph", func(t *testing.T) {
				edge = Edge[string]{
					v1: g.hash(v2),
					v2: g.hash(v3),
				}

				err = g.AddEdge(edge.v1, edge.v2)
				assert.NoError(t, err)

				m, err = g.WFI(maxN)
				assert.NoError(t, err)

				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m[hash(v2)][hash(v2)])
				assert.Equal(t, 0, m[hash(v3)][hash(v3)])

				assert.Equal(t, 1, m[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m[hash(v2)][hash(v1)])
				assert.Equal(t, 1, m[hash(v2)][hash(v3)])
				assert.Equal(t, 1, m[hash(v3)][hash(v2)])

				assert.Equal(t, 2, m[hash(v1)][hash(v3)])
				assert.Equal(t, 2, m[hash(v3)][hash(v1)])

			})

			t.Run("3 vertex + 3 edge graph", func(t *testing.T) {
				edge = Edge[string]{
					v1: g.hash(v2),
					v2: g.hash(v3),
				}

				m, err = g.WFI(maxN)
				assert.NoError(t, err)

				assert.Equal(t, 0, m[hash(v1)][hash(v1)])
				assert.Equal(t, 0, m[hash(v2)][hash(v2)])
				assert.Equal(t, 0, m[hash(v3)][hash(v3)])

				assert.Equal(t, 1, m[hash(v1)][hash(v2)])
				assert.Equal(t, 1, m[hash(v2)][hash(v1)])
				assert.Equal(t, 1, m[hash(v2)][hash(v3)])
				assert.Equal(t, 1, m[hash(v3)][hash(v2)])
				assert.Equal(t, 2, m[hash(v1)][hash(v3)])
				assert.Equal(t, 2, m[hash(v3)][hash(v1)])

			})
		})
	})
}
