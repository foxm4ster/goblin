package goblin_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/foxm4ster/goblin"
)

const (
	addr = ":9999"
	host = "http://localhost:9999"

	keyWord = "pong"
)

type Server struct {
	addr string
	http *http.Server
}

func (s Server) ID() string {
	return s.addr
}

func (s Server) Serve() error {
	return s.http.ListenAndServe()
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

type pingPongHandler struct{}

func (h pingPongHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, keyWord)
}

func TestGoblin_Run(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "silent success",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				opts := []goblin.Option{
					goblin.WithShutdownTimeout(time.Second),
				}

				if err := goblin.Run(opts, srv); err != nil {
					t.Errorf("goblin run: %v", err)
				}
			},
		},
		{
			name: "silent success 2",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				if err := goblin.Run(nil, srv); err != nil {
					t.Errorf("goblin run: %v", err)
				}
			},
		},
		{
			name: "success with logger",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

				opts := []goblin.Option{
					goblin.WithLogFuncs(logger.Info, logger.Error),
					goblin.WithShutdownTimeout(time.Second),
				}

				if err := goblin.Run(opts, srv); err != nil {
					t.Errorf("goblin run: %v", err)
				}
			},
		},
		{
			name: "with context success",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := goblin.RunContext(
					ctx,
					[]goblin.Option{goblin.WithShutdownTimeout(time.Second)},
					srv,
				); err != nil {
					t.Errorf("goblin run: %v", err)
				}
			},
		},
		{
			name: "server already started",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				handler := pingPongHandler{}
				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				srv2 := Server{
					addr: addr,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				if err := goblin.Run(
					[]goblin.Option{
						goblin.WithLogFuncs(logger.Info, logger.Error),
						goblin.WithShutdownTimeout(time.Second),
					},
					srv,
					srv2,
				); err == nil {
					t.Errorf("goblin expected err, got nil")
					return
				}

				if !bytes.Contains(buff.Bytes(), []byte("listen tcp :9999: bind: address already in use")) {
					t.Errorf("unexpected error message: %s", buff.String())
				}
			},
		},
		{
			name: "graceful shutdown",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr: addr,
						Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							time.Sleep(2 * time.Second)
							_, _ = fmt.Fprintf(w, keyWord)
						}),
						ReadTimeout:       time.Second * 10,
						ReadHeaderTimeout: time.Second * 10,
						WriteTimeout:      time.Second * 10,
						IdleTimeout:       time.Second * 10,
					},
				}

				first, second, third := make(chan struct{}), make(chan struct{}), make(chan struct{})

				go func() {
					<-first
					go func() {
						second <- struct{}{}
					}()

					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send first request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					<-third
					_, err := http.Get(host)
					if err == nil {
						t.Errorf("expected error as second request response, got nil")
						return
					}

					if !errors.Is(err, syscall.ECONNREFUSED) {
						t.Errorf("unexpected error, got %v, want %v", err, syscall.ECONNREFUSED)
					}
				}()

				go func() {
					<-second
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					third <- struct{}{}
				}()

				go func() {
					first <- struct{}{}
				}()

				if err := goblin.Run(
					[]goblin.Option{
						goblin.WithLogFuncs(logger.Info, logger.Error),
						goblin.WithShutdownTimeout(time.Second * 8),
					},
					srv,
				); err != nil {
					t.Errorf("goblin run: %v", err)
				}
			},
		},
		{
			name: "shutdown with error",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				srv := Server{
					addr: addr,
					http: &http.Server{
						Addr: addr,
						Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							time.Sleep(time.Second * 2)
							_, _ = fmt.Fprintf(w, keyWord)
						}),
						ReadTimeout:       time.Second * 10,
						ReadHeaderTimeout: time.Second * 10,
						WriteTimeout:      time.Second * 10,
						IdleTimeout:       time.Second * 10,
					},
				}

				first, second := make(chan struct{}), make(chan struct{})

				go func() {
					<-first
					go func() {
						second <- struct{}{}
					}()

					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send first request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					<-second
					time.Sleep(time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				go func() {
					first <- struct{}{}
				}()

				err := goblin.Run(
					[]goblin.Option{
						goblin.WithLogFuncs(logger.Info, logger.Error),
						goblin.WithShutdownTimeout(time.Second),
					},
					srv,
				)
				if err == nil {
					t.Error("goblin run expects error got nil")
				}

				if !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("goblin expected context deadline exceeded, got: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.f)
	}
}
