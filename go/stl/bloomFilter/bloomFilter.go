package bloomFilter

type BloomFilter struct {
	bits []uint64
	hash Hash
}

type BloomFilteror interface {
	Insert(v interface{})         //向布隆过滤器中插入v
	Check(v interface{}) (b bool) //检查该值是否存在于布隆过滤器中,该校验存在误差
	Clear()                       //清空该布隆过滤器
}

func New(h Hash) (bf *BloomFilter) {
	if h == nil {
		h = hash
	}
	return &BloomFilter{
		bits: make([]uint64, 0, 0),
		hash: h,
	}
}

func (bf *BloomFilter) Insert(v interface{}) {
	//bm不存在时直接结束
	if bf == nil {
		return
	}
	//开始插入
	h := bf.hash(v)
	if h/64+1 > uint32(len(bf.bits)) {
		//当前冗余量小于num位,需要扩增
		var tmp []uint64
		//通过冗余扩增减少扩增次数
		if h/64+1 < uint32(len(bf.bits)+1024) {
			//入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
			tmp = make([]uint64, len(bf.bits)+1024)
		} else {
			//直接增加到可以容纳第num位的位置
			tmp = make([]uint64, h/64+1)
		}
		//将原有元素复制到新增的切片内,并将bm所指向的修改为扩增后的
		copy(tmp, bf.bits)
		bf.bits = tmp
	}
	//将第num位设为1即实现插入
	bf.bits[h/64] ^= 1 << (h % 64)
}

func (bf *BloomFilter) Check(v interface{}) (b bool) {
	//bf不存在时直接返回false并结束
	if bf == nil {
		return false
	}
	h := bf.hash(v)
	//h超出范围,直接返回false并结束
	if h/64+1 > uint32(len(bf.bits)) {
		return false
	}
	//判断第num是否为1,为1返回true,否则为false
	if bf.bits[h/64]&(1<<(h%64)) > 0 {
		return true
	}
	return false
}

func (bf *BloomFilter) Clear() {
	if bf == nil {
		return
	}
	bf.bits = make([]uint64, 0, 0)
}
