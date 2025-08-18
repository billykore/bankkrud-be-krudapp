package pkg

import (
	"github.com/google/wire"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/db/postgres"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/httpclient"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/validation"
)

var ProviderSet = wire.NewSet(
	validation.New,
	postgres.New,
	httpclient.New,
)
