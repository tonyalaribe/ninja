package uilayer

import (
	"github.com/tonyalaribe/ninja/core"
	"github.com/tonyalaribe/ninja/uilayer/rest"
)

func Register(manager core.Manager) error {
	rest.Register(manager)
	return nil
}
