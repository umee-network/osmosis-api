package server

import (
	"bytes"
	"context"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/umee-network/osmosis-api/config"
)

const socketEndpoint = "localhost:8080"

type ServerTestSuite struct {
	suite.Suite
}

func (sts *ServerTestSuite) SetupSuite() {
	config := config.Server{
		WriteTimeout: "20ms",
		ReadTimeout:  "20ms",
		ListenAddr:   socketEndpoint,
	}
	server, err := New(zerolog.Nop(), config)
	sts.Require().NoError(err)

	ctx := context.Background()
	errCh := make(chan error, 1)

	go func() {
		errCh <- server.StartServer(ctx, config)
		for err := range errCh {
			sts.Require().NoError(err)
		}
	}()
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (sts *ServerTestSuite) pongHandler(connection *websocket.Conn, done chan interface{}) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		sts.Require().NoError(err)
		sts.Require().True(bytes.Equal(msg, []byte("Pong")))
		return
	}
}

func (sts *ServerTestSuite) TestPing() {
	socketUrl := "ws://" + socketEndpoint + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	sts.Require().NoError(err)
	defer conn.Close()

	done := make(chan interface{})
	go sts.pongHandler(conn, done)

	err = conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
	sts.Require().NoError(err)

	for range done {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		sts.Require().NoError(err)
		return
	}
}
