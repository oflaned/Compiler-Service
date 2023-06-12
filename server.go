package main

import (
	"net/http"
	"time"
)

import (
	"context"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           "localhost:" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          //1Mb
		ReadTimeout:    10 * time.Second, //10 Sec
		WriteTimeout:   10 * time.Second, //10 Sec
	}

	return s.httpServer.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
