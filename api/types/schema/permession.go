package schema

import "time"

type FilePermession struct {
	Fileid     string
	Username   string
	Shorturl   int64
	Expiredate time.Time
}

type FileExpireDate struct {
	Expiredate string
}
