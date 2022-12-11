package hel

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/raver119/retry"
)

type Server struct {
	Configuration
	r    *mux.Router
	w    *sync.WaitGroup
	s    *http.Server
	a    bool
	port int
}

// NewServer function creates new Server instance.
func NewServer(port int, router *mux.Router, options ...Option) (srv *Server, err error) {
	if !(port >= 1 && port < 65536) {
		return nil, fmt.Errorf("bad port value: %v", port)
	}

	srv = new(Server)

	for _, option := range options {
		option(&srv.Configuration)
	}

	if srv.Configuration.BindAddress == "" {
		srv.Configuration.BindAddress = "0.0.0.0"
	}

	srv.s = &http.Server{Addr: fmt.Sprintf("%v:%d", srv.Configuration.BindAddress, port), Handler: router}
	srv.w = &sync.WaitGroup{}
	srv.port = port

	return
}

// StartAsync function starts server asynchronously. Optionalluy blocks until server is accessible over the wire.
func (s Server) StartAsync() (err error) {
	s.a = true
	s.w.Add(1)

	go func() {
		defer s.w.Done()

		err = s.s.ListenAndServe()
	}()

	// optional blocking
	if s.BlockUntilLaunched {
		err = retry.ConnectionUntilConnectedOrTimeout(time.Minute, s.Configuration.BindAddress, s.port)
	}

	return
}

// Start function starts server synchronously.
func (s Server) Start() (err error) {
	if s.a {
		return fmt.Errorf("server was already started asynchronously")
	}

	s.a = false
	return s.s.ListenAndServe()
}

// Stop function stops server if it was started asynchronously.
func (s Server) Stop() (err error) {
	if s.a {
		_ = s.s.Shutdown(context.TODO())
		// do it
		s.w.Wait()
	}
	s.a = false
	return
}
