package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BNode struct {
	data []byte
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

const HEADER = 4
const BTREE_PAGE_SIZE = 4096 // 4k
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BTree struct {
	root uint64
	get  func(uint64) BNode
	new  func(BNode) uint64
	del  func(uint64)
}

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	fmt.Println(node1max)
	assert(node1max <= BTREE_PAGE_SIZE)
}

//header

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointers
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}
func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}
func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}
func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}
func (node BNode) getVal(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

// B-Tree Insertion
func nodeLoopupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}

	return found
}

func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	assert(srcOld+n <= old.nkeys())
	assert(dstNew+n <= new.nkeys())
	if n == 0 {
		return
	}

	// pointers
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}

	// offsets
	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(dstNew)
	for i := uint16(1); i <= n; i++ {
		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
		new.setOffset(dstNew+i, offset)
	}

	// KVs
	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)
	copy(new.data[new.kvPos(dstNew):1], old.data[begin:end])

}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	// ptrs
	new.setPtr(idx, ptr)
	// KVs
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(key)))
	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)

	new.setOffset(idx+1, new.getOffset(idx)+4+uint16((len(key)+len(val))))
}

func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

	// where to insrt the key
	idx := nodeLoopupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new, node, idx, key, val)
		} else {
			leafInsert(new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("bad node")
	}

	return new
}

func nodeInsert(
	tree *BTree, new BNode, node BNode, idx uint16,
	key []byte, val []byte) {

	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)

	knode = treeInsert(tree, knode, key, val)
	nsplit, splited := nodeSplite3(knode)
	nodeReplaceKidN(tree, new, node, idx, splited[:nsplit]...)
}

// split a bigger-than-allowed node into two.
// the second node always fits on a page.
func nodeSplit2(left BNode, right BNode, old BNode) {
	// 计算分割点
	mid := old.nkeys() / 2

	// 复制左半部分到left节点
	left.setHeader(old.btype(), mid)
	nodeAppendRange(left, old, 0, 0, mid)

	// 复制右半部分到right节点
	right.setHeader(old.btype(), old.nkeys()-mid)
	nodeAppendRange(right, old, 0, mid, old.nkeys()-mid)

	// 如果是内部节点,需要将中间键上移
	if old.btype() == BNODE_NODE {
		// 将中间键添加到left节点末尾
		nodeAppendKV(left, mid, right.getPtr(0), right.getKey(0), nil)
		// 从right节点移除第一个键值对
		nodeAppendRange(right, right, 0, 1, right.nkeys()-1)
		right.setHeader(BNODE_NODE, right.nkeys()-1)
	}
}

func nodeSplite3(old BNode) (uint16, [3]BNode) {
	if old.btype() <= BTREE_PAGE_SIZE {
		old.data = old.data[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}
	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)}
	right := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left.data = left.data[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	leftleft := BNode{make([]byte, BTREE_PAGE_SIZE)}
	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}

	nodeSplit2(leftleft, middle, left)

	assert(leftleft.nbytes() <= BTREE_PAGE_SIZE)
	return 3, [3]BNode{leftleft, middle, right}
}

// replace a link with multiple links

func nodeReplaceKidN(tree *BTree, new BNode, old BNode, idx uint16, kids ...BNode) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}


// remove a key from a leaf node

func leafDelete(new BNode, old BNode, idx uint16) { 
	new.setHeader(BNODE_LEAF, old.nkeys()-1) 
	nodeAppendRange(new, old, 0, 0, idx) 
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-(idx+1)) 
}

// delete a key from the tree

func treeDelete(tree * BTree, node BNode, key []byte) BNode { // where to find the key?

	idx := nodeLoopupLE(node, key) // act depending on the node type
	
	switch node.btype() { 
		case BNODE_LEAF:
			if !bytes.Equal(key, node.getKey(idx)) { 
				return BNode{} // not found 
				} 
				// delete the key in the leaf 
				new := BNode{data: make([]byte, BTREE_PAGE_SIZE)
					} 
					leafDelete(new, node, idx) 
					return new
				case BNODE_NODE:

					return nodeDelete(tree, node, idx, key) default:
					
					panic("bad node!") }
					}


					// part of the treeDelete()

func nodeDelete(tree * BTree, node BNode, idx uint16, key []byte) BNode { 
	// recurse into the kid 
	kptr := node.getPtr(idx) 
	updated := treeDelete(tree, tree.get(kptr), key) 
	if len(updated.data) == 0 { return BNode{} // not found 
	} 
	tree.del(kptr)
	new := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
	 // check for merging 
	 mergeDir, sibling := shouldMerge(tree, node, idx, updated) 
	 switch { 
		case mergeDir < 0: // left 
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, sibling, updated) 
		tree.del(node.getPtr(idx - 1)) 
		nodeReplace2Kid(new, node, idx-1, tree.new(merged), merged.getKey(0))
		 case mergeDir > 0: // right 
		 merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		 nodeMerge(merged, updated, sibling) 
		 tree.del(node.getPtr(idx + 1)) 
		 nodeReplace2Kid(new, node, idx, tree.new(merged), merged.getKey(0))
		case mergeDir == 0:

			assert(updated.nkeys() > 0)
			 nodeReplaceKidN(tree, new, node, idx, updated) 
			 } 
			 return new
			
			}

// merge 2 nodes into 1

func nodeMerge(new BNode, left BNode, right BNode) {
	new.setHeader(left.btype(), left.nkeys()+right.nkeys())
	nodeAppendRange(new, left, 0, 0, left.nkeys())
	nodeAppendRange(new, right, left.nkeys(), 0, right.nkeys())
}


func shouldMerge( tree * BTree, node BNode, idx uint16, updated BNode, ) (int, BNode) {
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}

	if idx > 0 {
		sibling := tree.get(node.getPtr(idx - 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}

	}
	if idx+1 < node.nkeys() {
		sibling := tree.get(node.getPtr(idx + 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling
		}

	}
	return 0, BNode{}
}

func main() {
	fmt.Println("vim-go")
}

func assert(val bool) {
	if !val {
		panic("sdfdsf")
	}
}
