package service

import (
	"transfer/internal/app/service/file"
	"transfer/internal/app/service/permession"
	"transfer/internal/app/service/user"

	"github.com/google/wire"
)

type (
	User       = user.Service
	File       = file.Service
	Permession = permession.Service
)

var ServiceSet = wire.NewSet(
	file.ServiceSet,
	user.ServiceSet,
	permession.ServiceSet,
)
