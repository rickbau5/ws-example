package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rickbau5/ws-example/internal/handlers"
	"golang.org/x/sync/errgroup"
)

var logger = log.New(os.Stdout, "[ws-server] ", log.LstdFlags)

func Server(ctx context.Context) error {
	grp, grpCtx := errgroup.WithContext(context.Background())

	host, err := os.Hostname()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/healthcheck", &handlers.Health{Start: time.Now()})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/ws", &handlers.WebSocketHandler{
		ID: host,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Logger: logger,
	})

	server := http.Server{
		Addr:    os.Getenv("HTTP_SERVER_ADDR"),
		Handler: mux,
		ConnState: func(conn net.Conn, state http.ConnState) {
			u, err := url.Parse(conn.RemoteAddr().String())
			if err != nil {
				return
			}
			// ignore connections from localhost
			if u.Hostname() == "127.0.0.1" {
				return
			}

			logger.Printf("connection state change: %s -> %s\n", conn.RemoteAddr(), state)
		},
		BaseContext: func(_ net.Listener) context.Context { return grpCtx },
	}

	grp.Go(server.ListenAndServe)

	// wait for shutdown / error
	grp.Go(func() error {
		select {
		case <-ctx.Done(): // parent cancelled
		case <-grpCtx.Done(): // server cancelled
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error in server shutdown: %w", err)
		}

		return errors.New("shutting down")
	})

	// all grp goroutines will return an error in all cases which ensures this will close
	return grp.Wait()
}
