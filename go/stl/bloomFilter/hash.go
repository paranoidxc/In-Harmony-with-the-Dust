package bloomFilter

import "fmt"

type Hash func(v interface{}) (h uint32)

func hash(v interface{}) (h uint32) {
	h = uint32(0)
	s := fmt.Sprintf("131-%v-%v", v, v)
	bs := []byte(s)
	for i := range bs {
		h += uint32(bs[i]) * 131
	}
	return h
}
