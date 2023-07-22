package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcgtrt/toll-tracker/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) AggregateInvoice(dist types.Distance) error {
	b, err := json.Marshal(dist)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the server responded with %d status code", resp.StatusCode)
	}
	return nil
}
