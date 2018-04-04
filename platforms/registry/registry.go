package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform

func Register(p goloc.Platform) {
	platforms = append(platforms, p)
}

func Platforms() []goloc.Platform {
	return platforms
}