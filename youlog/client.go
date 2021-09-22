package youlog

import (
	"context"
	"fmt"
	"github.com/project-xpolaris/youplustoolkit/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

const (
	LEVEL_DEBUG = 1
	LEVEL_INFO  = 2
	LEVEL_WARN  = 3
	LEVEL_ERROR = 4
	LEVEL_FATAL = 5
)

type YouLogClient struct {
	util.BaseRPCClient
	Address  string
	Client   LogServiceClient
	Conn     *grpc.ClientConn
	tryCount int
}

func NewYouLogClient(address string) *YouLogClient {
	return &YouLogClient{Address: address}
}

func (c *YouLogClient) Connect(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx,
		c.Address, grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	c.Client = NewLogServiceClient(conn)
	return nil
}
func (c *YouLogClient) StartDaemon(maxRetry int) {
	c.BaseRPCClient = util.BaseRPCClient{
		MaxRetry:  maxRetry,
		KeepAlive: true,
	}
	c.daemon()
}
func (c *YouLogClient) daemon() {
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
