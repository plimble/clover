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

func copyClientToGrant(grantData *GrantData, client Client) {
	if grantData.UserID == "" {
		grantData.UserID = client.GetUserID()
	}

	if grantData.ClientID == "" {
		grantData.ClientID = client.GetClientID()
	}

	if len(grantData.Scope) == 0 {
		grantData.Scope = client.GetScope()
	}

	if len(grantData.GrantType) == 0 {
		grantData.GrantType = client.GetGrantType()
	}
}
