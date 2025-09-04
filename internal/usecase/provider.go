package usecase

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/tapmoney"
)

var ProviderSet = wire.NewSet(
	tapmoney.NewUsecase,
)
