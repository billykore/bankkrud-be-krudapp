package usecase

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/tapmoney"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/transfer"
)

var ProviderSet = wire.NewSet(
	tapmoney.NewUsecase,
	transfer.NewUsecase,
)
