package validation

import (
	"crypto/md5"
	"fmt"
)

func CheckHttpSign(uri string, xSSOUID string, xSSORspIp string, userAgent string, tokenSecret string) string {
	strmd5 := []byte(uri + xSSOUID + xSSORspIp + userAgent + tokenSecret)
	return fmt.Sprintf("%x", md5.Sum(strmd5))
}
