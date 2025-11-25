# Ejercicio 03: Optimal Power Grid

## Enfoque
Se implementan los algoritmos de **Prim** y **Kruskal** para resolver el problema de conectar una red eléctrica (Power Grid) con el costo mínimo. Se utiliza un dataset real (`power-US-Grid.mtx`) para construir un grafo no dirigido donde los edificios son nodos y las líneas eléctricas son aristas ponderadas. El objetivo es encontrar el Árbol de Expansión Mínima (MST) que conecta todos los nodos sin ciclos y con el menor peso total.

## Complejidad
- **Temporal:** O(E log V) o O(E log E) dependiendo de la implementación (Prim con Heap o Kruskal con Union-Find).
- **Espacial:** O(V + E) para almacenar el grafo y estructuras auxiliares.

## Ejecución
```bash
go run main.go power-US-Grid.mtx
```

## Casos de prueba incluidos
- **TestMST**: Verifica que tanto el algoritmo de Prim como el de Kruskal calculen correctamente el costo mínimo en un grafo simple conocido.
- **TestUnionFind**: Verifica la correcta funcionalidad de la estructura de datos Union-Find (Find y Union).

## Dataset Citation
The dataset used in this exercise is from the Network Data Repository.

> Ryan A. Rossi and Nesreen K. Ahmed
> The Network Data Repository with Interactive Graph Analytics and Visualization
> Proceedings of the Twenty-Ninth AAAI Conference on Artificial Intelligence
> [http://networkrepository.com](http://networkrepository.com)
> 2015
