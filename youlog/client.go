package youlog

import (
	"context"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

const (
	LEVEL_DEBUG = 1
	LEVEL_INFO  = 2
	LEVEL_WARN  = 3
	LEVEL_ERROR = 4
	LEVEL_FATAL = 5
)

type YouLogClient struct {
	Address string
	pool    *grpcpool.Pool
}

func NewYouLogClient(address string) *YouLogClient {
	return &YouLogClient{Address: address}
}

func (c *YouLogClient) Init() error {
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
func (c *YouLogClient) GetClient() (LogServiceClient, error) {
	conn, err := c.pool.Get(context.Background())
	if err != nil {
		return nil, err
	}
	client := NewLogServiceClient(conn)
	return client, nil
}
