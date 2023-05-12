//go:build graphTest
// +build graphTest

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {

	t.Run("test storage functions", func(t *testing.T) {
		storage := newMemoryStorage[string]()

		var (
			h1, h2, h3 = "a", "b", "c"
			v1, v2, v3 = Vertex{Name: "aaa"}, Vertex{Name: "bbb"}, Vertex{Name: "ccc"}

			r1, r2      Vertex
			edge, rEdge Edge[string]
			err         error
			cnt         int
			l1, l2      []string
		)

		t.Run("AddVertex", func(t *testing.T) {

			t.Run("add first", func(t *testing.T) {
				cnt, err = storage.VertexCount()
				assert.NoError(t, err)
				assert.Equal(t, 0, cnt)

				l1, err = storage.ListVertices()
				assert.NoError(t, err)
				assert.Equal(t, l2, l1)

				err = storage.AddVertex(h1, v1)
				assert.NoError(t, err)

				cnt, err = storage.VertexCount()
				assert.NoError(t, err)
				assert.Equal(t, 1, cnt)

				r1, err = storage.Vertex(h1)
				assert.NoError(t, err)
				assert.Equal(t, v1, r1)

				l1, err = storage.ListVertices()
				assert.NoError(t, err)
				l2 = append(l2, h1)
				assert.Equal(t, l2, l1)
			})

			t.Run("add second", func(t *testing.T) {
				err = storage.AddVertex(h2, v2)
				assert.NoError(t, err)

				cnt, err = storage.VertexCount()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)

				r2, err = storage.Vertex(h2)
				assert.NoError(t, err)
				assert.Equal(t, v2, r2)

				l1, err = storage.ListVertices()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(l1))
			})

			t.Run("nonexistence", func(t *testing.T) {
				r2, err = storage.Vertex(h3)
				assert.Error(t, err)

				cnt, err = storage.VertexCount()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
			})

			t.Run("add existent", func(t *testing.T) {
				err = storage.AddVertex(h2, v2)
				assert.Error(t, err)

				cnt, err = storage.VertexCount()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)

				r2, err = storage.Vertex(h2)
				assert.NoError(t, err)
				assert.Equal(t, v2, r2)
			})
		})

		t.Run("AddEdge", func(t *testing.T) {
			edge = Edge[string]{
				v1: h1,
				v2: h2,
			}
			err = storage.AddEdge(h1, h2, edge)
			assert.NoError(t, err)

			rEdge, err = storage.Edge(h1, h2)
			assert.NoError(t, err)
			assert.Equal(t, edge, rEdge)
		})

		t.Run("VertexCount", func(t *testing.T) {
			cnt, err = storage.VertexCount()
			assert.NoError(t, err)
			assert.Equal(t, 2, cnt)

			err = storage.AddVertex(h3, v3)
			assert.NoError(t, err)

			cnt, err = storage.VertexCount()
			assert.NoError(t, err)
			assert.Equal(t, 3, cnt)
		})
	})
}
