package rpc

import (
	"context"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type YouPlusRPCClient struct {
	Address string
	pool    *grpcpool.Pool
}

func NewYouPlusRPCClient(address string) *YouPlusRPCClient {
	return &YouPlusRPCClient{
		Address: address,
	}
}

func (c *YouPlusRPCClient) Init() error {
	var factory grpcpool.Factory = func() (*grpc.ClientConn, error) {
		return grpc.DialContext(context.Background(), c.Address, grpc.WithInsecure())
	}
	pool, err := grpcpool.New(factory, 1, 3, 0)
	if err != nil {
		return err
	}
	c.pool = pool
	return nil
}
func (c *YouPlusRPCClient) GetClient() (YouPlusServiceClient, *grpcpool.ClientConn, error) {
	conn, err := c.pool.Get(context.Background())
	if err != nil {
		return nil, nil, err
	}
	client := NewYouPlusServiceClient(conn)
	return client, conn, nil
}
