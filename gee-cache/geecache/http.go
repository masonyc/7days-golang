package geecache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/masonyc/7days-golang/gee-cache/geecache/consistenthash"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

// HTTP SERVER / PEER
type HTTPPool struct {
	basePath    string
	self        string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

//init peers and getters in http pool object
func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaultReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		return h.httpGetters[key], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

// HTTP CLIENT
type httpGetter struct {
	baseURL string
}

var _ PeerGetter = (*httpGetter)(nil)

func (hg *httpGetter) Get(group string, key string) ([]byte, error) {
	//  set up cache url
	u := fmt.Sprintf("%v%v/%v",
		hg.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return b, nil
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<base>/<group>/<key>
	// ignore base and split after
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	group := parts[0]
	key := parts[1]
	g := GetGroup(group)
	if g == nil {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}
	b, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b.ByteSlice())
}
