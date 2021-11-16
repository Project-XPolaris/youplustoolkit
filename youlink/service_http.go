package youlink

import (
	"fmt"
	"github.com/allentom/haruka"
	"net/http"
)

type ResponseBody struct {
	CallbackId string      `json:"callbackId"`
	Inputs     []*Variable `json:"inputs"`
	Output     []*Variable `json:"outputs"`
}
type VariableDefinition struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}
type Variable struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
type Function struct {
	Endpoint          string                `json:"endpoint"`
	Name              string                `json:"name"`
	Desc              string                `json:"desc"`
	Template          string                `json:"template"`
	Inputs            []*Variable           `json:"inputs"`
	Outputs           []*Variable           `json:"outputs"`
	InputDefinitions  []*VariableDefinition `json:"inputDefinitions"`
	OutputDefinitions []*VariableDefinition `json:"outputDefinitions []"`
	HandlerFunc       func(f *Function) error
	CallbackFunc      func(variables []*Variable, err error) error
	Options           map[string]interface{} `json:"options"`
}
type Service struct {
	Functions  []*Function
	Client     *Client
	ServiceUrl string
}

func NewService(baseUrl string, serviceUrl string) *Service {
	client := NewYouLinkClient()
	client.Init(baseUrl)
	return &Service{
		Client:     client,
		Functions:  []*Function{},
		ServiceUrl: serviceUrl,
	}
}
func (s *Service) AddFunction(function ...*Function) {
	s.Functions = append(s.Functions, function...)
}
func (s *Service) RegisterHarukaHandler(engine *haruka.Engine) {
	for _, function := range s.Functions {
		function.Options = map[string]interface{}{
			"url": fmt.Sprintf("%s/youlink/%s", s.ServiceUrl, function.Endpoint),
		}
		engine.Router.POST(fmt.Sprintf("/youlink/%s", function.Endpoint), func(context *haruka.Context) {
			var requestBody ResponseBody
			err := context.ParseJson(&requestBody)
			if err != nil {
				AbortErrorWithStatus(err, context, http.StatusBadRequest)
				return
			}
			function.Inputs = requestBody.Inputs

			function.CallbackFunc = func(variables []*Variable, err error) error {
				s.Client.Callback(requestBody.CallbackId, variables, err)
				return nil
			}

			err = function.HandlerFunc(function)
			if err != nil {
				AbortErrorWithStatus(err, context, http.StatusBadRequest)
				return
			}
		})
	}
}
func (s *Service) RegisterFunction() error {
	err := s.Client.RegisterFunctions(s.Functions...)
	if err != nil {
		return err
	}
	return nil
}
func (f *Function) GetInput(key string) interface{} {
	for _, inputVar := range f.Inputs {
		if inputVar.Name == key {
			return inputVar.Value
		}
	}
	return nil
}
func (f *Function) GetInputString(key string) string {
	inputVar := f.GetInput(key)
	if inputVar == nil {
		return ""
	}
	if value, ok := inputVar.(string); ok {
		return value
	}
	return ""
}
func (f *Function) GetInputInt(key string) int {
	inputVar := f.GetInput(key)
	if inputVar == nil {
		return 0
	}
	if value, ok := inputVar.(int); ok {
		return value
	}
	return 0
}
func AbortErrorWithStatus(err error, context *haruka.Context, status int) {
	context.JSONWithStatus(map[string]interface{}{
		"success": false,
		"reason":  err.Error(),
	}, status)
}
