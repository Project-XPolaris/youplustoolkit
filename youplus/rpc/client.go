package rpc

import "google.golang.org/grpc"

type YouPlusRPCClient struct {
	Address string
	Client  ServiceClient
}

func NewYouPlusRPCClient(address string) *YouPlusRPCClient {
	return &YouPlusRPCClient{Address: address}
}

func (c *YouPlusRPCClient) Connect() error {
	conn, err := grpc.Dial(
		c.Address, grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	c.Client = NewServiceClient(conn)
	return nil
}
