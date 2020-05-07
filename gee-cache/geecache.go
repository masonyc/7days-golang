package geecache

import (
	"fmt"
	"log"
	"sync"
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
}

var (
	mu     sync.RWMutex //read write lock
	groups = make(map[string]*Group)
)

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

func (g *Group) load(key string) (ByteView, error) {
	//use Getter call back to get item
	return g.getLocally(key)
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
