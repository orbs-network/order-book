package rest

import (
	"context"
	"fmt"
	"log"
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
		// ListenAndServe returns ErrServerClosed on graceful shutdown
		if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("HTTP server error: %v", err))
		}
	}()
}

func (hs *HTTPServer) StopServer(ctx context.Context) error {
	log.Printf("Shutting down the HTTP server...\n")
	if err := hs.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %v", err)
	}
	return nil
}
