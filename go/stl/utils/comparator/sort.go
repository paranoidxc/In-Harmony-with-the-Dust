package comparator

func Sort(arr *[]interface{}, Cmp ...Comparator) {
	//如果传入一个空数组或nil,则直接结束
	if arr == nil || (*arr) == nil || len((*arr)) == 0 {
		return
	}
	var cmp Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = GetCmp((*arr)[0])
	}
	if cmp == nil {
		//未传入比较器且并非默认类型导致未找到默认比较器则直接终止排序
		return
	}
	//根据数组长度进行分类选择排序函数
	if len((*arr)) < 2^26 {
		//当长度小于2^16时使用二分排序
		binary(arr, 0, len((*arr))-1, cmp)
	} else {
		merge(arr, 0, len((*arr))-1, cmp)
	}
}

// 二分排序
// 对传入的待比较数组中的元素使用比较函数进行二分排序
func binary(arr *[]interface{}, l, r int, cmp Comparator) {
	//对当前部分进行预排序,使得两侧都大于或小于中间值
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
	//对分好的两侧进行迭代二分排序
	binary(arr, l, j, cmp)
	binary(arr, j+1, r, cmp)
}

// 归并排序
// 对传入的两个迭代器中的内容使用比较器进行归并排序
func merge(arr *[]interface{}, l, r int, cmp Comparator) {
	//对当前部分进行分组排序,将该部分近似平均的拆为两部分进行比较排序
	if l >= r {
		return
	}
	m := (r + l) / 2
	//对待排序内容进行二分
	merge(arr, l, m, cmp)
	merge(arr, m+1, r, cmp)
	//二分结束后依次比较进行归并
	i, j := l, m+1
	var tmp []interface{} = make([]interface{}, 0, r-l+1)
	for i <= m && j <= r {
		if cmp((*arr)[i], (*arr)[j]) <= 0 {
			tmp = append(tmp, (*arr)[i])
			i++
		} else {
			tmp = append(tmp, (*arr)[j])
			j++
		}
	}
	//当一方比较到头时将另一方剩余内容全部加入进去
	for ; i <= m; i++ {
		tmp = append(tmp, (*arr)[i])
	}
	for ; j <= r; j++ {
		tmp = append(tmp, (*arr)[j])
	}
	//将局部排序结果放入迭代器中
	for i, j = l, 0; i <= r; i, j = i+1, j+1 {
		(*arr)[i] = tmp[j]
	}
}
