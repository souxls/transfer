package storage

import (
	"transfer/internal/app/storage/file"

	"github.com/google/wire"
)

var StorageSet = wire.NewSet(
	file.StorageSet,
)
