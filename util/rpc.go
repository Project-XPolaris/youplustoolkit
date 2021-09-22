package util

type BaseRPCClient struct {
	KeepAlive bool
	MaxRetry  int
}

func (b *BaseRPCClient) IsKeepAlive() bool {
	return b.KeepAlive
}

func (b *BaseRPCClient) MaxRetryCount() int {
	return b.MaxRetry
}

type RPCClient interface {
	IsKeepAlive() bool
	MaxRetryCount() int
}
