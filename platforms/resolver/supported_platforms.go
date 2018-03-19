package resolver

import (
	"github.com/s0nerik/goloc/platforms"
	"github.com/s0nerik/goloc/goloc"
)

var supportedPlatforms = []goloc.Platform{
	&platforms.Android{},
}