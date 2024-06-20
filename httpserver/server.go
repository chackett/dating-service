package httpserver

import (
	"context"
	"fmt"
	"github.com/chackett/dating-service/datingservice"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	mux    http.Handler
}

func New(port int, ds *datingservice.DateService) (*HTTPServer, error) {
	h, err := newHandler(ds)
	if err != nil {
		return nil, fmt.Errorf("unable to create handler: %w", err)
	}

	mws := []func(handler2 http.Handler) http.Handler{
		h.middlewareAuth, h.middlewareAuth,
	}

	h.setupRoutes(mws)

	result := &HTTPServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: h.mux,
		},
	}

	return result, nil
}

func (s *HTTPServer) Serve() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("unable to start listening: %w", err)
	}
	return nil
}

func (s *HTTPServer) Close() error {
	fmt.Println("shutting down")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("unable to close web server listener: %w", err)
	}

	return nil
}
