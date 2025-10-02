package api

import (
	"context"
	"time"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
)

// CBSStatusAPI is the core banking system service API for getting CBSAuth status.
type CBSStatusAPI struct{}

func NewCBSStatusAPI() *CBSStatusAPI {
	return &CBSStatusAPI{}
}

func (cs *CBSStatusAPI) GetStatus(ctx context.Context) (cbs.Status, error) {
	return cbs.Status{
		SystemDate: time.Now().Format("2006-01-02"),
		IsEOD:      false,
		IsStandIn:  false,
	}, nil
}
