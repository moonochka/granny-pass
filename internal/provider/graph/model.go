package graph

import "errors"

var (
	ErrVertexNotFound      = errors.New("vertex not found")
	ErrVertexAlreadyExists = errors.New("vertex already exists")
	ErrEdgeNotFound        = errors.New("edge not found")
	ErrEdgeAlreadyExists   = errors.New("edge already exists")
	ErrNoVertices          = errors.New("no vertices")
	//ErrEdgeCreatesCycle    = errors.New("edge would create a cycle")
)

type Vertex struct {
	Name string
}

type Hash[K comparable, V Vertex] func(V) K

// why not Vertex? becouse ut connects hashes of Vertices
type Edge[V comparable] struct {
	v1 V
	v2 V
}

type Graph[K comparable, V Vertex] interface {
	AddVertex(value V) error
	Vertex(hash K) (V, error)
	AddEdge(source, target K) error
	Edge(source, target K) (Edge[V], error)
	// Order returns the number of vertices in the graph.
	Order() (int, error)
	// WFI
	WFI(nMax int) (map[K]map[K]int, error)
	// AdjacencyMapWithMaxWeight
	AdjacencyMapWithMaxWeight(nMax int) (map[K]map[K]int, error)
}

func New[K comparable, V Vertex](hash Hash[K, V]) Graph[K, V] {
	return NewWithStorage(hash, newMemoryStorage[K, V]())
}

func NewWithStorage[K comparable, V Vertex](hash Hash[K, V], storage Storage[K, V]) Graph[K, V] {
	return newUndirected(hash, storage)
}
