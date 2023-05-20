package board_validation

type Union struct {
	parent []int
	rank   []int
	count  int
}

func NewUnion(size int) *Union {
	u := Union{
		parent: make([]int, size),
		rank:   make([]int, size),
		count:  0,
	}
	for i := 0; i < size; i++ {
		u.parent[i] = -1
		u.rank[i] = 0
	}
	return &u
}

func (u *Union) Find(x int) int {
	if u.parent[x] == -1 {
		return x
	}
	u.parent[x] = u.Find(u.parent[x])
	return u.parent[x]
}

func (u *Union) Union(x, y int) {
	xRoot := u.Find(x)
	yRoot := u.Find(y)

	if xRoot == yRoot {
		return
	}

	switch {
	case u.rank[xRoot] < u.rank[yRoot]:
		u.parent[xRoot] = yRoot
	case u.rank[xRoot] > u.rank[yRoot]:
		u.parent[yRoot] = xRoot
	default:
		u.parent[yRoot] = xRoot
		u.rank[xRoot]++
	}

	u.count--
}

func ConnectedComponentLabeling(board [][]int) [][]int {
	rows := len(board)
	cols := len(board[0])
	labels := make([][]int, rows)
	for i := range labels {
		labels[i] = make([]int, cols)
	}

	u := NewUnion(rows * cols)

	// First pass
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == 1 {
				if i > 0 && board[i-1][j] == 1 {
					u.Union(i*cols+j, (i-1)*cols+j)
				}
				if j > 0 && board[i][j-1] == 1 {
					u.Union(i*cols+j, i*cols+j-1)
				}
			}
		}
	}

	// Second pass
	label := 1
	labelMap := make(map[int]int)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == 1 {
				root := u.Find(i*cols + j)
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
