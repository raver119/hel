package hel

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	r *mux.Router
	w *sync.WaitGroup
	s *http.Server
	a bool
}

func NewServer(port int, router *mux.Router) (srv *Server, err error) {
	if !(port >= 1 && port < 65536) {
		return nil, fmt.Errorf("bad port value: %v", port)
	}

	srv = new(Server)

	srv.s = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: router}
	srv.w = &sync.WaitGroup{}

	return
}

func (s Server) StartAsync() (err error) {
	s.a = true
	s.w.Add(1)

	go func() {
		defer s.w.Done()

		err = s.s.ListenAndServe()
	}()
	return
}

func (s Server) Start() (err error) {
	if s.a {
		return fmt.Errorf("server was already started asynchronously")
	}

	s.a = false
	return s.s.ListenAndServe()
}

func (s Server) Stop() (err error) {
	if s.a {
		_ = s.s.Shutdown(context.TODO())
		// do it
		s.w.Wait()
	}
	s.a = false
	return
}
