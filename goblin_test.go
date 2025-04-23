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
	addr    string
	http    *http.Server
	timeout time.Duration
}

func (s Server) Name() string {
	return s.addr
}

func (s Server) Serve() error {
	return s.http.ListenAndServe()
}

func (s Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.http.Shutdown(ctx)
}

type pingPongHandler struct{}

func (h pingPongHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, keyWord)
}

func TestGoblin_Awaken(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "silent success",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr:    addr,
					timeout: time.Second,
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

				if err := goblin.Awaken(goblin.WithDaemon(srv)); err != nil {
					t.Errorf("goblin awaken: %v", err)
				}
			},
		},
		{
			name: "success with logger",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr:    addr,
					timeout: time.Second,
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

				if err := goblin.Awaken(
					goblin.WithLogbook(logger),
					goblin.WithDaemon(srv),
				); err != nil {
					t.Errorf("goblin awaken: %v", err)
				}
			},
		},
		{
			name: "with context success",
			f: func(t *testing.T) {
				handler := pingPongHandler{}
				srv := Server{
					addr:    addr,
					timeout: time.Second,
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

				if err := goblin.AwakenContext(ctx, goblin.WithDaemon(srv)); err != nil {
					t.Errorf("goblin awaken: %v", err)
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
					addr:    addr,
					timeout: time.Second,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				srv2 := Server{
					addr:    addr,
					timeout: time.Second,
					http: &http.Server{
						Addr:    addr,
						Handler: handler,
					},
				}

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				if err := goblin.Awaken(
					goblin.WithLogbook(logger),
					goblin.WithDaemon(srv, srv2),
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
					addr:    addr,
					timeout: time.Second * 8,
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

				if err := goblin.Awaken(
					goblin.WithLogbook(logger),
					goblin.WithDaemon(srv),
				); err != nil {
					t.Errorf("goblin awaken: %v", err)
				}
			},
		},
		{
			name: "shutdown with error",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				srv := Server{
					addr:    addr,
					timeout: time.Second,
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

				err := goblin.Awaken(
					goblin.WithLogbook(logger),
					goblin.WithDaemon(srv),
				)
				if err == nil {
					t.Error("goblin aweken expects error got nil")
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
