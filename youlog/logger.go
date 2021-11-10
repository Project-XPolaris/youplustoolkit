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

func (c *LogClient) Init(address string, application string, instance string) {
	c.client = NewYouLogClient(address)
	c.Application = application
	c.Instance = instance
	c.Logger = logrus.New()
	return
}
func (c *LogClient) Connect(context context.Context) error {
	err := c.client.Connect(context)
	return err
}
func (c *LogClient) StartDaemon(maxRetry int) {
	c.client.StartDaemon(maxRetry)
}

func (c *LogClient) Info(message string) {
	err := c.NewScope(DefaultScope).write(LEVEL_INFO, message)
	if err != nil {
		c.Logger.Error(err)
	}
}
func (c *LogClient) Debug(message string) {
	err := c.NewScope(DefaultScope).write(LEVEL_DEBUG, message)
	if err != nil {
		c.Logger.Error(err)
	}
}
func (c *LogClient) Warn(message string) {
	err := c.NewScope(DefaultScope).write(LEVEL_WARN, message)
	if err != nil {
		c.Logger.Error(err)
	}
}
func (c *LogClient) Error(message string) {
	err := c.NewScope(DefaultScope).write(LEVEL_ERROR, message)
	if err != nil {
		c.Logger.Error(err)
	}
}
func (c *LogClient) Err(err error) {
	writeErr := c.NewScope(DefaultScope).write(LEVEL_ERROR, err.Error())
	if err != nil {
		c.Logger.Error(writeErr)
	}
}
func (c *LogClient) Fatal(message string) {
	err := c.NewScope(DefaultScope).write(LEVEL_FATAL, message)
	if err != nil {
		c.Logger.Error(err)
	}
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
func (c *Scope) SetFields(fields Fields) *Scope {
	c.Fields = fields
	return c
}
func (c *Scope) WithFields(fields Fields) *Scope {
	newScope := c.logClient.NewScope(c.Name)
	for key, value := range fields {
		newScope.Fields[key] = value
	}
	for key, value := range c.Fields {
		newScope.Fields[key] = value
	}
	return newScope
}
func (c *Scope) write(level int64, message string) error {
	raw, err := json.Marshal(c.Fields)
	if err != nil {
		return err
	}
	if c.logClient.client.Client != nil {
		_, err = c.logClient.client.Client.WriteLog(context.Background(), &LogData{
			Application: c.logClient.Application,
			Instance:    c.logClient.Instance,
			Scope:       c.Name,
			Extra:       string(raw),
			Message:     message,
			Level:       level,
			Time:        time.Now().Unix() * 1000,
		})
	}
	return err
}
func (c *Scope) getFieldsMap() map[string]interface{} {
	return c.Fields
}
func (c *Scope) Info(message string) {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Info(message)
	err := c.write(LEVEL_INFO, message)
	if err != nil {
		c.logClient.Logger.Error(err)
	}
}
func (c *Scope) Debug(message string) {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Debug(message)
	err := c.write(LEVEL_DEBUG, message)
	if err != nil {
		c.logClient.Logger.Error(err)
	}
}
func (c *Scope) Warn(message string) {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Warn(message)
	err := c.write(LEVEL_WARN, message)
	if err != nil {
		c.logClient.Logger.Error(err)
	}
}
func (c *Scope) Error(message string) {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Error(message)
	err := c.write(LEVEL_ERROR, message)
	if err != nil {
		c.logClient.Logger.Error(err)
	}
}
func (c *Scope) Fatal(message string) {
	c.logClient.Logger.WithFields(map[string]interface{}{
		"scope": c.Name,
	}).WithFields(c.getFieldsMap()).Fatal(message)
	err := c.write(LEVEL_FATAL, message)
	if err != nil {
		c.logClient.Logger.Error(err)
	}
}
