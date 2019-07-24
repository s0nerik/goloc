package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform
var sources map[string]goloc.Source

func RegisterPlatform(p goloc.Platform) {
	platforms = append(platforms, p)
}

func RegisterSource(s goloc.Source) {
	sources[s.Name()] = s
}

func GetPlatform(name string) goloc.Platform {
	for _, p := range platforms {
		for _, n := range p.Names() {
			if n == name {
				return p
			}
		}
	}
	return nil
}

func GetSource(name string) goloc.Source {
	return sources[name]
}
