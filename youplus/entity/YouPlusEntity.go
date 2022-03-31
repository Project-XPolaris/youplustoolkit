package entry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/project-xpolaris/youplustoolkit/youplus/rpc"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	ErrorCodeUnknown        = 9999
	ErrorCodeEntityNotFound = 6001
)

type EntityExport struct {
	Urls  []string    `json:"urls"`
	Extra interface{} `json:"extra"`
}
type EntityClient struct {
	Name                     string
	Version                  int64
	Instance                 string
	Export                   interface{}
	Client                   *rpc.YouPlusRPCClient
	HeartbeatRate            int64
	StopHeartbeatContext     context.Context
	StopHeartbeatContextFunc context.CancelFunc
}

func NewEntityClient(name string, version int64, export interface{}, client *rpc.YouPlusRPCClient) *EntityClient {
	instance := fmt.Sprintf("%s_%s", name, xid.New().String())
	return &EntityClient{
		Name:     name,
		Version:  version,
		Export:   export,
		Client:   client,
		Instance: instance,
	}
}

func (e *EntityClient) Register() error {
	client, conn, err := e.Client.GetClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	result, err := client.RegisterEntry(context.Background(), &rpc.RegisterEntryRequest{
		Name:     &e.Name,
		Version:  &e.Version,
		Instance: &e.Instance,
	})
	if err != nil {
		return err
	}
	if !result.GetSuccess() {
		return errors.New(result.GetReason())
	}
	return nil
}
func (e *EntityClient) Unregister() error {
	client, conn, err := e.Client.GetClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	result, err := client.UnregisterEntry(context.Background(), &rpc.UnregisterEntryRequest{
		Instance: &e.Instance,
	})
	if err != nil {
		return err
	}
	if !result.GetSuccess() {
		return errors.New(result.GetReason())
	}
	return nil
}
func (e *EntityClient) UpdateExport(data interface{}) error {
	client, conn, err := e.Client.GetClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	raw, err := json.Marshal(&data)
	e.Export = data
	if err != nil {
		return err
	}
	rawString := string(raw)
	result, err := client.UpdateEntryExport(context.Background(), &rpc.UpdateEntryExportRequest{
		Instance: &e.Instance,
		Data:     &rawString,
	})
	if err != nil {
		return err
	}
	if !result.GetSuccess() {
		return errors.New(result.GetReason())
	}
	return nil
}
func (e *EntityClient) StartHeartbeat(ctx context.Context) error {
	if e.StopHeartbeatContext != nil {
		return errors.New("only one heartbeat")
	}
	e.StopHeartbeatContext, e.StopHeartbeatContextFunc = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-time.After(time.Duration(e.HeartbeatRate) * time.Millisecond):
				client, conn, err := e.Client.GetClient()

				if err != nil {
					logrus.Info(err)
					conn.Close()
					continue
				}
				state := "online"
				reply, err := client.EntryHeartbeat(ctx, &rpc.HeartbeatRequest{
					Name:     &e.Name,
					Instance: &e.Instance,
					State:    &state,
				})
				if err != nil {
					logrus.Info(err)
					conn.Close()
					continue
				}
				if !*reply.Success {
					logrus.Info(*reply.Reason)
					if *reply.Code == ErrorCodeEntityNotFound {
						logrus.Info("try to register entity again")
						e.Register()
						e.UpdateExport(e.Export)
					}
				}
				conn.Close()
			case <-e.StopHeartbeatContext.Done():
				return
			}
		}
	}()
	return nil
}
