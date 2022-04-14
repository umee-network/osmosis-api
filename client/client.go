package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	tx "github.com/cosmos/cosmos-sdk/api/cosmos/tx/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is a struct for querying transactions for a cosmos chain.
type Client struct {
	grpcEndpoint string
	tls          bool
}

// NewClient returns a new Client instance.
func NewClient(grpcEndpoint string, tls bool) Client {
	return Client{
		grpcEndpoint: grpcEndpoint,
		tls:          tls,
	}
}

// GetTxs returns the transactions which match the events queried,
// e.g. ["message.action='/cosmos.bank.v1beta1.Msg/Send'"]
func (c *Client) GetTxs(events []string) (tx.GetTxsEventResponse, error) {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialerFunc),
	}
	if c.tls {
		grpcOpts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithContextDialer(dialerFunc),
		}
	}

	grpcConn, err := grpc.Dial(
		c.grpcEndpoint,
		grpcOpts...,
	)
	if err != nil {
		return tx.GetTxsEventResponse{}, fmt.Errorf("failed to dial Cosmos gRPC service: %w", err)
	}

	defer grpcConn.Close()
	queryClient := tx.NewServiceClient(grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	queryResponse, err := queryClient.GetTxsEvent(
		ctx, &tx.GetTxsEventRequest{
			Events: events,
		},
	)
	if err != nil {
		return tx.GetTxsEventResponse{}, fmt.Errorf("failed to make grpc query: %w", err)
	}

	return *queryResponse, nil
}
