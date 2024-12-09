package proxy

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Origin     string
	Port       int
	ClearCache bool
	Proxy      *Proxy
}

func (server *Server) StartServer() {
	if server.ClearCache {
		fmt.Println("clearing cache...")
		server.Proxy.ClearCache()
	}

	if server.Origin != "" || server.Port != 0 {
		fmt.Printf("Starting caching proxy server on port %d and forwarding requests to %s\n", server.Port, server.Origin)
		http.Handle("/", server.Proxy)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil))
	}
}
