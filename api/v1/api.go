package v1

import (
	"transfer/api/v1/file"
	"transfer/api/v1/user"

	"github.com/google/wire"
)

type (
	File = file.API
	User = user.API
)

var APISet = wire.NewSet(
	file.FileSet,
	user.UserSet,
)
