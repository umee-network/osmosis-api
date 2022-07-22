package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/osmosis-labs/osmosis/x/gamm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// Client is a struct for querying transactions for a cosmos chain.
type Client struct {
	grpcEndpoint   string
	tls            bool
	headers        []string
	grpcTimeout    time.Duration
	grpcConnection *grpc.ClientConn
}

// NewClient returns a new Client instance with an active gRPC connection.
// Headers are passed as map of strings and translated to a key-value based array.
func NewClient(
	grpcEndpoint string,
	tls bool,
	headers map[string]string,
	grpcTimeout time.Duration,
) (Client, error) {
	headersKV := make([]string, len(headers)*2)
	index := 0
	for key, value := range headers {
		headersKV[index] = key
		headersKV[index+1] = value
		index = index + 2
	}

	c := Client{
		grpcEndpoint: grpcEndpoint,
		tls:          tls,
		headers:      headersKV,
		grpcTimeout:  grpcTimeout,
	}
	err := c.connectGRPC()
	if err != nil {
		return Client{}, err
	}

	return c, nil
}

// connectGRPC dials up our grpc connection endpoint.
func (c *Client) connectGRPC() error {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithContextDialer(dial),
	}
	if c.tls {
		grpcOpts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithContextDialer(dial),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{}),
		}
	}

	grpcConn, err := grpc.Dial(
		c.grpcEndpoint,
		grpcOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to dial Cosmos gRPC service: %w", err)
	}

	c.grpcConnection = grpcConn
	return nil
}

// GetSpotPrice returns the spot price of an asset in a given pool.
func (c *Client) GetSpotPrice(ctx context.Context, req types.QuerySpotPriceRequest) (*types.QuerySpotPriceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.grpcTimeout)
	ctx = metadata.AppendToOutgoingContext(ctx, c.headers...)
	defer cancel()

	queryClient := types.NewQueryClient(c.grpcConnection)

	res, err := queryClient.SpotPrice(ctx, &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
