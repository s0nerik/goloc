package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform
var sources map[string]goloc.Source

// Register adds a new platform into the registry.
func Register(p goloc.Platform) {
	platforms = append(platforms, p)
}

// Platforms returns all platforms registered in the registry.
func Platforms() []goloc.Platform {
	return platforms
}

func RegisterSource(s goloc.Source) {
	sources[s.Name()] = s
}

func GetSource(name string) goloc.Source {
	return sources[name]
}
