package infra

import (
	"github.com/google/wire"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/account"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/cbs"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/payment"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transaction"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/transfer"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/api"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/handler"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/server"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/service"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/storage/repo"
)

var ProviderSet = wire.NewSet(
	api.NewCBSAccountAPI, wire.Bind(new(account.Repository), new(*api.CBSAccountAPI)),
	api.NewCBSStatusAPI, wire.Bind(new(cbs.Service), new(*api.CBSStatusAPI)),
	api.NewCBSTransferAPI, wire.Bind(new(transfer.Service), new(*api.CBSTransferAPI)),
	api.NewPaymentGateway, wire.Bind(new(payment.Service), new(*api.PaymentGateway)),
	repo.NewTransactionRepo, wire.Bind(new(transaction.Repository), new(*repo.TransactionRepo)),
	repo.NewUserRepo, wire.Bind(new(user.Repository), new(*repo.UserRepo)),
	service.NewAuthService, wire.Bind(new(user.AuthService), new(*service.AuthService)),
	handler.NewTransferHandler,
	handler.NewTapMoneyHandler,
	handler.NewAuthenticationHandler,
	handler.NewUserHandler,
	handler.NewTransactionHandler,
	server.NewHTTP,
)
