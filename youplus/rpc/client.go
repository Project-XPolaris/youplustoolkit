package rpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

type YouPlusRPCClient struct {
	Address   string
	Client    YouPlusServiceClient
	Conn      *grpc.ClientConn
	KeepAlive bool
	MaxRetry  int
	tryCount  int
}

func NewYouPlusRPCClient(address string) *YouPlusRPCClient {
	return &YouPlusRPCClient{
		Address: address,
	}
}

func (c *YouPlusRPCClient) daemon() {
	go func() {
		for {
			if c.Conn != nil && c.Conn.GetState() == connectivity.TransientFailure {
				if c.tryCount == c.MaxRetry {
					return
				}
				logrus.Info(fmt.Sprintf("youplus rpc connect lost,try to connect [%d of %d]", c.tryCount, c.MaxRetry))
				connContext, _ := context.WithTimeout(context.Background(), 3*time.Second)
				err := c.Connect(connContext)
				if err != nil {
					logrus.Error(err)
					c.tryCount += 1
					continue
				}
				c.tryCount = 0
			}
		}

	}()
}

func (c *YouPlusRPCClient) Connect(ctx context.Context) error {
	conn, err := grpc.DialContext(
		ctx,
		c.Address, grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	c.Client = NewYouPlusServiceClient(conn)
	c.Conn = conn
	if c.KeepAlive {
		c.daemon()
	}
	return nil
}
