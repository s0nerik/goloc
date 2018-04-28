package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform

// Register adds a new platform into the registry.
func Register(p goloc.Platform) {
	platforms = append(platforms, p)
}

// Platforms returns all platforms registered in the registry.
func Platforms() []goloc.Platform {
	return platforms
}
