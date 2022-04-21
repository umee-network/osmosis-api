package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/umee-network/osmosis-api/config"
)

type Server struct {
	http     http.Server
	upgrader websocket.Upgrader
	logger   zerolog.Logger
	cfg      config.Server
}

// SocketHandler upgrades our server to use websockets, and registers any msg
// handlers. Currently only supports "Ping" requests.
func (s *Server) SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if bytes.Equal(message, []byte("Ping")) {
			err = conn.WriteMessage(messageType, []byte("Pong"))
			if err != nil {
				break
			}
		} else {
			err = conn.WriteMessage(messageType, []byte("Invalid Request"))
			if err != nil {
				break
			}
		}
	}
}

// StartServer sets up our websocket server, given a ctx and config.
func (s *Server) StartServer(ctx context.Context, cfg config.Server) error {
	http.HandleFunc("/ws", s.SocketHandler)

	srvErrCh := make(chan error, 1)
	go func() {
		srvErrCh <- s.http.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			if err := s.http.Shutdown(shutdownCtx); err != nil {
				return err
			}

			return nil

		case err := <-srvErrCh:
			fmt.Println(err)
			return err
		}
	}
}

// New creates a new instance of the Server struct and returns any errors.
func New(logger zerolog.Logger, cfg config.Server) (Server, error) {
	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		return Server{}, err
	}
	readTimeout, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		return Server{}, err
	}

	return Server{
		logger: logger.With().Str("module", "server").Logger(),
		cfg:    cfg,
		http: http.Server{
			Addr:         cfg.ListenAddr,
			WriteTimeout: writeTimeout,
			ReadTimeout:  readTimeout,
		},
		upgrader: websocket.Upgrader{},
	}, nil
}
