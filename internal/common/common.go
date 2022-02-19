// Package common
package common

import (
	"time"

	"github.com/xdorro/golang-grpc-base-project/pkg/slug"
)

const (
	// KeyServiceRoles service roles key
	KeyServiceRoles       = "service:roles"
	KeyServiceRolesExpire = 1 * 24 * time.Hour
)

const (
	CtxUserID string = "user_id"
)

// GetSlugOrMakeSlug returns slug or makes slug
func GetSlugOrMakeSlug(title string, url string) string {
	if url == "" {
		url = slug.Make(title)
	}

	return url
}
