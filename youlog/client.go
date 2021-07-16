package youlog

import "google.golang.org/grpc"

const (
	LEVEL_DEBUG = 1
	LEVEL_INFO  = 2
	LEVEL_WARN  = 3
	LEVEL_ERROR = 4
	LEVEL_FATAL = 5
)

type YouLogClient struct {
	Address string
	Client  LogServiceClient
}

func NewYouLogClient(address string) *YouLogClient {
	return &YouLogClient{Address: address}
}

func (c *YouLogClient) Connect() error {
	conn, err := grpc.Dial(
		c.Address, grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	c.Client = NewLogServiceClient(conn)
	return nil
}
