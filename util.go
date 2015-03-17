package clover

import (
	"time"
)

func addSecondUnix(sec int) int64 {
	if sec == 0 {
		return 0
	}

	return time.Now().UTC().Truncate(time.Nanosecond).Add(time.Second * time.Duration(sec)).Unix()
}

func isExpireUnix(expires int64) bool {
	return time.Now().UTC().Truncate(time.Nanosecond).Unix() > expires
}
