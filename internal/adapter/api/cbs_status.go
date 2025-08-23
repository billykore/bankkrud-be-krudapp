package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.bankkrud.com/backend/svc/tapmoney/internal/domain/cbs"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
)

// CBSStatusAPI is the core banking system service API for getting CBSAuth status.
type CBSStatusAPI struct {
	client   *http.Client
	cbsAuth  *CBSAuth
	addr     string
	username string
	password string
}

func NewCBSStatusAPI(cfg *config.Configs, client *http.Client, cbsAuth *CBSAuth) *CBSStatusAPI {
	return &CBSStatusAPI{
		client:   client,
		cbsAuth:  cbsAuth,
		addr:     cfg.CBS.Addr,
		username: cfg.CBS.Username,
		password: cfg.CBS.Password,
	}
}

func (cs *CBSStatusAPI) GetStatus(ctx context.Context) (cbs.Status, error) {
	statusUrl := cs.addr + "/api/ref/core-status"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, statusUrl, nil)
	if err != nil {
		return cbs.Status{}, err
	}

	token, err := cs.cbsAuth.getToken(ctx)
	if err != nil {
		return cbs.Status{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cs.client.Do(req)
	if err != nil {
		return cbs.Status{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return cbs.Status{}, errors.New("failed to perform get status, status code: " + resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return cbs.Status{}, err
	}

	var apiRes CBSGetStatusResponse
	err = json.Unmarshal(b, &apiRes)
	if err != nil {
		return cbs.Status{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if apiRes.StatusCode != CBSSuccessResponseCode {
		return cbs.Status{}, fmt.Errorf("[%v]: %v", apiRes.StatusCode, apiRes.StatusDescription)
	}

	return cbs.Status{
		SystemDate: apiRes.Data.SystemDate,
		IsEOD:      apiRes.Data.IsEOD(),
		IsStandIn:  apiRes.Data.IsStandIn(),
	}, nil
}
