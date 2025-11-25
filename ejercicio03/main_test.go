package main

import (
	"testing"
)

func TestMST(t *testing.T) {
	// Grafo simple
	// 1 --(1)--> 2
	// 2 --(2)--> 3
	// 1 --(10)--> 3
	// MST debería usar (1,2) y (2,3) con costo total 3.
	edges := []Edge{
		{u: 1, v: 2, w: 1},
		{u: 2, v: 3, w: 2},
		{u: 1, v: 3, w: 10},
	}
	n := 3

	// Test Prim
	costPrim, _ := PrimMST(n, edges)
	if costPrim != 3 {
		t.Errorf("Prim MST cost expected 3, got %d", costPrim)
	}

	// Test Kruskal
	costKruskal, _ := KruskalMST(n, edges)
	if costKruskal != 3 {
		t.Errorf("Kruskal MST cost expected 3, got %d", costKruskal)
	}
}

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(5)

	if uf.Find(1) != 1 {
		t.Errorf("Initial parent should be self")
	}

	uf.Union(1, 2)
	if uf.Find(1) != uf.Find(2) {
		t.Errorf("1 and 2 should be connected")
	}

	uf.Union(3, 4)
	uf.Union(2, 4)

	if uf.Find(1) != uf.Find(3) {
		t.Errorf("1 and 3 should be connected transitively")
	}
}
