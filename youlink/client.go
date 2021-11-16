package youlink

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewYouLinkClient() *Client {
	return &Client{
		client: resty.New(),
	}
}

func (c *Client) Init(baseUrl string) {
	c.client.SetBaseURL(baseUrl)
}

func (c *Client) Callback(id string, output []*Variable, err error) error {
	if output == nil {
		output = []*Variable{}
	}
	callbackBody := map[string]interface{}{
		"id":     id,
		"output": output,
	}
	if err != nil {
		callbackBody["error"] = err.Error()
	}
	_, err = c.client.R().
		SetBody(callbackBody).
		Post("/callback")
	if err != nil {
		return err
	}
	return nil
}

type RegisterFunctionsRequestBody struct {
	Func []*FunctionTemplate `json:"func"`
}
type FunctionTemplate struct {
	Name     string                 `json:"name"`
	Template string                 `json:"template"`
	Desc     string                 `json:"desc"`
	Inputs   []*VariableDefinition  `json:"inputs"`
	Outputs  []*VariableDefinition  `json:"outputs"`
	Options  map[string]interface{} `json:"options"`
}

func (c *Client) RegisterFunctions(functions ...*Function) error {
	requestBody := RegisterFunctionsRequestBody{Func: []*FunctionTemplate{}}
	for _, function := range functions {
		requestBody.Func = append(requestBody.Func, &FunctionTemplate{
			Name:     function.Name,
			Template: function.Template,
			Desc:     function.Desc,
			Inputs:   function.InputDefinitions,
			Outputs:  function.OutputDefinitions,
			Options:  function.Options,
		})
	}
	_, err := c.client.R().
		SetBody(&requestBody).
		Post("/register")
	if err != nil {
		return err
	}
	return nil
}
