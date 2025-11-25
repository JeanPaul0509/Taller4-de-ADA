package main

import (
	"testing"
)

func TestDijkstra(t *testing.T) {
	// Crear un grafo simple
	graph := &Graph{
		Nodes: map[int64]Node{
			1: {ID: 1, Lat: 0, Lon: 0},
			2: {ID: 2, Lat: 0, Lon: 0},
			3: {ID: 3, Lat: 0, Lon: 0},
		},
		Edges: map[int64][]Edge{
			1: {{From: 1, To: 2, Distance: 10}, {From: 1, To: 3, Distance: 5}},
			2: {{From: 2, To: 3, Distance: 2}}, // 1 -> 2 -> 3 = 12, 1 -> 3 = 5. Shortest 1->3 is 5.
			3: {},
		},
	}

	// Test 1 -> 3 (Directo es mejor)
	dist, path := Dijkstra(graph, 1, 3)
	if dist != 5 {
		t.Errorf("Expected distance 5, got %f", dist)
	}
	if len(path) != 2 || path[0] != 1 || path[1] != 3 {
		t.Errorf("Expected path [1 3], got %v", path)
	}

	// Modificar grafo para que 1 -> 2 -> 3 sea mejor
	graph.Edges[1] = []Edge{{From: 1, To: 2, Distance: 2}, {From: 1, To: 3, Distance: 10}}
	// 1 -> 2 (2) + 2 -> 3 (2) = 4. Directo 1 -> 3 (10).

	dist, path = Dijkstra(graph, 1, 3)
	if dist != 4 {
		t.Errorf("Expected distance 4, got %f", dist)
	}
	if len(path) != 3 || path[0] != 1 || path[1] != 2 || path[2] != 3 {
		t.Errorf("Expected path [1 2 3], got %v", path)
	}
}

func TestHaversine(t *testing.T) {
	// Distancia aproximada entre Lima y Cusco
	// Lima: -12.0464, -77.0428
	// Cusco: -13.5320, -71.9675
	dist := haversine(-12.0464, -77.0428, -13.5320, -71.9675)

	// Valor esperado aprox 570-580 km
	if dist < 570000 || dist > 580000 {
		t.Errorf("Haversine calculation seems off, got %f meters", dist)
	}
}
