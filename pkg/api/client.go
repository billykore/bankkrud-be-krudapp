package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Client is the API client that performs all operations
// against a TapMoney service.
type Client struct {
	addr   string
	client *http.Client
}

// NewClient creates a new instance of the TapMoney API client with the given base URL.
func NewClient(client *http.Client, addr string) *Client {
	return &Client{
		addr:   addr,
		client: client,
	}
}

func (c *Client) Inquiry(ctx context.Context, request InquiryRequest) (Response[InquiryResponse], error) {
	url := c.addr + "/api/v1/tapmoney/inquiry"
	var apiRes Response[InquiryResponse]
	err := c.doRequest(ctx, http.MethodPost, url, request, &apiRes)
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	return apiRes, nil
}

func (c *Client) Payment(ctx context.Context, request PaymentRequest) (Response[PaymentResponse], error) {
	url := c.addr + "/api/v1/tapmoney/payment"
	var apiRes Response[PaymentResponse]
	err := c.doRequest(ctx, http.MethodPost, url, request, &apiRes)
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	return apiRes, nil
}

func (c *Client) doRequest(ctx context.Context, method, url string, requestBody any, responseBody any) error {
	bBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bBody))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, responseBody)
	if err != nil {
		return err
	}
	return nil
}
