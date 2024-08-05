package comparator

func NthElement(arr *[]interface{}, n int, Cmp ...Comparator) (value interface{}) {
	if arr == nil || (*arr) == nil || len((*arr)) == 0 {
		return nil
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = GetCmp((*arr)[0])
	}
	if cmp == nil {
		return nil
	}
	//判断待确认的第n位是否在该集合范围内
	if len((*arr)) < n || n < 0 {
		return nil
	}
	//进行查找
	nthElement(arr, 0, len((*arr))-1, n, cmp)
	return (*arr)[n]
}

func nthElement(arr *[]interface{}, l, r int, n int, cmp Comparator) {
	//二分该区域并对此进行预排序
	if l >= r {
		return
	}
	m := (*arr)[(r+l)/2]
	i, j := l-1, r+1
	for i < j {
		i++
		for cmp((*arr)[i], m) < 0 {
			i++
		}
		j--
		for cmp((*arr)[j], m) > 0 {
			j--
		}
		if i < j {
			(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
		}
	}
	//确认第n位的范围进行局部二分
	if n-1 >= i {
		nthElement(arr, j+1, r, n, cmp)
	} else {
		nthElement(arr, l, j, n, cmp)
	}
}
