package resolver

import (
	"errors"
	"fmt"
	"github.com/s0nerik/goloc/goloc"
)

func FindPlatform(name string) (goloc.Platform, error) {
	for _, p := range supportedPlatforms {
		for _, n := range p.Names() {
			if n == name {
				return p, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf(`Platform "%v" not found.`, name))
}