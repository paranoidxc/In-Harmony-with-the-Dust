package hash_map

type node struct {
	value interface{} //节点中存储的元素
	num   int         //该元素数量
	depth int         //该节点的深度
	left  *node       //左节点指针
	right *node       //右节点指针
}
