package model

import (
	"transfer/internal/app/model/file"
	"transfer/internal/app/model/permession"
	"transfer/internal/app/model/user"

	"github.com/google/wire"
)

type (
	UserModel       = file.FileRepo
	FileModel       = user.UserRepo
	PermessionModel = permession.PermessionRepo
	User            = user.User
	File            = file.File
	Permession      = permession.Permession
)

var ModelSet = wire.NewSet(
	file.FileSet,
	user.UserSet,
	permession.PermessionSet,
)
