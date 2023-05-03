package graph

import (
	"errors"
	"fmt"
)

type undirected[K comparable, V Vertex] struct {
	hash    Hash[K, V]
	storage Storage[K, V]
}

func newUndirected[K comparable, V Vertex](hash Hash[K, V], storage Storage[K, V]) *undirected[K, V] {
	return &undirected[K, V]{
		hash:    hash,
		storage: storage,
	}
}

func (u *undirected[K, V]) AddVertex(value V) error {
	hash := u.hash(value)

	return u.storage.AddVertex(hash, value)
}

func (u *undirected[K, T]) Vertex(hash K) (T, error) {
	vertex, err := u.storage.Vertex(hash)
	return vertex, err
}

func (u *undirected[K, V]) AddEdge(source, target K) error {
	if _, err := u.storage.Vertex(source); err != nil {
		return fmt.Errorf("could not find source vertex with hash %v: %w", source, err)
	}

	if _, err := u.storage.Vertex(target); err != nil {
		return fmt.Errorf("could not find target vertex with hash %v: %w", target, err)
	}

	// nolint: govet // false positive err shawdowing
	if _, err := u.Edge(source, target); !errors.Is(err, ErrEdgeNotFound) {
		return ErrEdgeAlreadyExists
	}

	edge := Edge[K]{
		v1: source,
		v2: target,
	}

	if err := u.addEdge(source, target, edge); err != nil {
		return fmt.Errorf("failed to add edge: %w", err)
	}

	return nil
}

func (u *undirected[K, V]) Edge(source, target K) (Edge[V], error) {
	// In an undirected graph, since multigraphs aren't supported, the edge AB is the same as BA.
	// Therefore, if source[target] cannot be found, this function also looks for target[source].

	//edge, err := u.storage.Edge(source, target)
	_, err := u.storage.Edge(source, target)
	if errors.Is(err, ErrEdgeNotFound) {
		_, err = u.storage.Edge(target, source)
	}

	if err != nil {
		return Edge[V]{}, err
	}

	sourceVertex, err := u.storage.Vertex(source)
	if err != nil {
		return Edge[V]{}, err
	}

	targetVertex, err := u.storage.Vertex(target)
	if err != nil {
		return Edge[V]{}, err
	}

	return Edge[V]{
		v1: sourceVertex,
		v2: targetVertex,
	}, nil
}

func (u *undirected[K, V]) addEdge(sourceHash, targetHash K, edge Edge[K]) error {
	err := u.storage.AddEdge(sourceHash, targetHash, edge)
	if err != nil {
		return err
	}

	rEdge := Edge[K]{
		v1: edge.v1,
		v2: edge.v2,
	}

	err = u.storage.AddEdge(targetHash, sourceHash, rEdge)
	if err != nil {
		return err
	}

	return nil
}

func (u *undirected[K, V]) Order() (int, error) {
	return u.storage.VertexCount()
}

func (u *undirected[K, V]) AdjacencyMapWithMaxWeight(maxN int) (map[K]map[K]int, error) {
	vertices, err := u.storage.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	m := make(map[K]map[K]int)

	for _, vertex := range vertices {
		m[vertex] = make(map[K]int)
		for _, vertex2 := range vertices {
			// zero in the diagonal
			if vertex == vertex2 {
				m[vertex][vertex2] = 0
				continue
			}

			_, err := u.storage.Edge(vertex, vertex2)
			if err != nil {
				m[vertex][vertex2] = maxN
			} else {
				m[vertex][vertex2] = 1
			}
		}
	}
	return m, nil
}

func (u *undirected[K, V]) WFI(maxN int) (map[K]map[K]int, error) {
	dist, err := u.AdjacencyMapWithMaxWeight(maxN)
	if err != nil {
		return nil, err
	}

	vertices, err := u.storage.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	if len(vertices) == 0 {
		return nil, ErrNoVertices
	}

	for _, k := range vertices {
		for _, i := range vertices {
			for _, j := range vertices {
				if dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist, nil
}
