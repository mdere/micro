package main

import (
	auth_plugin "bitbucket.org/appgoplaces/travelplatform-system/lib/plugins/auth"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(auth_plugin.New())
}
