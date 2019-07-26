package registry

import (
	"github.com/s0nerik/goloc/goloc"
)

var platforms []goloc.Platform

func RegisterPlatform(p goloc.Platform) {
	platforms = append(platforms, p)
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
