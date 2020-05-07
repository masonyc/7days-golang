package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[Slow DB] Search Key", key)
			if v, ok := db[key]; ok {
				return ByteView{b: v}, nil
			}
			return nil, fmt.Println("Key does not exist")
		}))
	addr := "localhost:9999"
	peer := geecache.NewHTTPPool(addr)
	log.Println("geecache running at ", peer.self)
	log.Fatal(http.ListenAndServe(addr, peer))
}
