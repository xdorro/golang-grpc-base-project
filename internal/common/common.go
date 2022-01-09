// Package common
package common

import (
	"time"
)

const (
	// KeyServiceRoles service roles key
	KeyServiceRoles       = "service:roles"
	KeyServiceRolesExpire = 1 * 24 * time.Hour
)

const (
	CtxUserID string = "user_id"
)
