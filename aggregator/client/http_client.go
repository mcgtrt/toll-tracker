package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcgtrt/toll-tracker/types"
)

type HTTPClient struct {
	endpoint string
}

func NewHTTPClient(endpoint string) Client {
	return &HTTPClient{
		endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", c.endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the server responded with %d status code", resp.StatusCode)
	}
	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := &types.GetInvoiceRequest{
		OBUID: int32(id),
	}
	b, err := json.Marshal(invReq)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/invoice?obu=%d", c.endpoint, id)
	httpReq, err := http.NewRequest("GET", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the server responded with %d status code", resp.StatusCode)
	}
	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	return &inv, err
}
