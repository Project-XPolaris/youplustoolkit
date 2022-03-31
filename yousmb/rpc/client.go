package rpc

import (
	"context"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type YouSMBRPCClient struct {
	Address string
	pool    *grpcpool.Pool
}

func NewYouSMBRPCClient(address string) *YouSMBRPCClient {
	return &YouSMBRPCClient{
		Address: address,
	}
}

func (c *YouSMBRPCClient) Init() error {
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
func (c *YouSMBRPCClient) GetClient() (YouSMBServiceClient, *grpcpool.ClientConn, error) {
	conn, err := c.pool.Get(context.Background())
	if err != nil {
		return nil, nil, err
	}
	client := NewYouSMBServiceClient(conn)
	return client, conn, nil
}
