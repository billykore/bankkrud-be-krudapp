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
	body, err := json.Marshal(request)
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	var apiRes Response[InquiryResponse]
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return Response[InquiryResponse]{}, err
	}
	return apiRes, nil
}

func (c *Client) Payment(ctx context.Context, request PaymentRequest) (Response[PaymentResponse], error) {
	url := c.addr + "/api/v1/tapmoney/payment"
	body, err := json.Marshal(request)
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	var apiRes Response[PaymentResponse]
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return Response[PaymentResponse]{}, err
	}
	return apiRes, nil
}
