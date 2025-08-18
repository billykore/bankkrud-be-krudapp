package usecase

import (
	"github.com/google/wire"
	"go.bankkrud.com/backend/svc/tapmoney/internal/usecase/tapmoney"
)

var ProviderSet = wire.NewSet(
	tapmoney.NewUsecase,
)
