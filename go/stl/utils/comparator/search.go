package comparator

func Search(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr == nil || (*arr) == nil || len((*arr)) == 0 {
		return
	}
	//判断比较函数是否有效,若无效则寻找默认比较函数
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		//若并非默认类型且未传入比较器则直接结束
		return -1
	}
	//查找开始
	return search(arr, e, cmp)
}

func search(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	//通过二分查找的方式寻找该元素
	l, m, r := 0, (len((*arr))-1)/2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) < 0 {
			l = m + 1
		} else {
			r = m
		}
	}
	//查找结束
	if l < len(*arr) && (*arr)[l] == e {
		//该元素存在,返回下标
		return l
	}
	//该元素不存在,返回-1
	return -1
}
