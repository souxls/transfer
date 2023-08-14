package errors

import "errors"

var (
	ErrUserExist = errors.New("user is exist")
	ErrPassword  = errors.New("user or password is incorrect")
	ErrSgin      = errors.New("sso sign is incorrect")
	ErrSaveFaild = errors.New("file save error")
	ErrDelFaild  = errors.New("delete failed")
)
