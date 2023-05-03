package graph

import "sync"

type Storage[K comparable, V Vertex] interface {
	AddVertex(hash K, value V) error
	Vertex(hash K) (V, error)
	AddEdge(source, target K, edge Edge[K]) error
	Edge(source, target K) (Edge[K], error)
	VertexCount() (int, error)
	ListVertices() ([]K, error)
}

type memoryStorage[K comparable, V Vertex] struct {
	lock     sync.RWMutex
	vertices map[K]V
	//vertexProperties map[K]VertexProperties
	//edges map[K]map[K]Edge[K]
	// outEdges and inEdges store all outgoing and ingoing edges for all vertices. For O(1) access,
	// these edges themselves are stored in maps whose keys are the hashes of the target vertices.
	outEdges map[K]map[K]Edge[K] // source -> target
	inEdges  map[K]map[K]Edge[K] // target -> source
}

func newMemoryStorage[K comparable, V Vertex]() Storage[K, V] {
	return &memoryStorage[K, V]{
		vertices: make(map[K]V),
		outEdges: make(map[K]map[K]Edge[K]),
		inEdges:  make(map[K]map[K]Edge[K]),
	}
}

func (s *memoryStorage[K, T]) AddVertex(k K, t T) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.vertices[k]; ok {
		return ErrVertexAlreadyExists
	}

	s.vertices[k] = t

	return nil
}

func (s *memoryStorage[K, T]) Vertex(k K) (T, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var v T
	var ok bool
	v, ok = s.vertices[k]
	if !ok {
		return v, ErrVertexNotFound
	}

	return v, nil
}

func (s *memoryStorage[K, T]) AddEdge(sourceHash, targetHash K, edge Edge[K]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.outEdges[sourceHash]; !ok {
		s.outEdges[sourceHash] = make(map[K]Edge[K])
	}

	s.outEdges[sourceHash][targetHash] = edge

	if _, ok := s.inEdges[targetHash]; !ok {
		s.inEdges[targetHash] = make(map[K]Edge[K])
	}

	s.inEdges[targetHash][sourceHash] = edge

	return nil
}

func (s *memoryStorage[K, T]) Edge(sourceHash, targetHash K) (Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return Edge[K]{}, ErrEdgeNotFound
	}

	edge, ok := sourceEdges[targetHash]
	if !ok {
		return Edge[K]{}, ErrEdgeNotFound
	}

	return edge, nil
}

func (s *memoryStorage[K, T]) VertexCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.vertices), nil
}

func (s *memoryStorage[K, T]) ListVertices() ([]K, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var hashes []K
	for k := range s.vertices {
		hashes = append(hashes, k)
	}

	return hashes, nil
}
