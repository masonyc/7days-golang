package geecache

//must be implemented to locate the peer that owns a specific key
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// must be implemented by a peer
// search data from cache
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
