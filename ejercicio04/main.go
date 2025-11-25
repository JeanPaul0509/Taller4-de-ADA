package main

import (
	"bufio"
	"container/heap"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// --- Estructuras de Datos ---

type Node struct {
	ID  int64
	Lat float64
	Lon float64
}

type Edge struct {
	From     int64
	To       int64
	Distance float64
}

// Graph representa el grafo como una lista de adyacencia
type Graph struct {
	Nodes map[int64]Node
	Edges map[int64][]Edge
}

// --- Cola de Prioridad para Dijkstra ---

type Item struct {
	NodeID   int64
	Priority float64 // Distancia acumulada
	Index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// --- Funciones Auxiliares ---

// Calcular distancia haversine entre dos coordenadas
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Radio de la Tierra en metros

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// --- Algoritmo de Dijkstra ---

func Dijkstra(graph *Graph, startID, endID int64) (float64, []int64) {
	dist := make(map[int64]float64)
	prev := make(map[int64]int64)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Inicializar distancias
	dist[startID] = 0
	heap.Push(&pq, &Item{NodeID: startID, Priority: 0})

	found := false

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item)
		uID := u.NodeID

		if uID == endID {
			found = true
			break
		}

		// Si encontramos un camino más largo, lo ignoramos
		if d, ok := dist[uID]; ok && u.Priority > d {
			continue
		}

		// Explorar vecinos
		if neighbors, ok := graph.Edges[uID]; ok {
			for _, edge := range neighbors {
				vID := edge.To
				alt := dist[uID] + edge.Distance

				if d, ok := dist[vID]; !ok || alt < d {
					dist[vID] = alt
					prev[vID] = uID
					heap.Push(&pq, &Item{NodeID: vID, Priority: alt})
				}
			}
		}
	}

	if !found {
		return -1, nil
	}

	// Reconstruir camino
	path := []int64{}
	curr := endID
	for curr != startID {
		path = append([]int64{curr}, path...) // Prepend
		curr = prev[curr]
	}
	path = append([]int64{startID}, path...)

	return dist[endID], path
}

// --- Funciones de Carga y Exportación ---

func loadGraphFromCSV(nodesFile, edgesFile string) (*Graph, error) {
	graph := &Graph{
		Nodes: make(map[int64]Node),
		Edges: make(map[int64][]Edge),
	}

	// Cargar Nodos
	fNodes, err := os.Open(nodesFile)
	if err != nil {
		return nil, err
	}
	defer fNodes.Close()

	readerNodes := csv.NewReader(fNodes)
	recordsNodes, err := readerNodes.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range recordsNodes {
		if i == 0 {
			continue
		} // Skip header
		id, _ := strconv.ParseInt(record[0], 10, 64)
		lat, _ := strconv.ParseFloat(record[1], 64)
		lon, _ := strconv.ParseFloat(record[2], 64)
		graph.Nodes[id] = Node{ID: id, Lat: lat, Lon: lon}
	}

	// Cargar Aristas
	fEdges, err := os.Open(edgesFile)
	if err != nil {
		return nil, err
	}
	defer fEdges.Close()

	readerEdges := csv.NewReader(fEdges)
	recordsEdges, err := readerEdges.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range recordsEdges {
		if i == 0 {
			continue
		} // Skip header
		from, _ := strconv.ParseInt(record[0], 10, 64)
		to, _ := strconv.ParseInt(record[1], 10, 64)
		dist, _ := strconv.ParseFloat(record[2], 64)

		graph.Edges[from] = append(graph.Edges[from], Edge{From: from, To: to, Distance: dist})
	}

	return graph, nil
}

