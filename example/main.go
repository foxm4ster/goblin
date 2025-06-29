package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/foxm4ster/goblin"
)

type Server struct {
	addr    string
	timeout time.Duration
	server  *http.Server
}

func NewServer(addr string, handler http.Handler, timeout time.Duration) Server {
	return Server{
		addr:    addr,
		timeout: timeout,
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s Server) ID() string {
	return s.addr
}

func (s Server) Serve() error {
	return s.server.ListenAndServe()
}

func (s Server) Shutdown(parent context.Context) error {
	if s.timeout <= 0 {
		s.timeout = time.Second * 3
	}

	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()

	time.Sleep(time.Duration(rand.IntN(5)) * time.Second)

	if rand.IntN(10)%2 == 0 {
		return errors.New("some error")
	}

	return s.server.Shutdown(ctx)
}

type pingPongHandler struct{}

func (h pingPongHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "some answer")
}

func main() {
	handler := pingPongHandler{}
	timeout := time.Second * 10

	srv := NewServer(":8080", handler, timeout)

	srv2 := NewServer(":8081", handler, timeout)
	srv3 := NewServer(":8082", handler, timeout)
	srv4 := NewServer(":8083", handler, timeout)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	conf := goblin.Config{}.WithLogFuncs(logger.Info, logger.Error)

	err := goblin.Run(conf, srv, srv2, srv3, srv4)
	if err != nil {
		logger.Error("goblin run", slog.Any("cause", err))
		return
	}
}
