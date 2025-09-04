package pkg

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/db/postgres"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/httpclient"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
)

var ProviderSet = wire.NewSet(
	validation.New,
	postgres.New,
	httpclient.New,
)
