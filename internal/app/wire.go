//go:build wireinject
// +build wireinject

package app

import (
	apiV1 "transfer/api/v1"
	"transfer/internal/app/model"
	"transfer/internal/app/router"
	"transfer/internal/app/service"
	"transfer/internal/app/storage"

	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGinEngine,
		InitDB,
		InitMinIO,
		storage.StorageSet,
		router.RouterSet,
		apiV1.APISet,
		service.ServiceSet,
		model.ModelSet,
		InjectorSet,
	)
	return new(Injector), nil, nil

}
