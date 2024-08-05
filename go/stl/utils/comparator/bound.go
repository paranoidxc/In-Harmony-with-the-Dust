package comparator

// 通过传入的比较函数对待查找数组进行查找以获取待查找元素的上界即不大于它的最大值的下标
// 以传入的比较函数进行比较
// 如果该元素存在,则上界指向元素为该元素
// 如果该元素不存在,上界指向元素为该元素的前一个元素
func UpperBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr == nil || (*arr) == nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的上界
	return upperBound(arr, e, cmp)
}

// 通过传入的比较函数对待查找数组进行查找以获取待查找元素的上界即不大于它的最大值的下标
// 以传入的比较函数进行比较
// 如果该元素存在,则上界指向元素为该元素,且为最右侧
// 如果该元素不存在,上界指向元素为该元素的前一个元素
// 以二分查找的方式寻找该元素的上界
func upperBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr))/2, len((*arr))-1
	for l < r {
		m = (l + r + 1) / 2
		if cmp((*arr)[m], e) <= 0 {
			l = m
		} else {
			r = m - 1
		}
	}
	return l
}

// 通过传入的比较函数对待查找数组进行查找以获取待查找元素的下界即不小于它的最小值的下标
// 以传入的比较函数进行比较
// 如果该元素存在,则上界指向元素为该元素
// 如果该元素不存在,上界指向元素为该元素的后一个元素
func LowerBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr == nil || (*arr) == nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的下界
	return lowerBound(arr, e, cmp)
}

// 通过传入的比较函数对待查找数组进行查找以获取待查找元素的下界即不小于它的最小值的下标
// 以传入的比较函数进行比较
// 如果该元素存在,则上界指向元素为该元素,且为最右侧
// 如果该元素不存在,上界指向元素为该元素的后一个元素
func lowerBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr))/2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) >= 0 {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
