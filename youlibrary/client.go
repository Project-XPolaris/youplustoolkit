package youlibrary

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewYouLibraryClient() *Client {
	return &Client{
		client: resty.New(),
	}
}

func (c *Client) Init(baseUrl string) {
	c.client.SetBaseURL(baseUrl)
}

type MatchSubjectRequestBody struct {
	Keyword     string `json:"keyword"`
	SubjectType string `json:"Type"`
}

func (c *Client) MatchVideoInfo(keyword string, subjectType string) (*MatchSubjectResponse, error) {
	body := MatchSubjectResponse{}
	_, err := c.client.R().SetResult(&body).SetBody(MatchSubjectRequestBody{
		Keyword:     keyword,
		SubjectType: subjectType,
	}).Post("/match")
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *Client) GetSubjectById(id uint) (*GetSubjectResponse, error) {
	body := GetSubjectResponse{}
	_, err := c.client.R().SetResult(&body).Get(fmt.Sprintf("/subject/%d", id))
	if err != nil {
		return nil, err
	}
	return &body, nil
}
