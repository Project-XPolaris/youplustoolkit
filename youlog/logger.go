package youlog

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

var DefaultScope = "Global"

type LogClient struct {
	client      *YouLogClient
	Logger      *logrus.Logger
	Application string
	Instance    string
}

type Scope struct {
	logClient *LogClient
	Name      string
	Fields    Fields
}
type Fields map[string]interface{}

func (c *LogClient) Init(address string, application string, instance string) error {
	c.client = NewYouLogClient(address)
	c.Application = application
	c.Instance = instance
	c.Logger = logrus.New()
	err := c.client.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (c *LogClient) Info(message string) error {
	return c.NewScope(DefaultScope).write(LEVEL_INFO, message)
}
func (c *LogClient) Debug(message string) error {
	return c.NewScope(DefaultScope).write(LEVEL_DEBUG, message)
}
func (c *LogClient) Warn(message string) error {
	return c.NewScope(DefaultScope).write(LEVEL_WARN, message)
}
func (c *LogClient) Error(message string) error {
	return c.NewScope(DefaultScope).write(LEVEL_ERROR, message)
}
func (c *LogClient) Fatal(message string) error {
	return c.NewScope(DefaultScope).write(LEVEL_FATAL, message)
}
func (c *LogClient) WithFields(field Fields) *Scope {
	return c.NewScope(DefaultScope).WithFields(field)
}

func (c *LogClient) NewScope(name string) *Scope {
	return &Scope{
		logClient: c,
		Name:      name,
		Fields:    Fields{},
	}
}

func (c *Scope) WithFields(field Fields) *Scope {
	for key, value := range field {
		c.Fields[key] = value
	}
	return c
}
func (c *Scope) write(level int64, message string) error {
	raw, err := json.Marshal(c.Fields)
	if err != nil {
		return err
	}
	_, err = c.logClient.client.Client.WriteLog(context.Background(), &LogData{
		Application: c.logClient.Application,
		Instance:    c.logClient.Instance,
		Scope:       c.Name,
		Extra:       string(raw),
		Message:     message,
		Level:       level,
		Time:        time.Now().Unix() * 1000,
	})
	return err
}
func (c *Scope) getFieldsMap() map[string]interface{} {
	return c.Fields
}
func (c *Scope) Info(message string) error {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Info(message)
	return c.write(LEVEL_INFO, message)
}
func (c *Scope) Debug(message string) error {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Debug(message)
	return c.write(LEVEL_DEBUG, message)
}
func (c *Scope) Warn(message string) error {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Warn(message)
	return c.write(LEVEL_WARN, message)
}
func (c *Scope) Error(message string) error {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Error(message)
	return c.write(LEVEL_ERROR, message)
}
func (c *Scope) Fatal(message string) error {
	err := c.write(LEVEL_FATAL, message)
	if err != nil {
		return err
	}
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Fatal(message)
	return nil
}
