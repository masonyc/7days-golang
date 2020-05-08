package geecache

import (
	"fmt"
	"log"
	"sync"

	"github.com/masonyc/7days-golang/gee-cache/geecache/singleflight"
)

//call back used when cache is not hit
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

//Implement Getter Interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	mainCache cache
	getter    Getter
	peers     PeerPicker
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex //read write lock
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if value, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return value, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (v ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if p, ok := g.peers.PickPeer(key); ok {
				if v, err := g.getFromPeer(p, key); err == nil {
					return v, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		//use Getter call back to get item
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil

}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		log.Printf("getter get failed %s", err)
		return ByteView{}, err
	}
	bv := ByteView{b: cloneBytes(bytes)}
	g.addToCache(key, bv)
	return bv, nil
}

func (g *Group) addToCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