func extractOSMData(pbfFile string) {
	file, err := os.Open(pbfFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := osmpbf.New(context.Background(), file, 3)
	defer scanner.Close()

	nodes := make(map[int64]Node)
	var edges []Edge

	// Bounding box para Lima
	const (
		minLat = -12.3
		maxLat = -11.7
		minLon = -77.3
		maxLon = -76.8
	)

	fmt.Println("Procesando archivo OSM (esto puede tardar)...")

	for scanner.Scan() {
		switch obj := scanner.Object().(type) {
		case *osm.Node:
			if obj.Lat >= minLat && obj.Lat <= maxLat &&
				obj.Lon >= minLon && obj.Lon <= maxLon {
				nodes[int64(obj.ID)] = Node{
					ID:  int64(obj.ID),
					Lat: obj.Lat,
					Lon: obj.Lon,
				}
			}

		case *osm.Way:
			highway := obj.Tags.Find("highway")
			if highway != "" {
				validTypes := map[string]bool{
					"motorway": true, "trunk": true, "primary": true,
					"secondary": true, "tertiary": true, "residential": true,
					"unclassified": true, "service": true,
				}

				if validTypes[highway] {
					for i := 0; i < len(obj.Nodes)-1; i++ {
						fromID := int64(obj.Nodes[i].ID)
						toID := int64(obj.Nodes[i+1].ID)

						from, fromOk := nodes[fromID]
						to, toOk := nodes[toID]

						if fromOk && toOk {
							distance := haversine(from.Lat, from.Lon, to.Lat, to.Lon)
							edges = append(edges, Edge{From: fromID, To: toID, Distance: distance})

							oneway := obj.Tags.Find("oneway")
							if oneway != "yes" {
								edges = append(edges, Edge{From: toID, To: fromID, Distance: distance})
							}
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Nodos extraídos: %d\n", len(nodes))
	fmt.Printf("✅ Aristas extraídas: %d\n", len(edges))

	exportNodesToCSV(nodes, "lima_nodes.csv")
	exportEdgesToCSV(edges, "lima_edges.csv")
}

func exportNodesToCSV(nodes map[int64]Node, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"id", "lat", "lon"})
	for _, node := range nodes {
		writer.Write([]string{fmt.Sprintf("%d", node.ID), fmt.Sprintf("%.6f", node.Lat), fmt.Sprintf("%.6f", node.Lon)})
	}
}

func exportEdgesToCSV(edges []Edge, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"from", "to", "distance"})
	for _, edge := range edges {
		writer.Write([]string{fmt.Sprintf("%d", edge.From), fmt.Sprintf("%d", edge.To), fmt.Sprintf("%.2f", edge.Distance)})
	}
}

func main() {
	// 1. Verificar si existen los CSVs, si no, extraerlos
	if _, err := os.Stat("lima_nodes.csv"); os.IsNotExist(err) {
		fmt.Println("⚠️ Archivos CSV no encontrados. Iniciando extracción desde PBF...")
		extractOSMData("peru-251122.osm.pbf")
	} else {
		fmt.Println("✅ Archivos CSV encontrados. Saltando extracción.")
	}

	// 2. Cargar Grafo
	fmt.Println("Cargando grafo en memoria...")
	graph, err := loadGraphFromCSV("lima_nodes.csv", "lima_edges.csv")
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Grafo cargado: %d nodos\n", len(graph.Nodes))

	// 3. Interacción con el usuario
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Navegación GPS (Dijkstra) ---")
		fmt.Print("Ingrese ID Nodo Origen (o 'q' para salir): ")
		inputStart, _ := reader.ReadString('\n')
		inputStart = strings.TrimSpace(inputStart)
		if inputStart == "q" {
			break
		}

		startID, err := strconv.ParseInt(inputStart, 10, 64)
		if err != nil {
			fmt.Println("❌ ID inválido.")
			continue
		}

		fmt.Print("Ingrese ID Nodo Destino: ")
		inputEnd, _ := reader.ReadString('\n')
		inputEnd = strings.TrimSpace(inputEnd)
		endID, err := strconv.ParseInt(inputEnd, 10, 64)
		if err != nil {
			fmt.Println("❌ ID inválido.")
			continue
		}

		// Verificar existencia
		if _, ok := graph.Nodes[startID]; !ok {
			fmt.Printf("❌ Nodo origen %d no existe en el grafo.\n", startID)
			continue
		}
		if _, ok := graph.Nodes[endID]; !ok {
			fmt.Printf("❌ Nodo destino %d no existe en el grafo.\n", endID)
			continue
		}

		// 4. Ejecutar Dijkstra
		fmt.Printf("Calculando ruta de %d a %d...\n", startID, endID)
		dist, path := Dijkstra(graph, startID, endID)

		if dist == -1 {
			fmt.Println("❌ No existe ruta entre estos nodos.")
		} else {
			fmt.Printf("✅ Ruta encontrada!\n")
			fmt.Printf("📏 Distancia total: %.2f metros\n", dist)
			fmt.Printf("🛣️  Pasos (%d nodos): %v\n", len(path), path)
		}
	}
}
