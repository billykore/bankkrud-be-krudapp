package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

// CBSAuth is the core banking system service API for authenticating.
type CBSAuth struct {
	client   *http.Client
	addr     string
	username string
	password string
}

func NewCBS(cfg *config.Configs, client *http.Client) *CBSAuth {
	return &CBSAuth{
		client:   client,
		addr:     cfg.CBS.Addr,
		username: cfg.CBS.Username,
		password: cfg.CBS.Password,
	}
}

func (c *CBSAuth) GetToken(ctx context.Context) (string, error) {
	tokenURL := c.addr + "/token"

	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)
	data.Set("grant_type", "password")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to perform get token, status code: " + resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiRes CBSGetTokenResponse
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return "", err
	}

	return apiRes.AccessToken, nil
}
