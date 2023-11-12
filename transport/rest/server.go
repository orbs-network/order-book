package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server      *http.Server
	router      *mux.Router
	StopChannel chan os.Signal
}

func NewHTTPServer(addr string, router *mux.Router) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		router:      router,
		StopChannel: make(chan os.Signal, 1),
	}
}

func (hs *HTTPServer) StartServer() {
	go func() {
		fmt.Printf("HTTP server started on %s\n", hs.server.Addr)
		if err := hs.server.ListenAndServe(); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}

func (hs *HTTPServer) StopServer(ctx context.Context) {
	fmt.Println("Shutting down the HTTP server...")
	if err := hs.server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}
}
