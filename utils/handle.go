package utils

import (
	"github.com/gosimple/slug"
)

func GetHandle(s string) string {
	return slug.Make(s)
}
