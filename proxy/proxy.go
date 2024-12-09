package proxy

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	Cache "github.com/Yelsnik/caching-server/cache"
)

type Proxy struct {
	Origin string
	Cache  map[string]*Cache.CacheDB
	Mutex  sync.RWMutex
}

func NewProxyServer(origin string, clearCache bool) *Proxy {
	return &Proxy{
		Origin: origin,
		Cache:  make(map[string]*Cache.CacheDB),
	}
}

func (proxy *Proxy) ClearCache() {

	proxy.Mutex.Lock() // Lock for writing
	defer proxy.Mutex.Unlock()
	proxy.Cache = make(map[string]*Cache.CacheDB)
	fmt.Println("Cache cleared successfully")

}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cacheKey := fmt.Sprintf("%s:%s", r.Method, r.URL.String())

	cache, ok := proxy.Cache[cacheKey]

	if ok {
		writeResponse(w, *cache.Response, cache.ResponseBody, "HIT", cacheKey)
		return
	}

	// if theres no cache do this
	orginURL := proxy.Origin + r.URL.String()
	var response *http.Response
	var err error
	if r.Method == "GET" {
		response, err = http.Get(orginURL)
		if err != nil {
			http.Error(w, "Error Forwarding Request", http.StatusInternalServerError)
			return
		}
	} else {
		response, err = http.Post(orginURL, "application/json", r.Body)
		if err != nil {
			http.Error(w, "Error sending Request", http.StatusInternalServerError)
			return
		}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Error Forwarding Request Body", http.StatusInternalServerError)
		return
	}

	// cache the response
	proxy.Mutex.Lock()
	defer proxy.Mutex.Unlock()
	proxy.Cache[cacheKey] = &Cache.CacheDB{
		Response:     response,
		ResponseBody: body,
		Created:      time.Now(),
	}

	writeResponse(w, *response, body, "MISS", cacheKey)
	return
}

func writeResponse(w http.ResponseWriter, response http.Response, body []byte, cacheHeader, key string) {
	fmt.Printf("X-Cache: %s %s\n", cacheHeader, key)
	w.WriteHeader(response.StatusCode)
	for k, v := range response.Header {
		w.Header()[k] = v
	}
	w.Write(body)
}
