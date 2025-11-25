# Ejercicio 04: Navegación GPS

## Enfoque
Se utiliza el algoritmo de Dijkstra para encontrar la ruta más corta entre dos puntos en un mapa real extraído de OpenStreetMap. Se construye un grafo donde las intersecciones son nodos y las calles son aristas ponderadas por la distancia.

## Complejidad
- **Temporal:** O((V + E) log V) utilizando una cola de prioridad (Heap).
- **Espacial:** O(V + E) para almacenar el grafo y las distancias.

## Ejecución
Asegúrate de estar en la carpeta `ejercicio04`:

```bash
go run main.go
```

## Casos de prueba incluidos
- **TestDijkstra**: Verifica que el algoritmo encuentre la ruta más corta en un grafo dirigido simple, probando caminos directos vs indirectos más cortos.
- **TestHaversine**: Verifica que el cálculo de distancia entre dos coordenadas geográficas (Lima - Cusco) sea aproximado al valor real.
