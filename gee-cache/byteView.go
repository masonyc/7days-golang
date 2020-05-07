package geecache

//readonly object represent cache
type ByteView struct {
	b []byte
}

//implement interface
func (bv ByteView) Len() int {
	return len(bv.b)
}

//use clone to return a copy
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (bv ByteView) String() string {
	return string(bv.b)
}
