package app

type UnionFind struct {
	parent []int
	rank   []int
	count  int
}

func NewUnionFind(size int) *UnionFind {
	uf := UnionFind{
		parent: make([]int, size),
		rank:   make([]int, size),
		count:  0,
	}
	for i := 0; i < size; i++ {
		uf.parent[i] = -1
		uf.rank[i] = 0
	}
	return &uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] == -1 {
		return x
	}
	uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
	xRoot := uf.Find(x)
	yRoot := uf.Find(y)
	if xRoot != yRoot {
		if uf.rank[xRoot] < uf.rank[yRoot] {
			uf.parent[xRoot] = yRoot
		} else if uf.rank[xRoot] > uf.rank[yRoot] {
			uf.parent[yRoot] = xRoot
		} else {
			uf.parent[yRoot] = xRoot
			uf.rank[xRoot]++
		}
		uf.count--
	}
}

func ConnectedComponentLabeling(image [][]int) [][]int {
	rows := len(image)
	cols := len(image[0])
	labels := make([][]int, rows)
	for i := range labels {
		labels[i] = make([]int, cols)
	}

	uf := NewUnionFind(rows * cols)

	// First pass
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if image[i][j] == 1 {
				if i > 0 && image[i-1][j] == 1 {
					uf.Union(i*cols+j, (i-1)*cols+j)
				}
				if j > 0 && image[i][j-1] == 1 {
					uf.Union(i*cols+j, i*cols+j-1)
				}
			}
		}
	}

	// Second pass
	label := 1
	labelMap := make(map[int]int)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if image[i][j] == 1 {
				root := uf.Find(i*cols + j)
				if _, ok := labelMap[root]; !ok {
					labelMap[root] = label
					label++
				}
				labels[i][j] = labelMap[root]
			}
		}
	}

	return labels
}
