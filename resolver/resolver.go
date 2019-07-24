package resolver

import (
	"fmt"

	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/registry"
	// It's needed to register all available platforms in the registry before running an app.
	_ "github.com/s0nerik/goloc/platforms"
)

// FindPlatform looks up and, if succeeds, returns a Platform given its name.
func FindPlatform(name string) (goloc.Platform, error) {
	for _, p := range registry.Platforms() {
		for _, n := range p.Names() {
			if n == name {
				return p, nil
			}
		}
	}
	return nil, fmt.Errorf(`platform "%v" not found`, name)
}
