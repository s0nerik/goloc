package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform
var sources map[string]goloc.Source

// RegisterPlatform adds a new platform into the registry.
func RegisterPlatform(p goloc.Platform) {
	platforms = append(platforms, p)
}

// FindPlatform looks up and, if succeeds, returns a Platform given its name.
func FindPlatform(name string) goloc.Platform {
	for _, p := range platforms {
		for _, n := range p.Names() {
			if n == name {
				return p
			}
		}
	}
	return nil
}

func RegisterSource(s goloc.Source) {
	sources[s.Name()] = s
}

func GetSource(name string) goloc.Source {
	return sources[name]
}
