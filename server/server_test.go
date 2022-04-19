package server

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/umee-network/osmosis-api/config"
)

const socketEndpoint = "localhost:8080"

type ServerTestSuite struct {
	suite.Suite

	server Server
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

	go server.StartServer(ctx, config)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (sts ServerTestSuite) pongHandler(connection *websocket.Conn, done chan interface{}) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		sts.Require().NoError(err)
		sts.Require().Equal(bytes.Compare(msg, []byte("Pong")), 0)
		return
	}
}

func (sts ServerTestSuite) TestPing() {
	socketUrl := "ws://" + socketEndpoint + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	sts.Require().NoError(err)
	defer conn.Close()

	done := make(chan interface{})
	go sts.pongHandler(conn, done)

	for {
		select {
		case <-done:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			sts.Require().NoError(err)
			return
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
			sts.Require().NoError(err)
		}
	}
}
