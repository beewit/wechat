package util

import (
	"time"
	"fmt"
	"crypto/md5"
	"strconv"
)

func GenerateNonceStr() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}