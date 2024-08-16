package bitmap

type Bitmap struct {
	bits []uint64
}

type bitmaper interface {
	Insert(num uint)         //在num位插入元素
	Delete(num uint)         //删除第num位
	Check(num uint) (b bool) //检查第num位是否有元素
	All() (nums []uint)      //返回所有存储的元素的下标
	Clear()                  //清空
}

func New() (bm *Bitmap) {
	return &Bitmap{
		bits: make([]uint64, 0, 0),
	}
}

func (bm *Bitmap) Insert(num uint) {
	//bm不存在时直接结束
	if bm == nil {
		return
	}
	//开始插入
	if num/64+1 > uint(len(bm.bits)) {
		//当前冗余量小于num位,需要扩增
		var tmp []uint64
		//通过冗余扩增减少扩增次数
		if num/64+1 < uint(len(bm.bits)+1024) {
			//入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
			tmp = make([]uint64, len(bm.bits)+1024)
		} else {
			//直接增加到可以容纳第num位的位置
			tmp = make([]uint64, num/64+1)
		}
		//将原有元素复制到新增的切片内,并将bm所指向的修改为扩增后的
		copy(tmp, bm.bits)
		bm.bits = tmp
	}
	//将第num位设为1即实现插入
	bm.bits[num/64] ^= 1 << (num % 64)
}

func (bm *Bitmap) Delete(num uint) {
	//bm不存在时直接结束
	if bm == nil {
		return
	}
	//num超出范围,直接结束
	if num/64+1 > uint(len(bm.bits)) {
		return
	}
	//将第num位设为0
	bm.bits[num/64] &^= 1 << (num % 64)
	if bm.bits[len(bm.bits)-1] == 0 {
		//最后一组为0,可能进行缩容
		//从后往前遍历判断可缩容内容是否小于总组数
		i := len(bm.bits) - 1
		for ; i >= 0; i-- {
			if bm.bits[i] == 0 && i != len(bm.bits)-1024 {
				continue
			} else {
				//不为0或到1024个时即可返回
				break
			}
		}
		if i <= len(bm.bits)/2 || i == len(bm.bits)-1024 {
			//小于总组数一半或超过1023个,进行缩容
			bm.bits = bm.bits[:i+1]
		}
	} else {
		return
	}
}

func (bm *Bitmap) Check(num uint) (b bool) {
	//bm不存在时直接返回false并结束
	if bm == nil {
		return false
	}
	//num超出范围,直接返回false并结束
	if num/64+1 > uint(len(bm.bits)) {
		return false
	}
	//判断第num是否为1,为1返回true,否则为false
	if bm.bits[num/64]&(1<<(num%64)) > 0 {
		return true
	}
	return false
}

func (bm *Bitmap) All() (nums []uint) {
	//对要返回的集合进行初始化,以避免返回nil
	nums = make([]uint, 0, 0)
	//bm不存在时直接返回并结束
	if bm == nil {
		return nums
	}
	//分组遍历判断某下标的元素是否存在于位图中,即其值是否为1
	for j := 0; j < len(bm.bits); j++ {
		for i := 0; i < 64; i++ {
			if bm.bits[j]&(1<<i) > 0 {
				//该元素存在,添加入结果集合内
				nums = append(nums, uint(j*64+i))
			}
		}
	}
	return nums
}

func (bm *Bitmap) Clear() {
	if bm == nil {
		return
	}
	bm.bits = make([]uint64, 0, 0)
}
