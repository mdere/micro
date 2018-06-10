package auth_plugin

// stats_auth enables basic auth on the /stats endpoint

import (
	"net/http"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"

	"github.com/micro/go-log"
)

type auth struct {
}

func (a *auth) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (a *auth) Commands() []cli.Command {
	return nil
}

func (a *auth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Log("------------------------")
			defer log.Log("------------------------")
			log.Log(r)
			return
		})
	}
}

func (a *auth) Init(ctx *cli.Context) error {
	return nil
}

func (a *auth) String() string {
	return "auth"
}

func New() plugin.Plugin {
	return new(auth)
}
