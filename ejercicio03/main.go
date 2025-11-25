package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Edge struct {
	u, v, w int
}

type Item struct {
	node  int
	cost  int
	from  int
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type UnionFind struct {
	parent []int
	rank   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank}
}

func (uf *UnionFind) Find(i int) int {
	if uf.parent[i] != i {
		uf.parent[i] = uf.Find(uf.parent[i])
	}
	return uf.parent[i]
}

func (uf *UnionFind) Union(i, j int) bool {
	rootI := uf.Find(i)
	rootJ := uf.Find(j)

	if rootI != rootJ {
		if uf.rank[rootI] < uf.rank[rootJ] {
			uf.parent[rootI] = rootJ
		} else if uf.rank[rootI] > uf.rank[rootJ] {
			uf.parent[rootJ] = rootI
		} else {
			uf.parent[rootJ] = rootI
			uf.rank[rootI]++
		}
		return true
	}
	return false
}

func loadGraph(filepath string) (int, []Edge, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var edges []Edge
	nodes := make(map[int]bool)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if strings.HasPrefix(parts[0], "%") || strings.HasPrefix(parts[0], "#") {
			continue
		}

		if len(parts) >= 2 {
			u, _ := strconv.Atoi(parts[0])
			v, _ := strconv.Atoi(parts[1])
			w := 1
			if len(parts) >= 3 {
				w, _ = strconv.Atoi(parts[2])
			}
			edges = append(edges, Edge{u, v, w})
			nodes[u] = true
			nodes[v] = true
		}
	}

	return len(nodes), edges, scanner.Err()
}

func PrimMST(n int, edges []Edge) (int, []Edge) {
	adj := make(map[int][]Edge)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], Edge{e.u, e.v, e.w})
		adj[e.v] = append(adj[e.v], Edge{e.v, e.u, e.w})
	}

	mstCost := 0
	var mstEdges []Edge
	visited := make(map[int]bool)

	startNode := edges[0].u

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{node: startNode, cost: 0, from: -1})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		u := item.node
		cost := item.cost
		from := item.from

		if visited[u] {
			continue
		}

		visited[u] = true
		mstCost += cost
		if from != -1 {
			mstEdges = append(mstEdges, Edge{from, u, cost})
		}

		for _, e := range adj[u] {
			v := e.v
			if !visited[v] {
				heap.Push(&pq, &Item{node: v, cost: e.w, from: u})
			}
		}
	}

	return mstCost, mstEdges
}

func KruskalMST(n int, edges []Edge) (int, []Edge) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	uf := NewUnionFind(n)

	maxID := 0
	for _, e := range edges {
		if e.u > maxID {
			maxID = e.u
		}
		if e.v > maxID {
			maxID = e.v
		}
	}
	uf = NewUnionFind(maxID)

	mstCost := 0
	var mstEdges []Edge

	for _, e := range edges {
		if uf.Union(e.u, e.v) {
			mstCost += e.w
			mstEdges = append(mstEdges, e)
		}
	}

	return mstCost, mstEdges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <ruta_del_dataset>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	n, edges, err := loadGraph(filepath)
	if err != nil {
		fmt.Printf("Error cargando el grafo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Grafo cargado con %d nodos y %d aristas.\n", n, len(edges))

	totalPossibleCost := 0
	for _, e := range edges {
		totalPossibleCost += e.w
	}

	fmt.Println("\n--- Algoritmo de Prim ---")
	primCost, primEdges := PrimMST(n, edges)
	fmt.Printf("Costo Mínimo: %d\n", primCost)
	fmt.Printf("Conexiones (%d)\n", len(primEdges))

	fmt.Println("\n--- Algoritmo de Kruskal (Union-Find) ---")

	edgesCopy := make([]Edge, len(edges))
	copy(edgesCopy, edges)
	kruskalCost, _ := KruskalMST(n, edgesCopy)
	fmt.Printf("Costo Mínimo: %d\n", kruskalCost)

	fmt.Printf("\nCosto total si se construyeran todas las conexiones: %d\n", totalPossibleCost)

	if primCost == kruskalCost {
		fmt.Println("\n[EXITO] Ambos algoritmos retornaron el mismo costo mínimo.")
	} else {
		fmt.Println("\n[ADVERTENCIA] ¡Los algoritmos retornaron costos diferentes!")
	}
}
