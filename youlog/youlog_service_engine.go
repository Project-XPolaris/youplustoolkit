package youlog

import (
	"context"
	"encoding/json"
	"errors"
	grpcpool "github.com/processout/grpc-go-pool"
	"github.com/project-xpolaris/youplustoolkit/youlog/logservice"
	"google.golang.org/grpc"
	"time"
)

type YouLogServiceEngineConfig struct {
	Address string `mapstructure:"address"`
}
type YouLogEngine struct {
	client *YouLogClient
	Config *YouLogServiceEngineConfig
}

func (e *YouLogEngine) Init() error {
	e.client = NewYouLogClient(e.Config.Address)
	return e.client.Init()
}

func (e *YouLogEngine) WriteLog(context context.Context, scope *Scope, message string, level int64) error {
	client, conn, err := e.client.GetClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	extra, err := json.Marshal(scope.Fields)
	if err != nil {
		return err
	}
	data := &logservice.LogData{
		Application: scope.LogClient.Application,
		Instance:    scope.LogClient.Instance,
		Scope:       scope.Name,
		Extra:       string(extra),
		Message:     message,
		Level:       level,
		Time:        time.Now().UnixMilli(),
	}

	reply, err := client.WriteLog(context, data)
	if err != nil {
		return err
	}
	if reply.Success == false {
		return errors.New("service error")
	}

	return nil
}

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
func (c *YouLogClient) GetClient() (logservice.LogServiceClient, *grpcpool.ClientConn, error) {
	conn, err := c.pool.Get(context.Background())
	if err != nil {
		return nil, nil, err
	}
	client := logservice.NewLogServiceClient(conn)
	return client, conn, nil
}
